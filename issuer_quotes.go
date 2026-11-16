// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// IssuerQuotesService provides access to the IssuerQuotes API.
//
// Issuer Quotes are "held rates" offered by Paxos to mint or redeem Paxos-
// issued assets at a specific price within a period of time - for example, the
// option to buy PAXG within the next 5 seconds for $3,000.
//
// The typical Issuer Quotes workflow is as follows:
//
//  1. Call [IssuerQuotesService.CreateIssuerQuote] to get a quote for the
//     asset you want to mint or redeem to a specific profile.
//  2. Call [IssuerQuotesService.CreateIssuerQuoteExecution] with the quote ID
//     and profile ID.
//
// You can call [IssuerQuotesService.ListIssuerQuoteExecutions] to review or
// construct reports on issuer quote execution activity.
type IssuerQuotesService struct {
	client *Client
}

// CreateIssuerQuote create Issuer Quote
//
// Requires scope: exchange:write_issuer_quote
func (s *IssuerQuotesService) CreateIssuerQuote(ctx context.Context, body *CreateIssuerQuoteRequest, opts ...RequestOption) (*CreateIssuerQuoteResponse, error) {
	return postTyped[CreateIssuerQuoteResponse](s.client, ctx, "/issuer-quote", body, opts...)
}

// GetIssuerQuoteSettlementAvailability get Issuer Quote Settlement Availability
//
// Requires scope: exchange:write_issuer_quote
func (s *IssuerQuotesService) GetIssuerQuoteSettlementAvailability(ctx context.Context, params *GetIssuerQuoteSettlementAvailabilityParams, opts ...RequestOption) (*GetIssuerQuoteSettlementAvailabilityResponse, error) {
	return getWithParamsTyped[GetIssuerQuoteSettlementAvailabilityResponse](s.client, ctx, "/issuer-quote/settlement-availability", params, opts...)
}

// CreateIssuerQuoteExecution create Issuer Quote Execution
//
// Requires scope: exchange:write_issuer_quote_execution
func (s *IssuerQuotesService) CreateIssuerQuoteExecution(ctx context.Context, quoteID string, body *CreateIssuerQuoteExecutionRequest, opts ...RequestOption) (*CreateIssuerQuoteExecutionResponse, error) {
	return postTyped[CreateIssuerQuoteExecutionResponse](s.client, ctx, buildPath("/issuer-quote/{quote_id}", "quote_id", quoteID), body, opts...)
}

// ListIssuerQuoteExecutions list Issuer Quote Executions
//
// Requires scope: exchange:read_issuer_quote_execution
func (s *IssuerQuotesService) ListIssuerQuoteExecutions(ctx context.Context, params *ListIssuerQuoteExecutionsParams, opts ...RequestOption) *Iterator[IssuerQuoteExecution] {
	return paginateTyped[IssuerQuoteExecution](s.client, ctx, "/issuer-quotes", "issuer_quote_executions", params, opts...)
}
