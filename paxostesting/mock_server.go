// Package paxostesting provides test utilities for the Paxos SDK.
//
// It includes a configurable mock HTTP server, pre-built response fixtures,
// and stub helpers for common error scenarios. Use it in tests to verify
// SDK integration without making real API calls.
//
// Example:
//
//	func TestMyFeature(t *testing.T) {
//	    server := paxostesting.NewMockServer(t)
//
//	    server.On("GET", "/transfer/transfers/txn_123").
//	        Respond(200, paxostesting.TransferFixture("txn_123"))
//
//	    client, _ := paxos.NewClient("id", "secret",
//	        paxos.WithScopes(paxos.ScopeTransferReadTransfer),
//	        server.ClientOption(),
//	    )
//
//	    transfer, err := client.Transfers.GetTransfer(ctx, "txn_123")
//	    assert.NoError(t, err)
//	    assert.Equal(t, "txn_123", transfer.ID)
//	}
package paxostesting

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	paxos "github.com/paxosglobal/paxos-sdk-go"
)

const tokenPath = "/oauth2/token"

// MockServer is an in-process HTTP server that records incoming requests
// and returns pre-configured responses. It wraps [httptest.Server] and
// provides a fluent API for registering route expectations.
type MockServer struct {
	t      testing.TB
	server *httptest.Server

	mu     sync.Mutex
	routes map[routeKey]*Route
	calls  []RecordedRequest
}

// routeKey uniquely identifies a registered mock route.
type routeKey struct {
	Method string
	Path   string
}

// RecordedRequest stores details about an HTTP request received by the mock server.
type RecordedRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   []byte
}

// NewMockServer creates a new mock HTTP server that automatically shuts down
// when the test completes. All unmatched requests return 404.
func NewMockServer(t testing.TB) *MockServer {
	m := &MockServer{
		t:      t,
		routes: make(map[routeKey]*Route),
	}

	m.server = httptest.NewServer(http.HandlerFunc(m.handler))
	t.Cleanup(m.server.Close)

	return m
}

// URL returns the base URL of the mock server (e.g. "http://127.0.0.1:12345").
func (m *MockServer) URL() string {
	return m.server.URL
}

// TokenURL returns the URL of the mock OAuth token endpoint.
func (m *MockServer) TokenURL() string {
	return m.server.URL + tokenPath
}

// ClientOption returns an option that points a client at this server for
// both API requests and OAuth token requests.
func (m *MockServer) ClientOption() paxos.Option {
	return func(c *paxos.Client) {
		paxos.WithBaseURL(m.URL())(c)
		paxos.WithTokenURL(m.TokenURL())(c)
	}
}

// On registers a route expectation for the given HTTP method and path.
// Returns a [Route] that can be configured with a response.
//
// Example:
//
//	server.On("GET", "/transfer/transfers/txn_123").
//	    Respond(200, paxostesting.TransferFixture("txn_123"))
func (m *MockServer) On(method, path string) *Route {
	m.mu.Lock()
	defer m.mu.Unlock()

	r := &Route{}
	m.routes[routeKey{Method: method, Path: path}] = r
	return r
}

// Calls returns all recorded requests received by the mock server.
func (m *MockServer) Calls() []RecordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]RecordedRequest, len(m.calls))
	copy(out, m.calls)
	return out
}

// Reset clears all registered routes and recorded requests.
func (m *MockServer) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.routes = make(map[routeKey]*Route)
	m.calls = nil
}

// ReturnNotFound is a convenience helper that registers a 404 response
// for the given method and path.
func (m *MockServer) ReturnNotFound(method, path string) {
	m.On(method, path).Respond(http.StatusNotFound, map[string]interface{}{
		"type":       "https://docs.paxos.com/problems/not-found",
		"title":      "Not Found",
		"status":     404,
		"detail":     "The requested resource was not found",
		"request_id": "req_mock_not_found",
	})
}

