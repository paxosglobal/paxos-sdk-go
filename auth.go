package paxos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// tokenManager handles OAuth2 client credentials token management.
type tokenCredentials struct {
	clientSecret string
	token        string
}

type tokenManager struct {
	clientID   string
	oauthURL   string
	scopes     []string
	httpClient *http.Client

	// sem is a binary semaphore (buffered channel of capacity 1) standing in
	// for a plain sync.Mutex: a plain mutex offers no way to bound a Lock()
	// call by a context deadline, which getToken needs so a caller with a
	// short remaining budget doesn't block on another goroutine's slow (or
	// stalled) in-flight refresh past its own deadline.
	sem          chan struct{}
	cred         *tokenCredentials
	expiresAt    time.Time
	lastFetchErr error
}

// tokenResponse represents the OAuth2 token response.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// newTokenManager creates a new token manager.
func newTokenManager(clientID, clientSecret, oauthURL string, httpClient *http.Client, scopes ...string) *tokenManager {
	return &tokenManager{
		clientID:   clientID,
		cred:       &tokenCredentials{clientSecret: clientSecret},
		oauthURL:   oauthURL,
		scopes:     scopes,
		httpClient: httpClient,
		sem:        make(chan struct{}, 1),
	}
}

// lockContext acquires tm.sem, bounded by ctx. Returns ctx.Err() (never nil)
// if ctx is done before the semaphore becomes available.
func (tm *tokenManager) lockContext(ctx context.Context) error {
	select {
	case tm.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (tm *tokenManager) unlock() {
	<-tm.sem
}

// getTokenWithClient uses the logical call's HTTP client, whose timeout is
// enforced by ctx. Standalone getToken callers retain the manager's default.
func (tm *tokenManager) getTokenWithClient(ctx context.Context, httpClient *http.Client) (string, error) {
	if err := tm.lockContext(ctx); err != nil {
		return "", err
	}
	defer tm.unlock()

	// Check if we have a valid token
	if tm.cred.token != "" && time.Now().Before(tm.expiresAt) {
		return tm.cred.token, nil
	}

	// Fetch a new token
	return tm.fetchTokenLocked(ctx, httpClient)
}

// peekToken returns the current cached token without triggering a refresh.
//
// Non-blocking: callers (do()'s retry loop) use this only to decide whether
// to log "refreshing OAuth2 token", immediately before also calling the
// ctx-bound getToken(ctx) for the actual (possibly blocking) token fetch. A
// blocking peek here — even one bounded by ctx — would let a caller queue
// behind another goroutine's in-flight refresh *twice* in a row, for no
// benefit: it's not this call's only, or even primary, way to get a token.
// If the semaphore isn't immediately available, returns "" rather than
// waiting for it at all; the caller falls back to comparing against
// getToken's result, so this only ever costs an extra debug log line.
func (tm *tokenManager) peekToken() string {
	select {
	case tm.sem <- struct{}{}:
	default:
		return ""
	}
	defer tm.unlock()
	if tm.cred.token != "" && time.Now().Before(tm.expiresAt) {
		return tm.cred.token
	}
	return ""
}

// invalidateToken clears the rejected token only if it is still cached.
// Used when receiving a 401 Unauthorized response.
//
// Bounded by ctx: if another goroutine is currently refreshing and ctx
// expires first, this is a no-op rather than blocking past the caller's own
// deadline. That's safe because the caller's ctx is already past its
// deadline by then, so the very next getToken(ctx) call fails fast on its
// own (via the same ctx-bound lock) rather than ever needing the
// invalidated state to have taken effect.
func (tm *tokenManager) invalidateToken(ctx context.Context, rejectedToken string) {
	if err := tm.lockContext(ctx); err != nil {
		return
	}
	defer tm.unlock()

	if tm.cred.token != rejectedToken {
		return
	}
	tm.cred.token = ""
	tm.expiresAt = time.Time{}
}

// fetchTokenLocked fetches a new access token from the OAuth2 server.
// Uses client_secret_post authentication method (credentials in request body).
// Caller must hold tm.sem.
func (tm *tokenManager) fetchTokenLocked(ctx context.Context, httpClient *http.Client) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", tm.clientID)
	data.Set("client_secret", tm.cred.clientSecret)
	if len(tm.scopes) > 0 {
		data.Set("scope", strings.Join(tm.scopes, " "))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tm.oauthURL, strings.NewReader(data.Encode()))
	if err != nil {
		tm.lastFetchErr = fmt.Errorf("failed to create token request: %w", err)
		return "", tm.lastFetchErr
	}

	req.Header.Set(headerContentType, "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		tm.lastFetchErr = newNetworkError("Token request failed", err, "", req.URL)
		return "", tm.lastFetchErr
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		tm.lastFetchErr = newNetworkError("failed to read token response", readErr, "", finalURL(resp, req))
		return "", tm.lastFetchErr
	}
	tokenResp, valid := parseTokenResponse(body)
	if resp.StatusCode != http.StatusOK || !valid {
		tm.lastFetchErr = newAuthenticationError(resp, body)
		return "", tm.lastFetchErr
	}

	// Store the token with an expiry buffer to allow for clock skew, clamped
	// to at most half the token's own lifetime — an unclamped 60s buffer
	// would leave a short-lived token with zero safety margin at all.
	tm.cred.token = tokenResp.AccessToken
	expiryDuration := time.Duration(tokenResp.ExpiresIn) * time.Second
	buffer := 60 * time.Second
	if half := expiryDuration / 2; half < buffer {
		buffer = half
	}
	tm.expiresAt = time.Now().Add(expiryDuration - buffer)
	tm.lastFetchErr = nil

	return tm.cred.token, nil
}

// parseTokenResponse validates the token schema without returning parser errors
// that could expose token response contents. Optional fields use OAuth defaults.
func parseTokenResponse(body []byte) (tokenResponse, bool) {
	var fields map[string]json.RawMessage
	token := tokenResponse{ExpiresIn: 3600}
	if json.Unmarshal(body, &fields) != nil || fields == nil {
		return token, false
	}
	if json.Unmarshal(fields["access_token"], &token.AccessToken) != nil || strings.TrimSpace(token.AccessToken) == "" {
		return token, false
	}
	if raw, present := fields["token_type"]; present {
		var tokenType string
		if json.Unmarshal(raw, &tokenType) != nil || !strings.EqualFold(tokenType, bearerTokenType) {
			return token, false
		}
	}
	if raw, present := fields["expires_in"]; present {
		var seconds *float64
		if json.Unmarshal(raw, &seconds) != nil || seconds == nil || math.IsNaN(*seconds) || math.IsInf(*seconds, 0) || math.Trunc(*seconds) != *seconds {
			return token, false
		}
		// Clamp before integer/duration conversion so even very large JSON numbers
		// cannot overflow. Nonpositive lifetimes remain immediately expired.
		const maxExpiresInSeconds = 100 * 365 * 24 * 60 * 60
		switch {
		case *seconds < 0:
			token.ExpiresIn = 0
		case *seconds > maxExpiresInSeconds:
			token.ExpiresIn = maxExpiresInSeconds
		default:
			token.ExpiresIn = int(*seconds)
		}
	}
	return token, true
}
