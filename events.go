// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// EventsService provides access to the Events API.
//
// The Events API allows you to fetch events that occurred on the Paxos
// platform with the full event payload. This REST API can be used as part of
// your webhook integrations (including rebuilding your event history should
// your webhook consumer exceed its retries), or standalone by polling for
// events related to your identities.
type EventsService struct {
	client *Client
}

// ListEvents list Events
//
// Requires scope: events:read_event
func (s *EventsService) ListEvents(ctx context.Context, params *ListEventsParams, opts ...RequestOption) *Iterator[Event] {
	return paginateTyped[Event](s.client, ctx, "/events", "events", params, opts...)
}

// GetEvent get Event
//
// Requires scope: events:read_event
func (s *EventsService) GetEvent(ctx context.Context, id string, params *GetEventParams, opts ...RequestOption) (*Event, error) {
	return getWithParamsTyped[Event](s.client, ctx, buildPath("/events/{id}", "id", id), params, opts...)
}
