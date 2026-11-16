package paxos

import (
	"bytes"
	"context"
	"crypto"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	mathrand "math/rand"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// defaultBaseURL is the production API base URL.
	defaultBaseURL = "https://api.paxos.com/v2"

	// sandboxBaseURL is the sandbox API base URL.
	sandboxBaseURL = "https://api.sandbox.paxos.com/v2"

	// defaultOAuthURL is the OAuth2 token endpoint.
	defaultOAuthURL = "https://oauth.paxos.com/oauth2/token"

	// sandboxOAuthURL is the OAuth2 token endpoint for the sandbox environment.
	sandboxOAuthURL = "https://oauth.sandbox.paxos.com/oauth2/token"

	// maxRetries is the maximum number of retry attempts for retryable requests.
	maxRetries = 3

	// defaultCallTimeout is the default total logical call budget.
	defaultCallTimeout = 30 * time.Second

	// sdkLanguage is the language identifier for telemetry headers.
	sdkLanguage = "go"
)

// retryableStatusCodes defines HTTP status codes that trigger automatic retry
// with exponential backoff.
var retryableStatusCodes = map[int]bool{
	http.StatusTooManyRequests:     true, // 429
	http.StatusInternalServerError: true, // 500
	http.StatusBadGateway:          true, // 502
	http.StatusServiceUnavailable:  true, // 503
	http.StatusGatewayTimeout:      true, // 504
}

// RequestHook is called once per HTTP attempt, before signing. It may add or
// change headers but must not change the method, URL or body.
type RequestHook func(req *http.Request) error

// ResponseHook is called once per attempt that received a response, before
// status handling. resp.Body is buffered and may be read freely.
type ResponseHook func(req *http.Request, resp *http.Response) error

// parseResponseMetadata extracts metadata from HTTP response headers.
func parseResponseMetadata(resp *http.Response) *ResponseMetadata {
	return &ResponseMetadata{
		StatusCode: resp.StatusCode,
		RequestID:  resp.Header.Get(headerRequestID),
	}
}

type requestSigner struct{ key *crypto.Signer }

// Client is the main SDK client for interacting with the Paxos API.
type Client struct {
	baseURL            string
	sandbox            bool
	customBaseURL      bool
	customOAuthURL     bool
	oauthURL           string
	scopes             []string
	httpClient         *http.Client
	tokenManager       *tokenManager
	noAuth             bool
	maxRetries         int
	callTimeout        time.Duration
	callAttemptTimeout time.Duration
	optionErr          error

	requestHook  RequestHook
	responseHook ResponseHook

	// Request signing fields. When both are set, every HTTP request is signed
	// with an Ed25519 JWS and the signature is sent in the Paxos-Signature header.
	// Request signing requires contacting Paxos support to enable on the OAuth client.
	signingKey   *requestSigner
	signingKeyID string
	signingErr   error

	// logger is an optional structured logger for debug output.
	logger *slog.Logger

	// Services holds all generated API service instances.
	// Embedded so fields are promoted — use client.Transfers directly.
	Services
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// RequestOption is a functional option for per-request overrides.
type RequestOption func(*requestOptions)

// ResponseMetadata contains metadata from an API response.
// Use WithResponseMetadata to capture this from any service method call.
// Metadata is populated for both success and error responses.
type ResponseMetadata struct {
	// RequestID is the unique identifier assigned to the request, useful for
	// debugging and support inquiries.
	RequestID string

	// StatusCode is the HTTP status code returned by the API.
	StatusCode int

	// RawBody is the unmodified body of the final response, kept even if decoding fails; empty (non-nil) for an empty body.
	RawBody []byte
}

// requestOptions holds per-request configuration that overrides client defaults.
type requestOptions struct {
	maxRetries         *int
	callTimeout        *time.Duration
	callAttemptTimeout *time.Duration
	optionErr          error
	extraHeaders       http.Header
	responseMetadata   *ResponseMetadata
	maxItems           int
	maxDuration        time.Duration
	replayProtected    bool
}

func withReplayProtected() RequestOption {
	return func(o *requestOptions) { o.replayProtected = true }
}

// WithMaxItems caps the total number of items an auto-paginating Iterator
// (e.g. from client.Transfers.ListTransfers) will yield, regardless of how
// many pages remain on the server. Ignored by non-paginated requests.
func WithMaxItems(n int) RequestOption {
	return func(o *requestOptions) {
		o.maxItems = n
		if n <= 0 {
			o.optionErr = fmt.Errorf("max items must be positive")
		}
	}
}

// WithMaxDuration bounds the total wall-clock time an auto-paginating
// Iterator may spend across all of its page fetches. Once exceeded,
// Iterator.Next returns false and Iterator.Err reports an error wrapping
// context.DeadlineExceeded. Ignored by non-paginated requests.
//
// A deadline on the ctx passed to the originating service method achieves
// the same effect and composes with WithMaxDuration — whichever deadline is
// earlier wins.
func WithMaxDuration(d time.Duration) RequestOption {
	return func(o *requestOptions) {
		o.maxDuration = d
		if d <= 0 {
			o.optionErr = fmt.Errorf("max duration must be positive")
		}
	}
}

// WithExtraHeaders adds additional HTTP headers for a single request.
func WithExtraHeaders(h http.Header) RequestOption {
	return func(o *requestOptions) {
		o.extraHeaders = h.Clone()
	}
}

// parseRequestOptions applies a variadic RequestOption list into a fresh
// requestOptions value. Shared by do, paginate, and paginateTyped.
func parseRequestOptions(opts ...RequestOption) requestOptions {
	var reqOpts requestOptions
	for _, opt := range opts {
		opt(&reqOpts)
	}
	return reqOpts
}

// WithResponseMetadata captures HTTP response metadata (RequestID, StatusCode)
// from a service method call into the provided ResponseMetadata value.
// Metadata is populated for both success and error responses, making it
// useful for debugging failed requests. Safe to use as long as the same
// instance is not shared across concurrent requests.
//
// Usage:
//
//	var meta paxos.ResponseMetadata
//	transfer, err := client.Transfers.GetTransfer(ctx, "txn_123", paxos.WithResponseMetadata(&meta))
//	if err != nil {
//	    log.Printf("request failed: request_id=%s status=%d", meta.RequestID, meta.StatusCode)
//	}
//	fmt.Println("request_id:", meta.RequestID, "status:", meta.StatusCode)
func WithResponseMetadata(meta *ResponseMetadata) RequestOption {
	return func(o *requestOptions) {
		o.responseMetadata = meta
	}
}

// WithBaseURL sets a custom base URL for the API. Trailing slashes are removed.
// It cannot be combined with WithSandbox and must be paired with WithTokenURL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL == "" {
			return
		}
		c.baseURL = trimPathSlashes(baseURL)
		c.customBaseURL = true
	}
}

