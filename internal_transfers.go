// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// InternalTransfersService provides access to the InternalTransfers API.
type InternalTransfersService struct {
	client *Client
}

// CreateInternalTransfer create Internal Transfer
//
// Requires scope: transfer:write_internal_transfer
func (s *InternalTransfersService) CreateInternalTransfer(ctx context.Context, body *CreateInternalTransferRequest, opts ...RequestOption) (*Transfer, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Transfer](s.client, ctx, "/transfer/internal", body, opts...)
}