// ReturnRateLimit is a convenience helper that registers a 429 response. The
// client retries a 429 with backoff, so a test that only checks the error should
// build its client with paxos.WithMaxRetries(0) to avoid waiting several seconds.
func (m *MockServer) ReturnRateLimit(method, path string) {
	m.On(method, path).
		Respond(http.StatusTooManyRequests, map[string]interface{}{
			"type":       "https://docs.paxos.com/problems/rate-limit",
			"title":      "Too Many Requests",
			"status":     429,
			"detail":     "Rate limit exceeded",
			"request_id": "req_mock_rate_limit",
		})
}

// ReturnServerError is a convenience helper that registers a 500 response
// for the given method and path.
func (m *MockServer) ReturnServerError(method, path string) {
	m.On(method, path).Respond(http.StatusInternalServerError, map[string]interface{}{
		"type":       "https://docs.paxos.com/problems/internal-error",
		"title":      "Internal Server Error",
		"status":     500,
		"detail":     "An unexpected error occurred",
		"request_id": "req_mock_server_error",
	})
}

// ReturnValidationError is a convenience helper that registers a 400 response
// for the given method and path.
func (m *MockServer) ReturnValidationError(method, path, detail string) {
	m.On(method, path).Respond(http.StatusBadRequest, map[string]interface{}{
		"type":       "https://docs.paxos.com/problems/validation-error",
		"title":      "Bad Request",
		"status":     400,
		"detail":     detail,
		"request_id": "req_mock_validation",
	})
}

// handler dispatches incoming requests to registered routes.
func (m *MockServer) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.URL.Path == tokenPath {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "mock-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
		return
	}

	// Read and record the request body.
	var body []byte
	if r.Body != nil {
		body, _ = readAll(r.Body)
		r.Body.Close()
	}

	m.mu.Lock()
	m.calls = append(m.calls, RecordedRequest{
		Method: r.Method,
		Path:   r.URL.Path,
		Header: r.Header.Clone(),
		Body:   body,
	})

	route, ok := m.routes[routeKey{Method: r.Method, Path: r.URL.Path}]
	m.mu.Unlock()

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"type":   "https://docs.paxos.com/problems/not-found",
			"title":  "Not Found",
			"status": 404,
			"detail": "No mock route registered for " + r.Method + " " + r.URL.Path,
		})
		return
	}

	route.mu.Lock()
	resp := route.nextResponse()
	route.mu.Unlock()

	// Write response headers.
	w.Header().Set("Content-Type", "application/json")
	if id := r.Header.Get("Paxos-Request-Id"); id != "" {
		w.Header().Set("Paxos-Request-Id", id)
	}
	for k, v := range resp.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(resp.Status)

	if resp.Body != nil {
		json.NewEncoder(w).Encode(resp.Body)
	}
}

// Route represents a registered mock route with one or more queued responses.
type Route struct {
	mu        sync.Mutex
	responses []routeResponse
	index     int
}

// routeResponse holds a single response configuration.
type routeResponse struct {
	Status  int
	Body    interface{}
	Headers map[string]string
}

// Respond configures the route to return the given status code and JSON body.
// If called multiple times, responses are returned in order. After all
// queued responses have been consumed, the last response is repeated.
func (r *Route) Respond(status int, body interface{}) *Route {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.responses = append(r.responses, routeResponse{
		Status: status,
		Body:   body,
	})
	return r
}

// RespondWithHeaders configures the route to return the given status code,
// JSON body, and additional HTTP headers.
func (r *Route) RespondWithHeaders(status int, body interface{}, headers map[string]string) *Route {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.responses = append(r.responses, routeResponse{
		Status:  status,
		Body:    body,
		Headers: headers,
	})
	return r
}

// nextResponse returns the next queued response. If all responses have been
// consumed, the last one is repeated. Must be called with r.mu held.
func (r *Route) nextResponse() routeResponse {
	if len(r.responses) == 0 {
		return routeResponse{Status: http.StatusOK}
	}
	resp := r.responses[r.index]
	if r.index < len(r.responses)-1 {
		r.index++
	}
	return resp
}