// WithSandbox configures the client to use the sandbox environment. It cannot
// be combined with WithBaseURL or WithTokenURL.
func WithSandbox() Option {
	return func(c *Client) {
		c.sandbox = true
	}
}

// WithHTTPClient sets a custom HTTP client. The SDK copies it and overwrites the
// copy's Timeout with its own per-attempt budget (WithDefaultCallAttemptTimeout),
// so a Timeout set on httpClient is ignored; use WithDefaultCallTimeout and
// WithDefaultCallAttemptTimeout instead. Its Transport is used as given.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithRequestHook sets the request hook. Setting it again replaces the previous hook.
func WithRequestHook(hook RequestHook) Option {
	return func(c *Client) {
		c.requestHook = hook
	}
}

// WithResponseHook sets the response hook. Setting it again replaces the previous hook.
func WithResponseHook(hook ResponseHook) Option {
	return func(c *Client) {
		c.responseHook = hook
	}
}

// WithTokenURL sets a custom OAuth2 token endpoint URL. Trailing slashes are
// removed. Useful for testing against a local token server. It cannot be
// combined with WithSandbox and must be paired with WithBaseURL.
func WithTokenURL(url string) Option {
	return func(c *Client) {
		if url == "" {
			return
		}
		c.oauthURL = trimPathSlashes(url)
		c.customOAuthURL = true
	}
}

func trimPathSlashes(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawPath = strings.TrimRight(u.RawPath, "/")
	return u.String()
}

// WithScopes sets the OAuth2 scopes to request when obtaining access tokens.
// Scopes control which API resources the client is permitted to access.
// At least one scope is required for any client that authenticates.
func WithScopes(scopes ...Scope) Option {
	return func(c *Client) {
		c.scopes = make([]string, len(scopes))
		for i, scope := range scopes {
			c.scopes[i] = string(scope)
		}
	}
}

