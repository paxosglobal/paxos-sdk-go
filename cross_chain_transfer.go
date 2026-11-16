// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// CrossChainTransferService provides access to the CrossChainTransfer API.
//
// Look up the status of a cross-chain transfer after burning on the source
// network, and the fees that apply to a destination network and asset. Each
// burn is identified by its receipt. Use
// [CrossChainTransferService.GetCrossChainAttestation] to look up a single
// burn by receipt, or [CrossChainTransferService.ListCrossChainAttestations]
// to look up the burns in a source transaction by `source_network` and
// `tx_hash`. [CrossChainTransferService.ListCrossChainTransfers] lists the
// burns a wallet submitted, newest first.
// [CrossChainTransferService.GetCrossChainFees] returns the `transfer_fee` and
// `forwarding_fee` for a `destination_network` and `asset`.
//
// These endpoints do not require authentication. Responses omit fields at
// their default value (`0`, `false`, `""`, empty list). They are rate limited
// per IP; on `429`, back off before retrying, for the `Retry-After` seconds
// when that header is present.
//
// A burn is `PENDING` until it is verified and signed, then `SIGNED` with an
// `authorization` to mint on the destination network. `COMPLETED`, `EXPIRED`
// and `FAILED` are terminal states, so stop polling when a burn reaches one of
// them. An `EXPIRED` burn was not redeemed before its authorization's
// `expires_at`. A receipt that has not been observed yet returns `404`, and a
// transaction that has not been observed yet returns no `items`, so retry for
// a bounded time before treating either as an error.
type CrossChainTransferService struct {
	client *Client
}

// ListCrossChainAttestations list Cross-Chain Attestations
func (s *CrossChainTransferService) ListCrossChainAttestations(ctx context.Context, params *ListCrossChainAttestationsParams, opts ...RequestOption) (*ListCrossChainAttestationsResponse, error) {
	return getWithParamsTyped[ListCrossChainAttestationsResponse](s.client, ctx, "/cross-chain/attestations", params, opts...)
}

// GetCrossChainAttestation get Cross-Chain Attestation
func (s *CrossChainTransferService) GetCrossChainAttestation(ctx context.Context, receipt string, opts ...RequestOption) (*CrossChainAttestation, error) {
	return getTyped[CrossChainAttestation](s.client, ctx, buildPath("/cross-chain/attestations/{receipt}", "receipt", receipt), opts...)
}

// GetCrossChainFees get Cross-Chain Fees
func (s *CrossChainTransferService) GetCrossChainFees(ctx context.Context, params *GetCrossChainFeesParams, opts ...RequestOption) (*CrossChainFees, error) {
	return getWithParamsTyped[CrossChainFees](s.client, ctx, "/cross-chain/fees", params, opts...)
}

// ListCrossChainTransfers list Cross-Chain Transfers
func (s *CrossChainTransferService) ListCrossChainTransfers(ctx context.Context, address string, params *ListCrossChainTransfersParams, opts ...RequestOption) *Iterator[CrossChainTransfer] {
	return paginateTyped[CrossChainTransfer](s.client, ctx, buildPath("/cross-chain/transfers/{address}", "address", address), "items", params, opts...)
}
