// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// IdentityControlsService provides access to the IdentityControls API.
type IdentityControlsService struct {
	client *Client
}

// ListIdentityControls list Identity Controls
//
// Requires scope: identity:read_identity_control
func (s *IdentityControlsService) ListIdentityControls(ctx context.Context, params *ListIdentityControlsParams, opts ...RequestOption) *Iterator[IdentityControl] {
	return paginateTyped[IdentityControl](s.client, ctx, "/identity/controls", "items", params, opts...)
}

// DeleteIdentityControl delete Identity Control
//
// Requires scope: identity:write_identity_control
func (s *IdentityControlsService) DeleteIdentityControl(ctx context.Context, params *DeleteIdentityControlParams, opts ...RequestOption) (*DeleteIdentityControlResponse, error) {
	return deleteWithParamsTyped[DeleteIdentityControlResponse](s.client, ctx, "/identity/controls", params, opts...)
}

// CreateIdentityControl create Identity Control
//
// Requires scope: identity:write_identity_control
func (s *IdentityControlsService) CreateIdentityControl(ctx context.Context, body *CreateIdentityControlRequest, opts ...RequestOption) (*IdentityControl, error) {
	return postTyped[IdentityControl](s.client, ctx, "/identity/controls", body, opts...)
}
