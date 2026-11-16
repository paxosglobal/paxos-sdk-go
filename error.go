package paxos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// Error represents an RFC 7807 Problem Details error from the Paxos API. An
// *AuthenticationError also matches it through errors.As.
type Error struct {
	// Type is a URI reference identifying the problem type.
	Type string `json:"type,omitempty"`

	// Title is a short, human-readable summary of the problem type.
	Title string `json:"title,omitempty"`

	// Status is the HTTP status code.
	Status int `json:"status,omitempty"`

	// Detail is a human-readable explanation specific to this occurrence.
	Detail string `json:"detail,omitempty"`

	// Instance is the URI reference identifying the specific occurrence of the problem.
	Instance string `json:"instance,omitempty"`

	// Meta contains additional metadata about the error.
	Meta json.RawMessage `json:"meta,omitempty"`

	// RequestID is the ID of the request that caused the error.
	RequestID string `json:"request_id,omitempty"`

	cause error
}

type sanitizedDecodeError struct{ kind string }

func (e sanitizedDecodeError) Error() string { return "JSON parse error (" + e.kind + ")" }

func newMalformedResponseError(status int, requestID string, decodeErr error) *Error {
	return &Error{
		Title:     "Malformed response",
		Status:    status,
		Detail:    "failed to decode response body",
		RequestID: requestID,
		cause:     sanitizedDecodeError{kind: fmt.Sprintf("%T", decodeErr)},
	}
}

// Unwrap returns the underlying cause, if any.
func (e *Error) Unwrap() error { return e.cause }

// Error implements the error interface.
func (e *Error) Error() string {
	title := e.Title
	if title == "" {
		title = "API Error"
	}
	parts := fmt.Sprintf("status %d", e.Status)
	if e.RequestID != "" {
		parts += ", request_id: " + e.RequestID
	}
	if e.Instance != "" {
		parts += ", instance: " + e.Instance
	}
	if e.Detail != "" {
		return fmt.Sprintf("%s (%s): %s", title, parts, e.Detail)
	}
	return fmt.Sprintf("%s (%s)", title, parts)
}

// IsNotFound returns true if the error (or any wrapped error) is a 404 Not Found error.
// This and the other Is* status helpers match *Error nodes anywhere in the error
// tree, including wrapped and joined errors, and skip *AuthenticationError. They
// walk Unwrap methods directly and do not consult custom As methods on caller
// wrapper types.
func IsNotFound(err error) bool {
	e := apiError(err)
	return e != nil && e.Status == http.StatusNotFound
}

// IsRateLimited returns true if the error (or any wrapped error) is a 429 Too Many Requests error.
func IsRateLimited(err error) bool {
	e := apiError(err)
	return e != nil && e.Status == http.StatusTooManyRequests
}

// IsUnauthorized returns true if the error (or any wrapped error) is a 401 Unauthorized error.
func IsUnauthorized(err error) bool {
	e := apiError(err)
	return e != nil && e.Status == http.StatusUnauthorized
}

// IsForbidden returns true if the error (or any wrapped error) is a 403 Forbidden error.
func IsForbidden(err error) bool {
	e := apiError(err)
	return e != nil && e.Status == http.StatusForbidden
}

// IsValidationError returns true if the error (or any wrapped error) is a 400 Bad Request error.
func IsValidationError(err error) bool {
	e := apiError(err)
	return e != nil && e.Status == http.StatusBadRequest
}

// IsConflict returns true if the error (or any wrapped error) is a 409 Conflict error.
func IsConflict(err error) bool {
	e := apiError(err)
	return e != nil && e.Status == http.StatusConflict
}

// IsUnprocessable returns true if the error (or any wrapped error) is a 422 Unprocessable Entity error.
func IsUnprocessable(err error) bool {
	e := apiError(err)
	return e != nil && e.Status == http.StatusUnprocessableEntity
}

// IsServerError returns true if the error (or any wrapped error) is a 5xx server error.
func IsServerError(err error) bool {
	e := apiError(err)
	return e != nil && e.Status >= 500
}

func apiError(err error) *Error {
	switch e := err.(type) {
	case nil, *AuthenticationError:
		return nil
	case *Error:
		return e
	case interface{ Unwrap() error }:
		return apiError(e.Unwrap())
	case interface{ Unwrap() []error }:
		for _, child := range e.Unwrap() {
			if found := apiError(child); found != nil {
				return found
			}
		}
	}
	return nil
}

// NetworkError reports a transport failure (DNS, connection refused, TLS,
// reset) where no HTTP response was received. It is not an *Error.
type NetworkError struct {
	// Detail is a human-readable explanation including the underlying error text.
	Detail string

	// RequestID is the ID of the request that failed, if known.
	RequestID string

	cause error
}