// WithRequestSigning enables EdDSA or ES256 JWS request signing for every HTTP request.
// When configured, each request will include a Paxos-Signature header containing
// a JWS compact serialization that signs the request method, path, timestamp,
// and body.
//
// Parameters:
//   - privateKeyHex: Ed25519 private key hex-encoded (64 bytes = 128 hex chars: the
//     32-byte seed followed by the 32-byte public key), or an unencrypted PEM: PKCS#8
//     ("BEGIN PRIVATE KEY") with an Ed25519 or P-256 key, or SEC1 ("BEGIN EC PRIVATE
//     KEY") with a P-256 key. P-256 keys sign with ES256 and Ed25519 keys with EdDSA.
//   - keyID: UUID of the registered public key.
//
// Note: Request signing requires contacting Paxos support to enable on the
// OAuth client before it can be used.
func WithRequestSigning(privateKeyHex, keyID string) Option {
	return func(c *Client) {
		c.signingKeyID = keyID
		key, err := parseSigningKey(privateKeyHex, "signing key")
		if err != nil {
			c.signingErr = err
			return
		}
		c.signingKey = &requestSigner{&key}
	}
}

// NewClient creates a new Paxos API client.
//
// clientID and clientSecret are the OAuth2 credentials used to authenticate
// with the Paxos API. Both are required.
// An error is returned if credentials are missing or empty.
func NewClient(clientID, clientSecret string, opts ...Option) (*Client, error) {
	client := &Client{
		baseURL:     defaultBaseURL,
		callTimeout: defaultCallTimeout,
		maxRetries:  maxRetries,
		oauthURL:    defaultOAuthURL,
		httpClient: &http.Client{
			Timeout: defaultCallTimeout,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	if client.optionErr != nil {
		return nil, client.optionErr
	}
	if client.httpClient == nil {
		return nil, fmt.Errorf("HTTP client must not be nil")
	}
	if client.sandbox {
		if client.customBaseURL || client.customOAuthURL {
			return nil, fmt.Errorf("sandbox cannot be combined with a custom base URL or token URL; set both URLs explicitly with sandbox disabled")
		}
		client.baseURL = sandboxBaseURL
		client.oauthURL = sandboxOAuthURL
	} else if client.customBaseURL != client.customOAuthURL {
		return nil, fmt.Errorf("base URL and token URL must be set together; pass both, or neither to use the production or sandbox defaults")
	}

	// Validate credentials if auth is enabled
	if !client.noAuth {
		if clientID == "" && clientSecret == "" {
			return nil, fmt.Errorf("client_id and client_secret are required")
		}
		if clientID == "" {
			return nil, fmt.Errorf("client_id is required")
		}
		if clientSecret == "" {
			return nil, fmt.Errorf("client_secret is required")
		}
		if len(client.scopes) == 0 || slices.ContainsFunc(client.scopes, func(s string) bool { return strings.TrimSpace(s) == "" }) {
			return nil, fmt.Errorf("scopes are required: provide at least one OAuth scope, e.g. paxos.WithScopes(paxos.ScopeTransferReadTransfer)")
		}
		client.tokenManager = newTokenManager(clientID, clientSecret, client.oauthURL, client.httpClient, client.scopes...)
	}

	// Surface any signing parse error from WithRequestSigning.
	if client.signingErr != nil {
		return nil, client.signingErr
	}

	// Validate request signing configuration: both or neither must be set.
	hasKey := client.signingKey != nil
	hasKeyID := client.signingKeyID != ""
	if hasKey != hasKeyID {
		return nil, fmt.Errorf("both signing key and key ID must be provided for request signing")
	}

	client.initDefaultLogger()

	client.Services = client.initServices()

	return client, nil
}

func bufferBody(resp *http.Response) ([]byte, func(), error) {
	buffered, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	reset := func() {
		var rest io.Reader = bytes.NewReader(buffered)
		if readErr != nil {
			rest = io.MultiReader(rest, errReader{readErr})
		}
		resp.Body = io.NopCloser(rest)
	}
	reset()
	return buffered, reset, readErr
}

// do performs an HTTP request with authentication, error handling, and retry logic.
// It handles:
// - OAuth2 token injection with automatic refresh
// - JSON marshaling/unmarshaling
// - Error parsing
// - Request ID extraction
// - Rate limiting with exponential backoff
// - Retry on 401 (single retry after token refresh)
// - Retry on 429/5xx (up to maxRetries with backoff)
const invalidPathPrefix = "\x00invalid-path-parameter:"

func (c *Client) do(ctx context.Context, method, path string, body, result any, opts ...RequestOption) (*ResponseMetadata, error) {
	// Apply per-request options
	reqOpts := parseRequestOptions(opts...)
	if reqOpts.responseMetadata != nil {
		*reqOpts.responseMetadata = ResponseMetadata{}
	}

	if reqOpts.optionErr != nil {
		return nil, reqOpts.optionErr
	}
	if c == nil {
		return nil, errors.New("paxos: service used without a client; obtain services from paxos.NewClient")
	}
	if name, ok := strings.CutPrefix(path, invalidPathPrefix); ok {
		return nil, fmt.Errorf("path parameter %q must not be empty", name)
	}
	for key := range reqOpts.extraHeaders {
		if sdkOwnedHeaders[http.CanonicalHeaderKey(key)] {
			return nil, fmt.Errorf("extra header %q is SDK-owned and cannot be set", key)
		}
	}
	maxRetries := c.maxRetries
	if reqOpts.maxRetries != nil {
		maxRetries = *reqOpts.maxRetries
	}
	effectiveTimeout := c.callTimeout
	if reqOpts.callTimeout != nil {
		effectiveTimeout = *reqOpts.callTimeout
	}
	if effectiveTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, effectiveTimeout)
		defer cancel()
	}
	if err := ctx.Err(); err != nil {
		return nil, withRequestID(err, "")
	}
	// Context bounds the logical call; Client.Timeout bounds each individual
	// HTTP attempt through full body consumption, including OAuth requests.
	httpClient := *c.httpClient
	httpClient.Timeout = c.callAttemptTimeout
	if reqOpts.callAttemptTimeout != nil {
		httpClient.Timeout = *reqOpts.callAttemptTimeout
	}

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = marshalRequestBody(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	url := c.baseURL + path

	logPath, _, _ := strings.Cut(path, "?")
	c.logDebug("paxos: " + method + " " + logPath)

	// Retry loop for 429 rate limiting
	authRefreshed := false
	lastSentID := ""
	requestStart := time.Now()
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Create a fresh reader for each attempt (required for retries)
		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		// Create request
		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set(headerContentType, mediaTypeJSON)
		req.Header.Set("Accept", mediaTypeJSON)

		// Inject telemetry headers
		req.Header.Set("Paxos-SDK-Language", sdkLanguage)
		req.Header.Set("Paxos-SDK-Version", version)
		requestID := generateUUID()
		req.Header.Set(headerRequestID, requestID)

		// Apply per-request extra headers
		for key, values := range reqOpts.extraHeaders {
			if strings.EqualFold(key, headerRequestID) {
				req.Header.Del(headerRequestID)
			}
			for _, v := range values {
				req.Header.Add(key, v)
			}
		}

		var rejectedToken string
		// Add authentication if enabled
		if !c.noAuth {
			tokenBefore := c.tokenManager.peekToken()
			token, err := c.tokenManager.getTokenWithClient(ctx, &httpClient)
			if err != nil {
				return nil, withRequestID(fmt.Errorf("failed to get access token: %w", err), lastSentID)
			}
			if tokenBefore == "" || tokenBefore != token {
				c.logDebug("paxos: refreshing OAuth2 token")
			}
			rejectedToken = token
			req.Header.Set("Authorization", bearerTokenType+" "+token)
		}

		// Call request hook if set
		if c.requestHook != nil {
			hookMethod, hookURL, hookHost := req.Method, req.URL.String(), req.Host
			ownedBefore := snapshotOwnedHeaders(req.Header)
			hookBody := &hookRequestBody{Reader: bytes.NewReader(bodyBytes)}
			if bodyBytes != nil {
				req.Body = hookBody
			}
			if err := c.requestHook(req); err != nil {
				return nil, fmt.Errorf("request hook failed: %w", err)
			}
			if req.URL == nil || req.Method != hookMethod || req.URL.String() != hookURL || req.Host != hookHost {
				return nil, errors.New("request hook must not change the request method or URL")
			}
			if !ownedHeadersUnchanged(req.Header, ownedBefore) {
				return nil, errors.New("request hook must not change or remove SDK-owned headers")
			}
			if req.Body != hookBody && req.Body != nil && req.Body != http.NoBody {
				changed, err := io.ReadAll(req.Body)
				req.Body.Close()
				if err != nil || !bytes.Equal(changed, bodyBytes) {
					return nil, errors.New("request hook must not change the request body")
				}
			} else if (req.Body == nil || req.Body == http.NoBody) && bodyBytes != nil {
				return nil, errors.New("request hook must not change the request body")
			}
			if bodyBytes != nil {
				req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				req.ContentLength = int64(len(bodyBytes))
				req.GetBody = func() (io.ReadCloser, error) {
					return io.NopCloser(bytes.NewReader(bodyBytes)), nil
				}
			} else {
				req.Body = nil
				req.ContentLength = 0
				req.GetBody = nil
			}
		}

		// Sign the final request after hooks have had a chance to change its
		// method, target or body. GetBody may still describe the original body
		// after a hook replaces Body, so snapshot Body itself and restore it
		// for both this attempt and any transport-level replay.
		if c.signingKey != nil {
			var signedBody []byte
			if req.Body != nil {
				signedBody, err = io.ReadAll(req.Body)
				req.Body.Close()
				if err != nil {
					return nil, fmt.Errorf("failed to read request body for signing: %w", err)
				}
			}
			req.ContentLength = int64(len(signedBody))
			req.GetBody = func() (io.ReadCloser, error) {
				if len(signedBody) == 0 {
					return http.NoBody, nil
				}
				return io.NopCloser(bytes.NewReader(signedBody)), nil
			}
			req.Body, _ = req.GetBody()
			c.logDebug("paxos: signing request with key " + c.signingKeyID)
			jws, err := signRequestWithKey(*c.signingKey.key, c.signingKeyID, req.Method, req.URL.RequestURI(), signedBody)
			if err != nil {
				return nil, fmt.Errorf("failed to sign request: %w", err)
			}
			req.Header.Set("Paxos-Signature", jws)
		}

		// Snapshot the request sent after hooks; hooks or transports may mutate req again.
		sentMethod := req.Method
		sentTarget := req.URL.Path
		sentRequestID := req.Header.Get(headerRequestID)
		lastSentID = sentRequestID
		resp, err := httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil, withRequestID(ctx.Err(), lastSentID)
			}
			if attempt < maxRetries && isNetworkRetrySafe(sentMethod) {
				if reqOpts.responseMetadata != nil {
					*reqOpts.responseMetadata = ResponseMetadata{}
				}
				if err := waitBackoff(ctx, attempt); err != nil {
					return nil, withRequestID(err, lastSentID)
				}
				continue
			}
			return nil, newNetworkError("request failed", err, sentRequestID, req.URL)
		}

		elapsed := time.Since(requestStart)

		// Extract response metadata
		respMeta := parseResponseMetadata(resp)
		if respMeta.RequestID == "" {
			respMeta.RequestID = sentRequestID
		}
		c.logDebug(fmt.Sprintf("paxos: %d %s %s (%dms)", resp.StatusCode, sentMethod, sentTarget, elapsed.Milliseconds()), "request_id", respMeta.RequestID)

		// Populate caller-provided response metadata immediately after receiving
		// the HTTP response. This ensures metadata is available for ALL responses
		// (2xx, 4xx, 5xx), which is critical for debugging failed requests.
		if reqOpts.responseMetadata != nil {
			reqOpts.responseMetadata.RequestID = respMeta.RequestID
			reqOpts.responseMetadata.StatusCode = respMeta.StatusCode
		}

		// Call response hook if set
		if c.responseHook != nil {
			_, resetBody, _ := bufferBody(resp)
			if err := c.responseHook(req, resp); err != nil {
				resp.Body.Close()
				return nil, fmt.Errorf("response hook failed: %w", err)
			}
			resetBody()
		}

		// Handle 401 Unauthorized - invalidate token and retry once
		if resp.StatusCode == http.StatusUnauthorized && !c.noAuth && !authRefreshed {
			resp.Body.Close()
			if reqOpts.responseMetadata != nil {
				*reqOpts.responseMetadata = ResponseMetadata{}
			}
			attempt-- // The single auth refresh does not consume the general retry budget.
			authRefreshed = true
			c.tokenManager.invalidateToken(ctx, rejectedToken)
			continue
		}

		// Handle retryable status codes (429, 500, 502, 503, 504) - retry with backoff
		if retryableStatusCodes[resp.StatusCode] && isStatusRetrySafe(resp.StatusCode, sentMethod, reqOpts.replayProtected) {
			if attempt < maxRetries {
				resp.Body.Close()
				if reqOpts.responseMetadata != nil {
					*reqOpts.responseMetadata = ResponseMetadata{}
				}
				// Calculate backoff with jitter
				backoff := calculateBackoff(attempt)
				c.logDebug(fmt.Sprintf("paxos: retry attempt %d/%d after %d (backoff %s)", attempt+1, maxRetries, resp.StatusCode, backoff.Round(time.Millisecond)))
				timer := time.NewTimer(backoff)
				select {
				case <-timer.C:
				case <-ctx.Done():
					timer.Stop()
					return nil, withRequestID(ctx.Err(), lastSentID)
				}
				requestStart = time.Now()
				continue
			}
		}

		// Parse terminal HTTP errors once, including exhausted status retries.
		if resp.StatusCode >= 400 {
			errBody, _, readErr := bufferBody(resp)
			if reqOpts.responseMetadata != nil {
				reqOpts.responseMetadata.RawBody = errBody
			}
			if readErr != nil {
				if ctx.Err() != nil {
					return nil, withRequestID(ctx.Err(), lastSentID)
				}
				return nil, newNetworkError("failed to read error response", readErr, sentRequestID, finalURL(resp, req))
			}
			apiErr := parseError(resp, errBody, respMeta.RequestID)
			if reqOpts.responseMetadata != nil {
				// The body may supply a request ID that supersedes the early
				// header metadata made available to response hooks.
				reqOpts.responseMetadata.RequestID = apiErr.RequestID
			}
			return nil, apiErr
		}

		// Read response body for logging and decoding.
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			if ctx.Err() != nil {
				return nil, withRequestID(ctx.Err(), lastSentID)
			}
			if attempt < maxRetries && isNetworkRetrySafe(sentMethod) {
				if reqOpts.responseMetadata != nil {
					*reqOpts.responseMetadata = ResponseMetadata{}
				}
				if err := waitBackoff(ctx, attempt); err != nil {
					return nil, withRequestID(err, lastSentID)
				}
				continue
			}
			return nil, newNetworkError("failed to read response body", err, sentRequestID, finalURL(resp, req))
		}

		if reqOpts.responseMetadata != nil {
			reqOpts.responseMetadata.RawBody = respBody
		}

		// Success - decode response if result is provided
		if result != nil {
			if string(bytes.TrimSpace(respBody)) == "null" {
				return nil, newMalformedResponseError(resp.StatusCode, respMeta.RequestID, errors.New("null body"))
			}
			if err := json.Unmarshal(respBody, result); err != nil {
				return nil, newMalformedResponseError(resp.StatusCode, respMeta.RequestID, err)
			}
			attachRawJSON(result, respBody)
		}

		if err := ctx.Err(); err != nil {
			return nil, withRequestID(err, lastSentID)
		}
		return respMeta, nil
	}

	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}

