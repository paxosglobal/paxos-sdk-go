// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// ExchangeService provides access to the Exchange API.
type ExchangeService struct {
	client *Client
}

// ListProfileExecutions list Executions
//
// Requires scope: exchange:read_order
func (s *ExchangeService) ListProfileExecutions(ctx context.Context, profileID string, params *ListProfileExecutionsParams, opts ...RequestOption) *Iterator[Execution] {
	return paginateTyped[Execution](s.client, ctx, buildPath("/profiles/{profile_id}/executions", "profile_id", profileID), "items", params, opts...)
}

// ListProfileOrders list Orders
//
// Requires scope: exchange:read_order
func (s *ExchangeService) ListProfileOrders(ctx context.Context, profileID string, params *ListProfileOrdersParams, opts ...RequestOption) *Iterator[Order] {
	return paginateTyped[Order](s.client, ctx, buildPath("/profiles/{profile_id}/orders", "profile_id", profileID), "items", params, opts...)
}