// Error implements the error interface.
func (e *NetworkError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("Network error (request_id: %s): %s", e.RequestID, e.Detail)
	}
	return "Network error: " + e.Detail
}

// Unwrap returns the underlying transport error.
func (e *NetworkError) Unwrap() error { return e.cause }

// IsNetworkError returns true if the error (or any wrapped error) is a *NetworkError.
func IsNetworkError(err error) bool {
	var e *NetworkError
	return errors.As(err, &e)
}

// TimeoutError reports an SDK call, token or pagination deadline that expired.
// It is not an *Error and wraps context.DeadlineExceeded.
type TimeoutError struct {
	// Detail is a human-readable explanation including the underlying error text.
	Detail string

	// RequestID is the Paxos-Request-Id of the most recent attempt built, or empty if
	// no attempt was built.
	RequestID string

	cause error
}

// Error implements the error interface.
func (e *TimeoutError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("Timeout error (request_id: %s): %s", e.RequestID, e.Detail)
	}
	return "Timeout error: " + e.Detail
}

// Unwrap returns the underlying deadline error.
func (e *TimeoutError) Unwrap() error { return e.cause }

// CanceledError reports a call or pagination cancelled by the caller's context.
// It is not an *Error and wraps context.Canceled.
type CanceledError struct {
	// Detail is a human-readable explanation including the underlying error text.
	Detail string

	// RequestID is the Paxos-Request-Id of the most recent attempt built, or empty if
	// no attempt was built. A cancelled mutation may still have reached Paxos, so quote
	// the ID to reconcile.
	RequestID string

	cause error
}

// Error implements the error interface.
func (e *CanceledError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("Canceled (request_id: %s): %s", e.RequestID, e.Detail)
	}
	return "Canceled: " + e.Detail
}

// Unwrap returns the underlying cancellation error.
func (e *CanceledError) Unwrap() error { return e.cause }

func withRequestID(err error, requestID string) error {
	var ne *NetworkError
	var te *TimeoutError
	var ce *CanceledError
	if errors.As(err, &ne) {
		c := *ne
		c.RequestID = requestID
		return &c
	}
	if errors.As(err, &te) {
		c := *te
		c.RequestID = requestID
		return &c
	}
	if errors.As(err, &ce) {
		c := *ce
		c.RequestID = requestID
		return &c
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return &TimeoutError{Detail: err.Error(), RequestID: requestID, cause: err}
	}
	if errors.Is(err, context.Canceled) {
		return &CanceledError{Detail: err.Error(), RequestID: requestID, cause: err}
	}
	return err
}

func newNetworkError(prefix string, cause error, requestID string, target *url.URL) error {
	if errors.Is(cause, context.DeadlineExceeded) {
		return &TimeoutError{Detail: prefix + ": " + cause.Error(), RequestID: requestID, cause: cause}
	}
	return &NetworkError{Detail: prefix + ": " + describeCause(cause, target), RequestID: requestID, cause: cause}
}

func finalURL(resp *http.Response, req *http.Request) *url.URL {
	if resp != nil && resp.Request != nil && resp.Request.URL != nil {
		return resp.Request.URL
	}
	return req.URL
}

func describeCause(cause error, target *url.URL) string {
	text := cause.Error()
	root := cause
	for next := errors.Unwrap(root); next != nil; next = errors.Unwrap(root) {
		root = next
	}
	if root.Error() == "" {
		text = strings.TrimSpace(strings.TrimRight(text, " ") + " " + fmt.Sprintf("%T", root))
	}
	if target == nil || target.Hostname() == "" || strings.Contains(text, target.Hostname()) {
		return text
	}
	port := target.Port()
	if port == "" {
		port = "80"
		if target.Scheme == "https" {
			port = "443"
		}
	}
	return text + ": " + net.JoinHostPort(target.Hostname(), port)
}

// RequestID returns the Paxos-Request-Id of any SDK error, or empty.
func RequestID(err error) string {
	var e *Error
	var ne *NetworkError
	var te *TimeoutError
	var ce *CanceledError
	var ae *AuthenticationError
	switch {
	case errors.As(err, &e):
		return e.RequestID
	case errors.As(err, &ne):
		return ne.RequestID
	case errors.As(err, &te):
		return te.RequestID
	case errors.As(err, &ce):
		return ce.RequestID
	case errors.As(err, &ae):
		return ae.RequestID
	}
	return ""
}

// StatusCode returns the HTTP status of an API error, or 0 when err has none.
func StatusCode(err error) int {
	var e *Error
	var ae *AuthenticationError
	switch {
	case errors.As(err, &e):
		return e.Status
	case errors.As(err, &ae):
		return ae.Status
	}
	return 0
}