// isNetworkRetrySafe reports whether an ambiguous network failure may be retried.
func isNetworkRetrySafe(method string) bool {
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}

func isStatusRetrySafe(status int, method string, replayProtected bool) bool {
	if status == http.StatusTooManyRequests || replayProtected {
		return true
	}
	return method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch
}

// calculateBackoff calculates the backoff duration for retry attempts.
// It implements exponential backoff with jitter, consistent across all SDK languages.
// Formula: base = 0.5s * 2^attempt, jitter = uniform [0.75, 1.25), capped at 30s.
func calculateBackoff(attempt int) time.Duration {
	base := time.Duration(float64(500*time.Millisecond) * math.Pow(2, float64(min(attempt, 7))))

	jitter := 0.75 + mathrand.Float64()*0.5
	delay := time.Duration(float64(base) * jitter)

	return min(delay, 30*time.Second)
}

// generateUUID generates a UUID v4 string.
func generateUUID() string {
	return uuid.NewString()
}

// Ptr returns a pointer to the given value of any type.
func Ptr[T any](v T) *T { return &v }

// String returns a pointer to the given string value.
// Use this to set optional string fields in typed params structs.
func String(v string) *string { return &v }

// Int returns a pointer to the given int value.
// Use this to set optional int fields in typed params structs.
func Int(v int) *int { return &v }

