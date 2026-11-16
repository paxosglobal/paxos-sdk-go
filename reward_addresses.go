// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// RewardAddressesService provides access to the RewardAddresses API.
//
// Reward addresses are blockchain addresses registered for reward attribution.
//
// Custody Rewards: For EVM chains, Paxos registers these addresses on-chain
// for reward calculation and accrual in the token contract. For Solana,
// provide the owner account (not a token account) — Paxos will derive and
// register the associated token account on-chain for the given asset.
// Providing a token account will result in a derived address with no meaning.
// For other non-EVM chains, Paxos tracks on-chain activity for custody reward
// calculation.
//
// Non-Custody Rewards: Paxos tracks on-chain activity for these addresses to
// calculate and distribute non-custody rewards to your organization.
//
// *These endpoints are only relevant for whitelisted Rewards partners.*
type RewardAddressesService struct {
	client *Client
}

// ListAttributionAddresses list reward addresses
//
// Requires scope: rewards:read_monitoring_address
func (s *RewardAddressesService) ListAttributionAddresses(ctx context.Context, params *ListAttributionAddressesParams, opts ...RequestOption) *Iterator[AttributionAddress] {
	return paginateTyped[AttributionAddress](s.client, ctx, "/rewards/addresses", "items", params, opts...)
}

// CreateAttributionAddress create a reward address
//
// Requires scope: rewards:write_monitoring_address
func (s *RewardAddressesService) CreateAttributionAddress(ctx context.Context, body *CreateAttributionAddressRequest, opts ...RequestOption) (*CreateAttributionAddressResponse, error) {
	return postTyped[CreateAttributionAddressResponse](s.client, ctx, "/rewards/addresses", body, opts...)
}

// BatchCreateAttributionAddresses batch create reward addresses
//
// Requires scope: rewards:write_monitoring_address
func (s *RewardAddressesService) BatchCreateAttributionAddresses(ctx context.Context, body *BatchCreateAttributionAddressesRequest, opts ...RequestOption) (*BatchCreateAttributionAddressesResponse, error) {
	return postTyped[BatchCreateAttributionAddressesResponse](s.client, ctx, "/rewards/addresses/batch", body, opts...)
}

// DeleteAttributionAddress delete a reward address
//
// Requires scope: rewards:write_monitoring_address
func (s *RewardAddressesService) DeleteAttributionAddress(ctx context.Context, id string, opts ...RequestOption) (*DeleteAttributionAddressResponse, error) {
	return deleteTyped[DeleteAttributionAddressResponse](s.client, ctx, buildPath("/rewards/addresses/{id}", "id", id), opts...)
}

// UpdateAttributionAddress update a reward address
//
// Requires scope: rewards:write_monitoring_address
func (s *RewardAddressesService) UpdateAttributionAddress(ctx context.Context, id string, body *AttributionAddress, opts ...RequestOption) (*UpdateAttributionAddressResponse, error) {
	return putTyped[UpdateAttributionAddressResponse](s.client, ctx, buildPath("/rewards/addresses/{id}", "id", id), body, opts...)
}