// Detail returns the detail of any SDK error (the OAuth code for authentication errors), or empty.
func Detail(err error) string {
	var e *Error
	var ne *NetworkError
	var te *TimeoutError
	var ce *CanceledError
	var ae *AuthenticationError
	switch {
	case errors.As(err, &e):
		return e.Detail
	case errors.As(err, &ne):
		return ne.Detail
	case errors.As(err, &te):
		return te.Detail
	case errors.As(err, &ce):
		return ce.Detail
	case errors.As(err, &ae):
		return ae.Code
	}
	return ""
}

// Title returns the title of an API error, "OAuth authentication failed" for
// an authentication error, or empty for network and timeout errors, other
// errors, and API errors with no title.
func Title(err error) string {
	var e *Error
	var ae *AuthenticationError
	switch {
	case errors.As(err, &e):
		return e.Title
	case errors.As(err, &ae):
		return "OAuth authentication failed"
	}
	return ""
}

// IsTimeoutError returns true if the error (or any wrapped error) is a
// request or pagination timeout — either a per-request timeout
// (WithCallTimeout) or an auto-paginating Iterator's WithMaxDuration
// deadline, both of which surface as an error wrapping context.DeadlineExceeded.
func IsTimeoutError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}

// IsCanceledError returns true if the error (or any wrapped error) is a call or
// pagination cancelled by the caller's context.
func IsCanceledError(err error) bool {
	return errors.Is(err, context.Canceled)
}

// parseError attempts to parse an HTTP response as an Error.
func parseError(resp *http.Response, body []byte, requestID string) *Error {
	var apiErr Error
	var fields map[string]json.RawMessage
	if json.Unmarshal(body, &fields) != nil || fields == nil {
		apiErr.Title = resp.Status
	} else {
		// Decode fields independently so one malformed field cannot hide useful
		// diagnostics or replace the HTTP status with a parser error.
		for key, dest := range map[string]*string{
			"type": &apiErr.Type, "title": &apiErr.Title, "detail": &apiErr.Detail,
			"instance": &apiErr.Instance, "request_id": &apiErr.RequestID,
		} {
			_ = json.Unmarshal(fields[key], dest)
		}
		apiErr.Meta = fields["meta"]
	}

	// The transport status is authoritative, even when the body disagrees.
	apiErr.Status = resp.StatusCode
	if apiErr.RequestID == "" {
		apiErr.RequestID = requestID
	}

	return &apiErr
}

// AuthenticationError reports an OAuth HTTP failure or invalid token response.
// It is distinct from an API endpoint's 401 Error, but like the other SDKs it is
// also an API error: errors.As with a *Error target succeeds and yields the
// status, request ID, the title "OAuth authentication failed" and the OAuth code
// as detail. The IsNotFound, IsUnauthorized and other status helpers do not
// match it; use IsAuthenticationError. Transport failures reaching the token
// endpoint are a *NetworkError or *TimeoutError, not an AuthenticationError.
// Diagnostics never include the raw token response, error_description, or
// credentials.
type AuthenticationError struct {
	Status    int
	RequestID string
	// Code is a recognized OAuth error code, or empty.
	Code string
}

func (e *AuthenticationError) Error() string {
	message := fmt.Sprintf("OAuth authentication failed (status %d", e.Status)
	if e.RequestID != "" {
		message += ", request_id: " + e.RequestID
	}
	message += ")"
	if e.Code != "" {
		message += ": " + e.Code
	}
	return message
}

// As lets errors.As match an *AuthenticationError with a *Error target.
func (e *AuthenticationError) As(target any) bool {
	t, ok := target.(**Error)
	if !ok || e == nil {
		return false
	}
	*t = &Error{Title: "OAuth authentication failed", Status: e.Status, Detail: e.Code, RequestID: e.RequestID}
	return true
}

// IsAuthenticationError reports whether err wraps an OAuth AuthenticationError.
func IsAuthenticationError(err error) bool {
	var e *AuthenticationError
	return errors.As(err, &e)
}

func newAuthenticationError(resp *http.Response, body []byte) *AuthenticationError {
	e := &AuthenticationError{Status: resp.StatusCode, RequestID: resp.Header.Get(headerRequestID)}
	var fields map[string]json.RawMessage
	if json.Unmarshal(body, &fields) == nil {
		var requestID string
		if json.Unmarshal(fields["request_id"], &requestID) == nil && requestID != "" {
			e.RequestID = requestID
		}
	}
	var payload struct {
		Code string `json:"error"`
	}
	if json.Unmarshal(body, &payload) == nil {
		switch payload.Code {
		case "invalid_request", "invalid_client", "invalid_grant", "unauthorized_client",
			"unsupported_grant_type", "invalid_scope", "server_error", "temporarily_unavailable":
			e.Code = payload.Code
		}
	}
	return e
}
