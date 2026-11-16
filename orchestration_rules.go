// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// OrchestrationRulesService provides access to the OrchestrationRules API.
type OrchestrationRulesService struct {
	client *Client
}

// ListOrchestrationRules list Orchestration Rules
//
// Requires scope: orchestration:read_orchestration_rule
func (s *OrchestrationRulesService) ListOrchestrationRules(ctx context.Context, params *ListOrchestrationRulesParams, opts ...RequestOption) *Iterator[OrchestrationRule] {
	return paginateTyped[OrchestrationRule](s.client, ctx, "/orchestration/rules", "items", params, opts...)
}

// CreateOrchestrationRule create Orchestration Rule
//
// Requires scope: orchestration:write_orchestration_rule
func (s *OrchestrationRulesService) CreateOrchestrationRule(ctx context.Context, body *CreateOrchestrationRuleRequest, opts ...RequestOption) (*CreateOrchestrationRuleResponse, error) {
	if body != nil && body.RefID == "" {
		bodyCopy := *body
		bodyCopy.RefID = uuid.NewString()
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[CreateOrchestrationRuleResponse](s.client, ctx, "/orchestration/rules", body, opts...)
}

// DeleteOrchestrationRule delete Orchestration Rule
//
// Requires scope: orchestration:write_orchestration_rule
func (s *OrchestrationRulesService) DeleteOrchestrationRule(ctx context.Context, id string, opts ...RequestOption) (*DeleteOrchestrationRuleResponse, error) {
	return deleteTyped[DeleteOrchestrationRuleResponse](s.client, ctx, buildPath("/orchestration/rules/{id}", "id", id), opts...)
}
