// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// MonitoringAddressesService provides access to the MonitoringAddresses API.
//
// A Monitoring Address is a blockchain address that Paxos monitors daily for
// eligible stablecoin activity and determines the amount of rewards your
// organization will earn. We recommend adding any addresses you expect to
// hold, receive, or mint Paxos stablecoins as a monitoring address in order to
// receive rewards for all of your qualifying on-chain activity.
//
// > Monitoring addresses must be added by month end (`23:59 UTC+0`) in order
// to earn rewards for that month. Addresses added after month end will start
// earning rewards for the following month.
//
// *These endpoints are only relevant for Global Dollar Network (GDN)
// partners.*
type MonitoringAddressesService struct {
	client *Client
}

// ListMonitoringAddress list Monitoring Address
//
// Requires scope: rewards:read_monitoring_address
//
// Deprecated: This API operation is deprecated.
func (s *MonitoringAddressesService) ListMonitoringAddress(ctx context.Context, params *ListMonitoringAddressParams, opts ...RequestOption) *Iterator[MonitoringAddress] {
	return paginateTyped[MonitoringAddress](s.client, ctx, "/rewards/monitor/address", "addresses", params, opts...)
}

// CreateMonitoringAddress create Monitoring Address
//
// Requires scope: rewards:write_monitoring_address
//
// Deprecated: This API operation is deprecated.
func (s *MonitoringAddressesService) CreateMonitoringAddress(ctx context.Context, body *CreateMonitoringAddressRequest, opts ...RequestOption) (*CreateMonitoringAddressResponse, error) {
	return postTyped[CreateMonitoringAddressResponse](s.client, ctx, "/rewards/monitor/address", body, opts...)
}

// GetMonitoringAddress get Monitoring Address
//
// Requires scope: rewards:read_monitoring_address
//
// Deprecated: This API operation is deprecated.
func (s *MonitoringAddressesService) GetMonitoringAddress(ctx context.Context, id string, opts ...RequestOption) (*GetMonitoringAddressResponse, error) {
	return getTyped[GetMonitoringAddressResponse](s.client, ctx, buildPath("/rewards/monitor/address/{id}", "id", id), opts...)
}

// UpdateMonitoringAddress update Monitoring Address
//
// Requires scope: rewards:write_monitoring_address
//
// Deprecated: This API operation is deprecated.
func (s *MonitoringAddressesService) UpdateMonitoringAddress(ctx context.Context, id string, body *UpdateMonitoringAddressRequest, opts ...RequestOption) (*UpdateMonitoringAddressResponse, error) {
	return putTyped[UpdateMonitoringAddressResponse](s.client, ctx, buildPath("/rewards/monitor/address/{id}", "id", id), body, opts...)
}
