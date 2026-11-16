// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// PricingService provides access to the Pricing API.
//
// Pricing provides historical data related to charting of asset prices.
//
// Note: Some Pricing API markets may be unavailable for trading.
type PricingService struct {
	client *Client
}

// ListPrices list Prices
func (s *PricingService) ListPrices(ctx context.Context, params *ListPricesParams, opts ...RequestOption) (*ListPricesResponse, error) {
	return getWithParamsTyped[ListPricesResponse](s.client, ctx, "/all-markets/prices", params, opts...)
}

// ListHistoricalPrices list Historical Prices
//
// Requires scope: exchange:historical_prices
func (s *PricingService) ListHistoricalPrices(ctx context.Context, market string, params *ListHistoricalPricesParams, opts ...RequestOption) (*ListHistoricalPricesResponse, error) {
	return getWithParamsTyped[ListHistoricalPricesResponse](s.client, ctx, buildPath("/markets/{market}/historical-prices", "market", market), params, opts...)
}

// ListHistoricalCandles list Historical Candles
//
// Requires scope: exchange:historical_prices
func (s *PricingService) ListHistoricalCandles(ctx context.Context, market string, params *ListHistoricalCandlesParams, opts ...RequestOption) *Iterator[CandlestickPrice] {
	return paginateTyped[CandlestickPrice](s.client, ctx, buildPath("/markets/{market}/historical-candles", "market", market), "items", params, opts...)
}
