// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// OrdersService provides access to the Orders API.
//
// There are four types of orders: market, limit, stop market, and stop limit.
//
// Market orders guarantee execution at a variable price and quantity. Limit
// orders guarantee price and quantity at the time of execution but do not
// guarantee execution. Stop-market orders convert to market orders when the
// stop price is reached. Stop-limit orders convert to limit orders at the
// specified price or better when the stop price is reached.
type OrdersService struct {
	client *Client
}

// ListExecutions list Executions
//
// Requires scope: exchange:read_order
func (s *OrdersService) ListExecutions(ctx context.Context, params *ListExecutionsParams, opts ...RequestOption) *Iterator[Execution] {
	return paginateTyped[Execution](s.client, ctx, "/executions", "items", params, opts...)
}

// ListOrders list Orders
//
// Requires scope: exchange:read_order
func (s *OrdersService) ListOrders(ctx context.Context, params *ListOrdersParams, opts ...RequestOption) *Iterator[Order] {
	return paginateTyped[Order](s.client, ctx, "/orders", "items", params, opts...)
}

// CreateOrder create Order
//
// Requires scope: exchange:write_order
func (s *OrdersService) CreateOrder(ctx context.Context, profileID string, body *CreateOrderRequest, opts ...RequestOption) (*Order, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Order](s.client, ctx, buildPath("/profiles/{profile_id}/orders", "profile_id", profileID), body, opts...)
}

// GetOrder get Order
//
// Requires scope: exchange:read_order
func (s *OrdersService) GetOrder(ctx context.Context, profileID string, id string, opts ...RequestOption) (*Order, error) {
	return getTyped[Order](s.client, ctx, buildPath("/profiles/{profile_id}/orders/{id}", "profile_id", profileID, "id", id), opts...)
}

// CancelOrder cancel Order
//
// Requires scope: exchange:write_order
func (s *OrdersService) CancelOrder(ctx context.Context, profileID string, id string, opts ...RequestOption) (*CancelOrderResponse, error) {
	return deleteTyped[CancelOrderResponse](s.client, ctx, buildPath("/profiles/{profile_id}/orders/{id}", "profile_id", profileID, "id", id), opts...)
}