// Int32 returns a pointer to the given int32 value.
// Use this to set optional int32 fields in typed params structs (e.g., limit).
func Int32(v int32) *int32 { return &v }

// Int64 returns a pointer to the given int64 value.
// Use this to set optional int64 fields in typed params structs.
func Int64(v int64) *int64 { return &v }

// Bool returns a pointer to the given bool value.
// Use this to set optional bool fields in typed params structs.
func Bool(v bool) *bool { return &v }

// buildPath interpolates path parameters into a URL template, encoding each
// value for safe use in a URL path segment. Parameters are passed as
// alternating key/value pairs: buildPath("/transfers/{id}", "id", transferID).
func buildPath(template string, keysAndValues ...string) string {
	for i := 0; i < len(keysAndValues)-1; i += 2 {
		if strings.TrimSpace(strings.ReplaceAll(keysAndValues[i+1], "\ufeff", "")) == "" {
			return invalidPathPrefix + keysAndValues[i]
		}
		template = strings.Replace(template, "{"+keysAndValues[i]+"}", url.PathEscape(keysAndValues[i+1]), 1)
	}
	return template
}

func marshalRequestBody(body any) ([]byte, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var tree any
	if err := dec.Decode(&tree); err != nil {
		return raw, nil
	}
	return json.Marshal(omitEmptyCollections(tree, reflect.TypeOf(body), reflect.ValueOf(body)))
}

