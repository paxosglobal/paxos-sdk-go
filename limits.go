// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// LimitsService provides access to the Limits API.
type LimitsService struct {
	client *Client
}

// ListTransferLimits list Transfer Limits
//
// Requires scope: transfer:read_transfer_limit
func (s *LimitsService) ListTransferLimits(ctx context.Context, params *ListTransferLimitsParams, opts ...RequestOption) (*ListTransferLimitsResponse, error) {
	return getWithParamsTyped[ListTransferLimitsResponse](s.client, ctx, "/transfer/limits/utilizations", params, opts...)
}
