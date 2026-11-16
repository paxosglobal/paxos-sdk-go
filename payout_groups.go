// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// PayoutGroupsService provides access to the PayoutGroups API.
//
// Payout groups encapsulate common rewards characteristics. All monitored
// addresses that reference a payout group accumulate rewards at the same rate
// and claim to the payout destination specified for that group.
//
// *These endpoints are only relevant for whitelisted Rewards partners.*
type PayoutGroupsService struct {
	client *Client
}

// ListPayoutGroups list Payout Groups
//
// Requires scope: rewards:read_payout_groups
func (s *PayoutGroupsService) ListPayoutGroups(ctx context.Context, params *ListPayoutGroupsParams, opts ...RequestOption) *Iterator[PayoutGroup] {
	return paginateTyped[PayoutGroup](s.client, ctx, "/rewards/payout-groups", "items", params, opts...)
}

// CreatePayoutGroup create Payout Group
//
// Requires one of: rewards:create_payout_groups; rewards:write_payout_groups
func (s *PayoutGroupsService) CreatePayoutGroup(ctx context.Context, body *CreatePayoutGroupRequest, opts ...RequestOption) (*CreatePayoutGroupResponse, error) {
	return postTyped[CreatePayoutGroupResponse](s.client, ctx, "/rewards/payout-groups", body, opts...)
}

// DeletePayoutGroup delete Payout Group
//
// Requires one of: rewards:delete_payout_groups; rewards:write_payout_groups
func (s *PayoutGroupsService) DeletePayoutGroup(ctx context.Context, id string, opts ...RequestOption) (*DeletePayoutGroupResponse, error) {
	return deleteTyped[DeletePayoutGroupResponse](s.client, ctx, buildPath("/rewards/payout-groups/{id}", "id", id), opts...)
}

// UpdatePayoutGroup update Payout Group
//
// Requires one of: rewards:update_payout_groups; rewards:write_payout_groups
func (s *PayoutGroupsService) UpdatePayoutGroup(ctx context.Context, id string, body *UpdatePayoutGroupRequest, opts ...RequestOption) (*UpdatePayoutGroupResponse, error) {
	return putTyped[UpdatePayoutGroupResponse](s.client, ctx, buildPath("/rewards/payout-groups/{id}", "id", id), body, opts...)
}
