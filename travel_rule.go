// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// TravelRuleService provides access to the TravelRule API.
type TravelRuleService struct {
	client *Client
}

// ListVasps calls the ListVasps API operation.
func (s *TravelRuleService) ListVasps(ctx context.Context, params *ListVaspsParams, opts ...RequestOption) *Iterator[Vasp] {
	return paginateTyped[Vasp](s.client, ctx, "/travelrule/vasps", "items", params, opts...)
}

// GetVasp calls the GetVasp API operation.
func (s *TravelRuleService) GetVasp(ctx context.Context, id string, opts ...RequestOption) (*GetVaspResponse, error) {
	return getTyped[GetVaspResponse](s.client, ctx, buildPath("/travelrule/vasps/{id}", "id", id), opts...)
}
