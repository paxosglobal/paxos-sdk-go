// Code generated from the Paxos OpenAPI specification. DO NOT EDIT.

package paxos

import (
	"context"
)

// UserAuditLogsService provides access to the UserAuditLogs API.
//
// The User Audit Logs API provides a paginated audit trail of actions
// performed by human users in the Paxos dashboard, with optional filtering by
// HTTP method, path, domain, time range, and status code.
type UserAuditLogsService struct {
	client *Client
}

// ListUserAuditLogs list user audit logs
//
// Requires scope: user_audit_logs:read_audit_log
func (s *UserAuditLogsService) ListUserAuditLogs(ctx context.Context, params *ListUserAuditLogsParams, opts ...RequestOption) *Iterator[UserAuditLog] {
	return paginateTyped[UserAuditLog](s.client, ctx, "/user-audit-logs", "items", params, opts...)
}