func omitEmptyCollections(v any, t reflect.Type, rv reflect.Value) any {
	t = indirectType(t)
	rv = indirectValue(rv)
	if rv.IsValid() && (t == nil || t.Kind() == reflect.Interface) {
		t = rv.Type()
	}
	if t != nil && t != timeType && (t.Implements(marshalerType) || reflect.PointerTo(t).Implements(marshalerType)) {
		t = nil
		rv = reflect.Value{}
	}
	switch val := v.(type) {
	case map[string]any:
		var fields map[string]structField
		var elem reflect.Type
		if t != nil {
			switch t.Kind() {
			case reflect.Struct:
				fields = structFields(t, rv)
			case reflect.Map:
				elem = t.Elem()
			}
		}
		var mapValues map[string]reflect.Value
		if rv.IsValid() && rv.Kind() == reflect.Map {
			mapValues = make(map[string]reflect.Value, rv.Len())
			iter := rv.MapRange()
			for iter.Next() {
				if name, ok := mapKeyName(iter.Key()); ok {
					mapValues[name] = iter.Value()
				}
			}
		}
		out := make(map[string]any, len(val))
		for k, child := range val {
			childType := elem
			var childValue reflect.Value
			if fields != nil {
				childType = fields[k].typ
				childValue = fields[k].val
			} else {
				childValue = mapValues[k]
			}
			pruned := omitEmptyCollections(child, childType, childValue)
			if isEmptyList(pruned) || (isEmptyObject(pruned) && (isMapType(childType) || (isValueStructType(childType) && isEmptyObject(child)))) {
				continue
			}
			out[k] = pruned
		}
		return out
	case []any:
		var elem reflect.Type
		isList := t != nil && (t.Kind() == reflect.Slice || t.Kind() == reflect.Array)
		if isList {
			elem = t.Elem()
		}
		out := make([]any, len(val))
		for i, child := range val {
			var childValue reflect.Value
			if isList && rv.IsValid() && i < rv.Len() {
				childValue = rv.Index(i)
			}
			out[i] = omitEmptyCollections(child, elem, childValue)
		}
		return out
	case string:
		if t == timeType && rv.IsValid() && rv.CanInterface() {
			if original, ok := rv.Interface().(time.Time); ok && val == original.Format(time.RFC3339Nano) {
				return original.UTC().Format(time.RFC3339Nano)
			}
		}
	}
	return v
}

