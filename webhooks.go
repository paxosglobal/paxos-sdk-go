// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// WebhooksService provides access to the Webhooks API.
type WebhooksService struct {
	client *Client
}

// ListWebhooks list webhooks
//
// Requires scope: webhooks:read_webhook
func (s *WebhooksService) ListWebhooks(ctx context.Context, params *ListWebhooksParams, opts ...RequestOption) *Iterator[Webhook] {
	return paginateTyped[Webhook](s.client, ctx, "/webhooks", "items", params, opts...)
}

// ListDeliveryHistories list delivery history across webhooks
//
// Requires scope: webhooks:read_webhook
func (s *WebhooksService) ListDeliveryHistories(ctx context.Context, params *ListDeliveryHistoriesParams, opts ...RequestOption) *Iterator[DeliveryHistorySummary] {
	return paginateTyped[DeliveryHistorySummary](s.client, ctx, "/webhooks/delivery-history", "items", params, opts...)
}

// GetWebhook get webhook by ID
//
// Requires scope: webhooks:read_webhook
func (s *WebhooksService) GetWebhook(ctx context.Context, webhookID string, opts ...RequestOption) (*Webhook, error) {
	return getTyped[Webhook](s.client, ctx, buildPath("/webhooks/{webhook_id}", "webhook_id", webhookID), opts...)
}
