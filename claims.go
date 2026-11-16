// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// ClaimsService provides access to the Claims API.
//
// Claims represent rewards paid out on-chain. Use claim schedules to automate
// daily reward claims.
//
// *These endpoints are only relevant for whitelisted Rewards partners.*
type ClaimsService struct {
	client *Client
}

// ListClaimSchedules list claim schedules
//
// Requires scope: rewards:read_claims_schedule
func (s *ClaimsService) ListClaimSchedules(ctx context.Context, params *ListClaimSchedulesParams, opts ...RequestOption) *Iterator[ClaimSchedule] {
	return paginateTyped[ClaimSchedule](s.client, ctx, "/rewards/claim-schedules", "items", params, opts...)
}

// CreateClaimSchedule create a claim schedule
//
// Requires scope: rewards:write_claims_schedule
func (s *ClaimsService) CreateClaimSchedule(ctx context.Context, body *CreateClaimScheduleRequest, opts ...RequestOption) (*CreateClaimScheduleResponse, error) {
	return postTyped[CreateClaimScheduleResponse](s.client, ctx, "/rewards/claim-schedules", body, opts...)
}

// DeleteClaimSchedule delete a claim schedule
//
// Requires scope: rewards:write_claims_schedule
func (s *ClaimsService) DeleteClaimSchedule(ctx context.Context, id string, opts ...RequestOption) (*DeleteClaimScheduleResponse, error) {
	return deleteTyped[DeleteClaimScheduleResponse](s.client, ctx, buildPath("/rewards/claim-schedules/{id}", "id", id), opts...)
}

// UpdateClaimSchedule update a claim schedule
//
// Requires scope: rewards:write_claims_schedule
func (s *ClaimsService) UpdateClaimSchedule(ctx context.Context, id string, body *UpdateClaimScheduleRequest, opts ...RequestOption) (*UpdateClaimScheduleResponse, error) {
	return putTyped[UpdateClaimScheduleResponse](s.client, ctx, buildPath("/rewards/claim-schedules/{id}", "id", id), body, opts...)
}

// ListClaims list claims
//
// Requires scope: rewards:read_claims
func (s *ClaimsService) ListClaims(ctx context.Context, params *ListClaimsParams, opts ...RequestOption) *Iterator[Claim] {
	return paginateTyped[Claim](s.client, ctx, "/rewards/claims", "items", params, opts...)
}