var (
	timeType      = reflect.TypeOf(time.Time{})
	marshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
)

func indirectType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func mapKeyName(k reflect.Value) (string, bool) {
	if k.Kind() == reflect.String {
		return k.String(), true
	}
	if k.CanInterface() {
		if m, ok := k.Interface().(encoding.TextMarshaler); ok {
			if k.Kind() == reflect.Ptr && k.IsNil() {
				return "", false
			}
			b, err := m.MarshalText()
			return string(b), err == nil
		}
	}
	switch k.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(k.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(k.Uint(), 10), true
	}
	return "", false
}

func indirectValue(rv reflect.Value) reflect.Value {
	for rv.IsValid() && (rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface) {
		if rv.IsNil() {
			return reflect.Value{}
		}
		rv = rv.Elem()
	}
	return rv
}

func isMapType(t reflect.Type) bool {
	t = indirectType(t)
	return t != nil && t.Kind() == reflect.Map
}

func isValueStructType(t reflect.Type) bool {
	return t != nil && t.Kind() == reflect.Struct
}

type structField struct {
	typ reflect.Type
	val reflect.Value
}

func structFields(t reflect.Type, rv reflect.Value) map[string]structField {
	fields := map[string]structField{}
	embedded := map[string]structField{}
	ambiguous := map[string]bool{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "-" {
			continue
		}
		var fv reflect.Value
		if rv.IsValid() && rv.Kind() == reflect.Struct && rv.Type() == t {
			fv = rv.Field(i)
		}
		if name == "" && f.Anonymous {
			if inner := indirectType(f.Type); inner.Kind() == reflect.Struct {
				for k, ff := range structFields(inner, indirectValue(fv)) {
					if _, dup := embedded[k]; dup {
						ambiguous[k] = true
					}
					embedded[k] = ff
				}
				continue
			}
		}
		if !f.IsExported() {
			continue
		}
		if name == "" {
			name = f.Name
		}
		fields[name] = structField{f.Type, fv}
	}
	for k, ff := range embedded {
		if _, ok := fields[k]; !ok && !ambiguous[k] {
			fields[k] = ff
		}
	}
	return fields
}

