// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// PaxosTransfersService provides access to the PaxosTransfers API.
//
// Move assets between two Entities belonging to the same Organization or to a
// different Organization on the Paxos platform.
//
// > Transferring USD between Entities is prohibited in some jurisdictions.
// Contact Support if you run into any problems with this restriction.
//
// Ensure the destination Profile already exists before beginning the transfer.
// This may require contacting someone outside your Organization to get the
// destination Profile ID.
type PaxosTransfersService struct {
	client *Client
}

// CreatePaxosTransfer create Paxos Transfer
//
// Requires scope: transfer:write_paxos_transfer
func (s *PaxosTransfersService) CreatePaxosTransfer(ctx context.Context, body *CreatePaxosTransferRequest, opts ...RequestOption) (*Transfer, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[Transfer](s.client, ctx, "/transfer/paxos", body, opts...)
}
