// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// QuotesService provides access to the Quotes API.
//
// Quotes are "held rates" offered by Paxos to buy or sell assets at a specific
// price within a period of time - for example, the option to buy BTC within
// the next 30 seconds for $8,000.
//
// The typical Quotes workflow is as follows:
//
//  1. Call [QuotesService.ListQuotes] to get the latest available quoted
//     prices for the assets you want to buy or sell.
//  2. Present the offered price to one or more end-users, with a timer
//     indicating the time to expiration.
//  3. If a user accepts the price, call
//     [QuoteExecutionsService.CreateQuoteExecution] with the amount to buy or
//     sell.
//  4. Call [QuoteExecutionsService.GetQuoteExecution] to monitor for
//     completion and know when the funds from the execution are available.
//
// It's important to show end-users the latest available price. If you cache
// prices to show them to multiple users, you should refresh the cache once per
// second.
//
// You can call [QuoteExecutionsService.ListQuoteExecutions] to review or
// construct reports on quote execution activity.
type QuotesService struct {
	client *Client
}

// ListQuotes list Quotes
//
// Requires scope: exchange:read_quote
func (s *QuotesService) ListQuotes(ctx context.Context, params *ListQuotesParams, opts ...RequestOption) (*ListQuotesResponse, error) {
	return getWithParamsTyped[ListQuotesResponse](s.client, ctx, "/quotes", params, opts...)
}
