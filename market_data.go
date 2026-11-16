// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// MarketDataService provides access to the MarketData API.
//
// Market data provides various parameters of the order book and historical
// order data.
type MarketDataService struct {
	client *Client
}

// ListMarkets list Markets
func (s *MarketDataService) ListMarkets(ctx context.Context, opts ...RequestOption) (*ListMarketsResponse, error) {
	return getTyped[ListMarketsResponse](s.client, ctx, "/markets", opts...)
}
