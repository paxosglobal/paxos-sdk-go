// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// QuoteExecutionsService provides access to the QuoteExecutions API.
//
// Quote Executions buy or sell assets using a Quote obtained from the Quotes
// flow.
type QuoteExecutionsService struct {
	client *Client
}

// ListQuoteExecutions list Quote Executions
//
// Requires scope: exchange:read_quote_execution
func (s *QuoteExecutionsService) ListQuoteExecutions(ctx context.Context, profileID string, params *ListQuoteExecutionsParams, opts ...RequestOption) *Iterator[QuoteExecution] {
	return paginateTyped[QuoteExecution](s.client, ctx, buildPath("/profiles/{profile_id}/quote-executions", "profile_id", profileID), "items", params, opts...)
}

// CreateQuoteExecution create Quote Execution
//
// Requires scope: exchange:write_quote_execution
func (s *QuoteExecutionsService) CreateQuoteExecution(ctx context.Context, profileID string, body *CreateQuoteExecutionRequest, opts ...RequestOption) (*QuoteExecution, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[QuoteExecution](s.client, ctx, buildPath("/profiles/{profile_id}/quote-executions", "profile_id", profileID), body, opts...)
}

// GetQuoteExecution get Quote Execution
//
// Requires scope: exchange:read_quote_execution
func (s *QuoteExecutionsService) GetQuoteExecution(ctx context.Context, profileID string, id string, opts ...RequestOption) (*QuoteExecution, error) {
	return getTyped[QuoteExecution](s.client, ctx, buildPath("/profiles/{profile_id}/quote-executions/{id}", "profile_id", profileID, "id", id), opts...)
}
