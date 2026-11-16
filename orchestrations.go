// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// OrchestrationsService provides access to the Orchestrations API.
type OrchestrationsService struct {
	client *Client
}

// ListOrchestrations list Orchestrations
//
// Requires scope: orchestration:read_orchestration
func (s *OrchestrationsService) ListOrchestrations(ctx context.Context, params *ListOrchestrationsParams, opts ...RequestOption) *Iterator[Orchestration] {
	return paginateTyped[Orchestration](s.client, ctx, "/orchestration/orchestrations", "items", params, opts...)
}

// CreateOrchestration create Orchestration
//
// Requires scope: orchestration:write_orchestration
func (s *OrchestrationsService) CreateOrchestration(ctx context.Context, body *CreateOrchestrationRequest, opts ...RequestOption) (*Orchestration, error) {
	if body != nil && body.RefID == "" {
		bodyCopy := *body
		bodyCopy.RefID = uuid.NewString()
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Orchestration](s.client, ctx, "/orchestration/orchestrations", body, opts...)
}

// GetOrchestration get Orchestration
//
// Requires scope: orchestration:read_orchestration
func (s *OrchestrationsService) GetOrchestration(ctx context.Context, id string, opts ...RequestOption) (*Orchestration, error) {
	return getTyped[Orchestration](s.client, ctx, buildPath("/orchestration/orchestrations/{id}", "id", id), opts...)
}
