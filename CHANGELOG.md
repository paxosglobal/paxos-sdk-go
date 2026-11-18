# Changelog

All notable changes to the Paxos Go SDK are documented here.

## 0.1.1 - 2026-11-18

### API changes

#### New and changed

- `CreateCryptoWithdrawal`, `CreateFiatWithdrawal`, `CreateInternalTransfer`, `CreatePaxosTransfer`, `GetTransfer`, `ListTransfers`: added the optional property `sender_address` to the response.
- `GetRewardsStatement`: response field `statement.type` can now also be `STATEMENT_TYPE_EBP_REWARDS` or `STATEMENT_TYPE_MANUAL_REWARDS`.
- `ListRewardsStatements`: `type` filter now also accepts `STATEMENT_TYPE_EBP_REWARDS` or `STATEMENT_TYPE_MANUAL_REWARDS`.
- `ListRewardsStatements`: response field `type` can now also be `STATEMENT_TYPE_EBP_REWARDS` or `STATEMENT_TYPE_MANUAL_REWARDS`.

## 0.1.0

Initial beta release of the Paxos Go SDK.

Beta: we aim to avoid breaking changes, but they may occur before 1.0.

### Features

- OAuth2 client-credentials authentication with automatic token caching and refresh; at least one typed scope is required
- Automatic retries with exponential backoff and jitter on 429 and 500, 502, 503 and 504 responses (default 3 retries, configurable), plus one token-refresh retry on 401; network failures are retried only for GET, HEAD and OPTIONS
- Typed models, enums, request parameters and service clients for every covered endpoint, accessed from one `paxos.Client`
- Empty, blank or null path parameters are rejected before any request is sent
- Auto-paginating `Iterator[T]` for list methods, with `Next`, `Collect`, `All` (range-over-func) and per-page access via `Pages`, bounded by `WithMaxItems` and `WithMaxDuration`
- `paxos.Decimal` for monetary amounts, serialized as strings on the wire
- Timestamps as `time.Time`
- Idempotent mutations: an omitted `ref_id` is filled with a UUID once per call and reused across retries
- Access to the raw JSON of every top-level response object
- Typed errors for API, network, timeout, cancellation and authentication failures, each carrying the Paxos request ID
- Per-client and per-request timeouts, including an overall call budget and an optional per-attempt cap
- Request and response hooks
- Optional JWS request signing with Ed25519 (EdDSA) or P-256 (ES256) keys
- Debug logging, enabled with `PAXOS_LOG_LEVEL` or a custom logger; logs paths without query strings and never includes request or response bodies
- Response metadata (request ID, status code, raw body) available alongside typed results
- Custom HTTP client support (`WithHTTPClient`) for proxies and other transports
- Testing helpers in `paxostesting`: a local mock server that records requests, plus typed response fixtures
- Sandbox and production environments, plus a custom base URL and token URL (set together) for local or test servers
- Coverage of every endpoint in the Paxos API v2 OpenAPI specification (see the [API Reference](https://docs.paxos.com/api-reference/introduction))
- Requirements: Go 1.26 or later
- Runtime dependencies: `github.com/google/uuid` and `github.com/shopspring/decimal` (not exposed in the public API)
- Free-form JSON object fields such as `UserAuditLog.RequestPayload` and `Event.Object` are typed `map[string]interface{}`
- `paxos.NewDecimal`, `paxos.NewDecimalFromString` and `paxos.MustDecimal` constructors, and an explicit `paxos.Decimal` API (`Add`, `Sub`, `Mul`, `Div`, `Neg`, `Abs`, `Round`, `Truncate`, `IntPart`, `InexactFloat64`, `Cmp`, `Equal`, `IsZero`, `Sign`, `String`, `StringFixed`) with text, binary (gob) and `database/sql` encodings
