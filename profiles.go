// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// ProfilesService provides access to the Profiles API.
//
// Profiles hold asset balances, and every Paxos transaction is on a particular
// profile. Depending on your integration type, Profiles may be used to reflect
// both corporate balances or individual user balances. There are two types of
// Profiles:
//
//   - `DEFAULT` Profiles are system-generated.
//   - `NORMAL` Profiles are created using the [ProfilesService.CreateProfile]
//     endpoint and are the only type of Profile you can create.
type ProfilesService struct {
	client *Client
}

// ListProfiles list Profiles
//
// Requires scope: funding:read_profile
func (s *ProfilesService) ListProfiles(ctx context.Context, params *ListProfilesParams, opts ...RequestOption) *Iterator[Profile] {
	return paginateTyped[Profile](s.client, ctx, "/profiles", "items", params, opts...)
}

// CreateProfile create Profile
//
// Requires scope: funding:write_profile
func (s *ProfilesService) CreateProfile(ctx context.Context, body *CreateProfileRequest, opts ...RequestOption) (*Profile, error) {
	return postTyped[Profile](s.client, ctx, "/profiles", body, opts...)
}

// GetProfile get Profile
//
// Requires scope: funding:read_profile
func (s *ProfilesService) GetProfile(ctx context.Context, profileID string, params *GetProfileParams, opts ...RequestOption) (*Profile, error) {
	return getWithParamsTyped[Profile](s.client, ctx, buildPath("/profiles/{profile_id}", "profile_id", profileID), params, opts...)
}

// UpdateProfile update Profile
//
// Requires scope: funding:write_profile
func (s *ProfilesService) UpdateProfile(ctx context.Context, profileID string, body *UpdateProfileRequest, opts ...RequestOption) (*Profile, error) {
	return putTyped[Profile](s.client, ctx, buildPath("/profiles/{profile_id}", "profile_id", profileID), body, opts...)
}

// ListProfileBalances list Profile Balances
//
// Requires scope: funding:read_profile
func (s *ProfilesService) ListProfileBalances(ctx context.Context, profileID string, params *ListProfileBalancesParams, opts ...RequestOption) (*ListProfileBalancesResponse, error) {
	return getWithParamsTyped[ListProfileBalancesResponse](s.client, ctx, buildPath("/profiles/{profile_id}/balances", "profile_id", profileID), params, opts...)
}

// GetProfileBalance get Profile Balance
//
// Requires scope: funding:read_profile
func (s *ProfilesService) GetProfileBalance(ctx context.Context, profileID string, asset string, params *GetProfileBalanceParams, opts ...RequestOption) (*ProfileBalance, error) {
	return getWithParamsTyped[ProfileBalance](s.client, ctx, buildPath("/profiles/{profile_id}/balances/{asset}", "profile_id", profileID, "asset", asset), params, opts...)
}

// DeactivateProfile deactivate Profile
//
// Requires scope: funding:write_profile
func (s *ProfilesService) DeactivateProfile(ctx context.Context, profileID string, opts ...RequestOption) (*DeactivateProfileResponse, error) {
	return putTyped[DeactivateProfileResponse](s.client, ctx, buildPath("/profiles/{profile_id}/deactivate", "profile_id", profileID), nil, opts...)
}
