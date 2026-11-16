// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// PaymentsService provides access to the Payments API.
//
// Payments represent actual transfers of rewards to the payout address
// associated with a given statement. Payments are made on-chain to a specified
// payout address. Payments tie 1:1 with a statement amount.
//
// Use the `ref_id` to link payments to statements.
//
// > Some fields (e.g., `account_id`, `profile_id`) are populated for
// consistency but are not relevant to reward-specific flows.
//
// *These endpoints are only relevant for Global Dollar Network (GDN)
// partners.*
type PaymentsService struct {
	client *Client
}

// ListPayments list Payments
//
// Requires scope: statements:read_payment
//
// Deprecated: This API operation is deprecated.
func (s *PaymentsService) ListPayments(ctx context.Context, params *ListPaymentsParams, opts ...RequestOption) *Iterator[Payment] {
	return paginateTyped[Payment](s.client, ctx, "/statements/payments", "payments", params, opts...)
}
