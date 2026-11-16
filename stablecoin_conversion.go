// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"

	"github.com/google/uuid"
)

// StablecoinConversionService provides access to the StablecoinConversion API.
//
// Convert between fiat and stablecoin using
// [StablecoinConversionService.CreateStablecoinConversion] and check the
// status of a [StablecoinConversionService.GetStablecoinConversion] or
// retrieve a [StablecoinConversionService.ListStablecoinConversions].
type StablecoinConversionService struct {
	client *Client
}

// ListStablecoinConversions list Stablecoin Conversions
//
// Requires scope: conversion:read_conversion_stablecoin
//
// Deprecated: This API operation is deprecated.
func (s *StablecoinConversionService) ListStablecoinConversions(ctx context.Context, params *ListStablecoinConversionsParams, opts ...RequestOption) *Iterator[StablecoinConversion] {
	return paginateTyped[StablecoinConversion](s.client, ctx, "/conversion/stablecoins", "items", params, opts...)
}

// CreateStablecoinConversion create Stablecoin Conversion
//
// Requires scope: conversion:write_conversion_stablecoin
//
// Deprecated: This API operation is deprecated.
func (s *StablecoinConversionService) CreateStablecoinConversion(ctx context.Context, body *CreateStablecoinConversionRequest, opts ...RequestOption) (*StablecoinConversion, error) {
	if body != nil && (body.RefID == nil || *body.RefID == "") {
		bodyCopy := *body
		refID := uuid.NewString()
		bodyCopy.RefID = &refID
		body = &bodyCopy
	}
	opts = append([]RequestOption{withReplayProtected()}, opts...)
	return postTyped[StablecoinConversion](s.client, ctx, "/conversion/stablecoins", body, opts...)
}

// GetStablecoinConversion get Stablecoin Conversion
//
// Requires scope: conversion:read_conversion_stablecoin
//
// Deprecated: This API operation is deprecated.
func (s *StablecoinConversionService) GetStablecoinConversion(ctx context.Context, id string, opts ...RequestOption) (*StablecoinConversion, error) {
	return getTyped[StablecoinConversion](s.client, ctx, buildPath("/conversion/stablecoins/{id}", "id", id), opts...)
}

// CancelStablecoinConversion cancel Stablecoin Conversion
//
// Requires scope: conversion:write_conversion_stablecoin
//
// Deprecated: This API operation is deprecated.
func (s *StablecoinConversionService) CancelStablecoinConversion(ctx context.Context, id string, opts ...RequestOption) (*StablecoinConversion, error) {
	return deleteTyped[StablecoinConversion](s.client, ctx, buildPath("/conversion/stablecoins/{id}", "id", id), opts...)
}
