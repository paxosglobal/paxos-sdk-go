// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// RewardsStatementsService provides access to the RewardsStatements API.
type RewardsStatementsService struct {
	client *Client
}

// ListRewardsStatements list Rewards Statements
//
// Requires scope: rewards:read_statements
func (s *RewardsStatementsService) ListRewardsStatements(ctx context.Context, params *ListRewardsStatementsParams, opts ...RequestOption) *Iterator[RewardsStatement] {
	return paginateTyped[RewardsStatement](s.client, ctx, "/rewards/statements", "items", params, opts...)
}

// GetRewardsStatement get Rewards Statement
//
// Requires scope: rewards:read_statements
func (s *RewardsStatementsService) GetRewardsStatement(ctx context.Context, id string, opts ...RequestOption) (*GetRewardsStatementResponse, error) {
	return getTyped[GetRewardsStatementResponse](s.client, ctx, buildPath("/rewards/statements/{id}", "id", id), opts...)
}
