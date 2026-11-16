// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// RewardsService provides access to the Rewards API.
//
// Rewards details including current reward rates, next accrual time, lifetime
// earned amounts, and estimated next accrual.
//
// *These endpoints are only relevant for whitelisted Rewards partners.*
type RewardsService struct {
	client *Client
}

// GetRewardDetails get reward details
//
// Requires scope: rewards:read_rewards_details
func (s *RewardsService) GetRewardDetails(ctx context.Context, params *GetRewardDetailsParams, opts ...RequestOption) (*GetRewardDetailsResponse, error) {
	return getWithParamsTyped[GetRewardDetailsResponse](s.client, ctx, "/rewards/details", params, opts...)
}

// GetRewardsRates get reward rates
//
// Requires scope: rewards:read_rewards_details
func (s *RewardsService) GetRewardsRates(ctx context.Context, params *GetRewardsRatesParams, opts ...RequestOption) (*GetRewardsRatesResponse, error) {
	return getWithParamsTyped[GetRewardsRatesResponse](s.client, ctx, "/rewards/rates", params, opts...)
}