func isEmptyList(v any) bool {
	list, ok := v.([]any)
	return ok && len(list) == 0
}

func isEmptyObject(v any) bool {
	obj, ok := v.(map[string]any)
	return ok && len(obj) == 0
}

func queryValue(v any) string {
	switch t := v.(type) {
	case time.Time:
		return t.UTC().Format(time.RFC3339Nano)
	case *time.Time:
		if t != nil {
			return t.UTC().Format(time.RFC3339Nano)
		}
	}
	return fmt.Sprintf("%v", v)
}

// appendQueryParams adds encoded params to any existing URL query parameters.
// Slice values are serialised as repeated keys (e.g. ids=a&ids=b).
func appendQueryParams(path string, params map[string]any) string {
	if len(params) == 0 {
		return path
	}
	q := url.Values{}
	for k, v := range params {
		switch val := v.(type) {
		case []string:
			for _, s := range val {
				q.Add(k, s)
			}
		case []any:
			for _, elem := range val {
				q.Add(k, queryValue(elem))
			}
		case []time.Time:
			for _, t := range val {
				q.Add(k, queryValue(t))
			}
		case []*time.Time:
			for _, t := range val {
				if t != nil {
					q.Add(k, queryValue(t))
				}
			}
		case *time.Time:
			if val != nil {
				q.Set(k, queryValue(val))
			}
		default:
			q.Set(k, queryValue(v))
		}
	}
	encoded := q.Encode()
	if encoded == "" {
		return path
	}

	// Preserve the existing query's encoding and repeated values. Insert new
	// parameters before the fragment, which is not part of the HTTP target.
	target, fragment, hasFragment := strings.Cut(path, "#")
	separator := "?"
	if _, query, hasQuery := strings.Cut(target, "?"); hasQuery {
		separator = "&"
		if query == "" || strings.HasSuffix(query, "&") {
			separator = ""
		}
	}
	target += separator + encoded
	if hasFragment {
		target += "#" + fragment
	}
	return target
}

// WithMaxRetries sets the ordinary retry limit (default 3). Zero disables retries.
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		if n < 0 || n > 2147483646 {
			c.optionErr = fmt.Errorf("max retries must be between 0 and 2147483646")
			return
		}
		c.maxRetries = n
	}
}

// WithRequestMaxRetries overrides the ordinary retry limit for one call.
// The single 401 token refresh remains independent of this limit.
func WithRequestMaxRetries(n int) RequestOption {
	return func(o *requestOptions) {
		if n < 0 || n > 2147483646 {
			o.optionErr = fmt.Errorf("max retries must be between 0 and 2147483646")
			return
		}
		o.maxRetries = &n
	}
}

type errReader struct{ err error }

func (e errReader) Read([]byte) (int, error) { return 0, e.err }

type hookRequestBody struct{ io.Reader }

func (hookRequestBody) Close() error { return nil }

var sdkOwnedHeaders = map[string]bool{
	"Authorization":      true,
	"Content-Type":       true,
	"Accept":             true,
	"Paxos-Sdk-Language": true,
	"Paxos-Sdk-Version":  true,
	"Paxos-Signature":    true,
}

func snapshotOwnedHeaders(h http.Header) map[string][]string {
	owned := make(map[string][]string)
	for key, values := range h {
		name := http.CanonicalHeaderKey(key)
		if sdkOwnedHeaders[name] && len(values) > 0 {
			owned[name] = append(owned[name], values...)
		}
	}
	for _, values := range owned {
		slices.Sort(values)
	}
	return owned
}

func ownedHeadersUnchanged(h http.Header, before map[string][]string) bool {
	after := snapshotOwnedHeaders(h)
	if len(after) != len(before) {
		return false
	}
	for name, values := range before {
		if !slices.Equal(values, after[name]) {
			return false
		}
	}
	return true
}
