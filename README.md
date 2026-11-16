# Paxos Go SDK

Official Go SDK for the [Paxos API](https://docs.paxos.com).

> Beta: we aim to avoid breaking changes, but they may occur before 1.0. For support, contact Paxos through standard Paxos support channels; we do not handle issues through GitHub.

## Installation

Requires Go 1.26 or later.

```bash
go get github.com/paxosglobal/paxos-sdk-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	paxos "github.com/paxosglobal/paxos-sdk-go"
)

func main() {
	client, err := paxos.NewClient(
		os.Getenv("PAXOS_CLIENT_ID"),
		os.Getenv("PAXOS_CLIENT_SECRET"),
		paxos.WithSandbox(),
		paxos.WithScopes(paxos.ScopeTransferReadTransfer),
	)
	if err != nil {
		log.Fatal(err)
	}

	transfer, err := client.Transfers.GetTransfer(context.Background(), "txn_123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(transfer.ID, transfer.Status)
}
```

## Sandbox

Use `paxos.WithSandbox()` to target the Paxos sandbox environment for testing:

```go
client, err := paxos.NewClient(
	os.Getenv("PAXOS_CLIENT_ID"),
	os.Getenv("PAXOS_CLIENT_SECRET"),
	paxos.WithSandbox(),
	paxos.WithScopes(paxos.ScopeTransferReadTransfer),
)
```

## Authentication

The SDK uses OAuth2 client credentials. Token acquisition, caching, and refresh are handled automatically. Pass your credentials and at least one OAuth scope (required) to `NewClient`. Use the `paxos.Scope*` constants; each method's GoDoc states the scope it requires (for example `Requires scope: transfer:read_transfer`):

```go
client, err := paxos.NewClient(clientID, clientSecret, paxos.WithScopes(paxos.ScopeTransferReadTransfer))
```

## Error Handling

All API errors are returned as `*paxos.Error` with RFC 7807 fields. Use helper functions for specific error types:

```go
transfer, err := client.Transfers.GetTransfer(ctx, "txn_123")
if err != nil {
	if paxos.IsNotFound(err) {
		fmt.Println("transfer not found")
	} else if paxos.IsRateLimited(err) {
		fmt.Println("rate limited after SDK retries; this call has finished")
	} else {
		var apiErr *paxos.Error
		if errors.As(err, &apiErr) {
			// Always log request_id -- it is needed for Paxos support
			log.Printf("API error: status=%d request_id=%s detail=%s",
				apiErr.Status, apiErr.RequestID, apiErr.Detail)
		}
	}
}
```

`paxos.RequestID(err)`, `paxos.StatusCode(err)` (0 when none, e.g. network and timeout errors), `paxos.Detail(err)` and `paxos.Title(err)` (the API error's title; empty for network and timeout errors; `OAuth authentication failed` for authentication errors) work on any SDK error (`*Error`, `*NetworkError`, `*TimeoutError`, `*AuthenticationError`), including wrapped ones, so one branch can map every failure to a response.

An `*paxos.Error` renders as `<title> (status N, request_id: ...): <detail>`. When the response body is empty or not a JSON object, the title is the HTTP status line (for example `409 Conflict`); `API Error (status N)` appears when a JSON error body carries no string title, or for an `Error` built directly with an empty title.

`paxos.IsUnauthorized` matches an API endpoint's 401 only. OAuth token failures are a distinct `*paxos.AuthenticationError`; test for them with `paxos.IsAuthenticationError`. As in the other SDKs, an authentication error is also an API error: `errors.As(err, &apiErr)` with `var apiErr *paxos.Error` matches it and yields the status, request ID, the title `OAuth authentication failed` and the OAuth code as detail, while the status helpers (`IsNotFound`, `IsServerError` and the rest) do not match it. A transport failure reaching the token endpoint is a `*paxos.NetworkError` or `*paxos.TimeoutError`, not an authentication error.

Error helpers: `IsCanceledError`, `IsNotFound`, `IsRateLimited`, `IsUnauthorized`, `IsForbidden`, `IsValidationError`, `IsConflict`, `IsUnprocessable`, `IsServerError`, `IsTimeoutError`, `IsNetworkError`, `IsAuthenticationError`.

A 2xx body that cannot be decoded, including a JSON `null` body on a call with a response type, a bad list envelope, a non-array collection, a null or non-object list element, or a bad item during pagination, returns an `*paxos.Error` titled `Malformed response` with the real HTTP status and request ID; its message and cause name only the parser's error type, never body text, and the request is not retried.

Transport failures (DNS, connection refused, TLS, reset) return `*paxos.NetworkError`, not `*paxos.Error`; `errors.Is`/`errors.As` reach the underlying cause. `*paxos.NetworkError` and `*paxos.TimeoutError` carry `RequestID`, the `Paxos-Request-Id` of the most recent attempt built (empty if none); a network failure or timeout does not prove the request never reached Paxos, so quote it to support.

The SDK automatically retries on 429 (rate limit) and 500, 502, 503 and 504 (server errors) with exponential backoff, up to 3 retries.

HTTP 500, 502, 503 and 504 responses are retried for GET, HEAD, OPTIONS and DELETE, and for POST, PUT and PATCH requests on operations whose ref_id the SDK manages (create operations with a ref_id). Other POST, PUT and PATCH requests, including UpdateIdentity and UpdateAccount, are returned to the caller without a retry because the request may already have taken effect; check the resource's state before retrying such a call yourself.

## Pagination

List methods return `*Iterator[T]` that automatically fetches subsequent pages:

```go
iter := client.Transfers.ListTransfers(ctx, &paxos.ListTransfersParams{
	Limit:      paxos.Int32(50),
	ProfileIDs: &[]string{"p_123"},
})
for iter.Next() {
	t := iter.Current()
	fmt.Println(t.ID, t.Status)
}
if err := iter.Err(); err != nil {
	log.Fatal(err)
}
```

`iter.Collect()` returns all items as a slice (on failure, the items fetched so far together with the error), and `iter.All()` is an `iter.Seq2[T, error]` for `for t, err := range iter.All()`; both release the iterator, and breaking out of `All()` closes it. Use `paxos.Ptr(v)` to take a pointer to any value for optional fields.

`paxos.Ptr(100)` is `*int`, so typed fields such as `Limit` need `paxos.Int32(100)` (also `paxos.Int64`, `paxos.String`, `paxos.Bool`).

Always check `iter.Err()` after the loop completes. Termination is guaranteed even if the
server misbehaves — a non-advancing cursor, or an empty page alongside a cursor, both stop
iteration instead of looping forever.

`Limit` controls page size, not the total number of results — it's re-sent unchanged on
every page. To bound the whole iteration instead, pass `WithMaxItems` and/or
`WithMaxDuration`:

```go
iter := client.Transfers.ListTransfers(ctx,
	&paxos.ListTransfersParams{Limit: paxos.Int32(100), ProfileIDs: &[]string{"p_123"}},
	paxos.WithMaxItems(250),               // stop after 250 items total (~3 pages)
	paxos.WithMaxDuration(10*time.Second), // or after 10s, whichever comes first
)
```

A `WithMaxDuration` deadline that elapses surfaces as an error wrapping
`context.DeadlineExceeded` from `iter.Err()` (check with `paxos.IsTimeoutError(err)`),
rather than silently returning a partial result. If you break out of a loop before an
iterator with `WithMaxDuration` naturally finishes, call `iter.Close()` to release its
deadline timer promptly instead of waiting for it to expire on its own.

### Pages

`iter.Pages()` is an `iter.Seq2[paxos.Page[T], error]` that yields one fetched page at a time, lazily, with that page's `Items`, `NextCursor` (empty on the last page, except that an empty page reports the cursor as received, unless it repeats the cursor just sent) and its own `Metadata` (request ID, status code, raw body). It honours the same termination guarantees, `WithMaxItems` (a page trimmed by it reports an empty `NextCursor`; a page ending exactly on the cap keeps its cursor) and `WithMaxDuration`.

```go
iter := client.Transfers.ListTransfers(ctx, &paxos.ListTransfersParams{Limit: paxos.Int32(50)})
defer iter.Close()
for page, err := range iter.Pages() {
	if err != nil {
		return err
	}
	w.Header().Set("X-Next-Cursor", page.NextCursor)
	return json.NewEncoder(w).Encode(page.Items)
}
```

A caller can tell a capped result from end of data by comparing the item count with the cap: a page trimmed by `WithMaxItems` reports an empty `NextCursor`, while a page ending exactly on the cap keeps its cursor.

Unlike `All()`, breaking out of `Pages()` leaves the iterator open, so `defer iter.Close()`. An empty page (including an empty first page) is yielded once, then iteration stops. Breaking out after the first page fetches exactly one page; resume later by passing its `NextCursor` as `PageCursor` on a new list call. Mixing `Pages()` with `Next`, `Collect` or `All` on one iterator is an error in both directions while it is live (a failed or closed iterator instead ends in either view, with no request: `Err()` keeps the original error, and `Pages()` yields nothing once closed, or re-yields the stored error once per range after a failure, where TypeScript and Python end silently and keep the error on `Err()`). Calling `Pages()` again after a break continues from the next unfetched page. To resume later, break after the pages you want and keep the last `NextCursor`; do not rely on `WithMaxItems` for that, because a page trimmed by it has no cursor.

## Raw JSON

Every top-level object the SDK returns or yields from a response (each list item, and the result of each get, create and update call) has `RawJSON() json.RawMessage`: the object's JSON exactly as it appeared in the response body, original bytes and whitespace within the object included. It reflects the wire, independent of the typed fields, so a passthrough can forward it without re-encoding, and it keeps fields the SDK does not model.

```go
transfer, err := client.Transfers.GetTransfer(ctx, "txn_123")
if err != nil {
	return err
}
w.Write(transfer.RawJSON())
```

`RawJSON()` returns a fresh copy on each call, so modifying it cannot affect the object or later pages. It is nil for a value you constructed yourself. Nested objects have no raw JSON, and it never appears in `json.Marshal` output or request bodies. Responses that are plain maps (the delete and sandbox operations returning `map[string]interface{}`) have no `RawJSON()`.

Models that carry it include one unexported field, which affects equality in tests: `reflect.DeepEqual` and testify's `assert.Equal` return false for a decoded model versus an identically hand-built one, and `cmp.Equal` panics on the unexported field unless you pass `cmpopts.IgnoreUnexported(paxos.Transfer{})`. Compare marshalled values with `assert.JSONEq` instead, or ignore the field. Other SDKs expose the same content as a parsed dict in Python, compact JSON text in Java and a parsed value in TypeScript. The envelope of a page is available as `Page.Metadata.RawBody`.

## Decimal Handling

All monetary amounts use `paxos.Decimal` (arbitrary precision; arithmetic returns `paxos.Decimal`) to avoid floating-point precision loss:

```go
amount := paxos.MustDecimal("1.50")          // panics on invalid input
amount, err := paxos.NewDecimalFromString("1.50") // returns error
zero := paxos.NewDecimal(0)                       // from int64
```

The Decimal type serializes as JSON strings on the wire, matching the Paxos API convention.

Arithmetic (`Add`, `Sub`, `Mul`, `Div`) returns a `paxos.Decimal`, so results pass straight to the SDK, for example `amount.Add(fee)`.

## Per-Request Options

Override timeout or add extra headers for individual requests:

The default total call budget is 30 seconds, covering authentication, retries,
backoff and complete response-body receipt. Configure the client with
`WithDefaultCallTimeout`, and override individual calls with `WithCallTimeout`.
`WithDefaultCallAttemptTimeout` and `WithCallAttemptTimeout` optionally cap each
HTTP exchange; no additional attempt cap applies by default.

Request options override client defaults. Omitted options inherit. Duration
setters reject zero and negative values; use `WithoutCallTimeout` or
`WithoutCallAttemptTimeout` to disable a budget per request, and the corresponding
`WithoutDefaultCallTimeout` / `WithoutDefaultCallAttemptTimeout` client options.
Caller context cancellation still applies. The SDK copies a supplied HTTP client
and overwrites the copy's `Timeout` with its own per-attempt budget, so a
`Timeout` set on a client passed to `WithHTTPClient` is ignored; use
`WithDefaultCallTimeout` and `WithDefaultCallAttemptTimeout` instead. Your
client's `Transport` (and its own phase limits) is kept.
Timeout errors are `*paxos.TimeoutError` (wrapping `context.DeadlineExceeded`), recognized by `paxos.IsTimeoutError(err)`. A call cancelled through its context returns `*paxos.CanceledError` (wrapping `context.Canceled`, so `errors.Is(err, context.Canceled)` holds) with the `RequestID` of the most recent attempt built; a cancelled mutation may still have reached Paxos, so reconcile with that ID.

```go
import "time"

transfer, err := client.Transfers.GetTransfer(ctx, "txn_123",
	paxos.WithCallTimeout(60 * time.Second),
)

// Extra headers
transfer, err = client.Transfers.GetTransfer(ctx, "txn_123",
	paxos.WithExtraHeaders(http.Header{"X-Trace-Id": {"trace-abc"}}),
)
```

## Request/Response Hooks

Register one request hook and one response hook (registering again replaces the previous one):

```go
client, err := paxos.NewClient(clientID, clientSecret,
    paxos.WithScopes(paxos.ScopeTransferReadTransfer),
    paxos.WithRequestHook(func(req *http.Request) error {
        req.Header.Set("X-Trace-Id", traceID)
        return nil
    }),
    paxos.WithResponseHook(func(req *http.Request, resp *http.Response) error {
        log.Printf("<- %d %s", resp.StatusCode, req.URL.Path)
        return nil
    }),
)
```

- Hooks run once per HTTP attempt, including retries and the single 401 refresh retry. Neither runs for the OAuth token request.
- The request hook runs before signing and may add or change headers. Changing the method, URL or body fails the call and nothing is sent. So does changing or removing an SDK-owned header (`Authorization`, `Content-Type`, `Accept`, `Paxos-SDK-Language`, `Paxos-SDK-Version`) or adding `Paxos-Signature`. Adding headers and overriding `Paxos-Request-Id` are allowed.
- The response hook runs before 401, retry and status handling, and only when a response was received. `resp.Body` is buffered, so reading it does not affect the SDK.
- A hook error fails the call, is returned wrapped (`errors.Is`/`errors.As` work), and is not retried.

## Concurrency

The client is safe to share across goroutines. Go's concurrency model (goroutines + channels) replaces the async/await pattern found in other languages -- simply call SDK methods from multiple goroutines as needed.

## Request Signing

Request signing is configured per OAuth client. Setup follows the guide (https://docs.paxos.com/guides/developer/request-signing), which says request signing "must be activated by Paxos" and to contact Support (https://support.paxos.com). Credentials with the crypto-withdrawal scope also need Paxos Support to register the public key before signing works.

1. Generate a private key with OpenSSL 3 or later (check with `openssl version`; the macOS system `openssl` is LibreSSL and does not support Ed25519). For Ed25519 (EdDSA): `openssl genpkey -algorithm ed25519 -out private.pem`. For P-256 (ES256): `openssl ecparam -name prime256v1 -genkey -noout -out private.pem`
2. Derive the public key. For Ed25519: `openssl pkey -in private.pem -pubout -out public.pem`. For P-256: `openssl ec -in private.pem -pubout -out public.pem`
3. In the Dashboard, go to Developer > API, open the credential and choose Edit. Click "Add Public Key", enter a nickname, select the algorithm that matches your key (EdDSA or ES256), paste the contents of `public.pem`, click "Save Key", then save the credential with "Save API Credential".
4. Copy the Key ID (a UUID) shown on the saved public key.
5. Activate signing. For credentials without the crypto-withdrawal scope, turn on "Enable Request Signing" (available once the credential has a public key), save the credential, then contact Support to complete activation. For credentials with the crypto-withdrawal scope, contact Support to register the key with the signing stack, and turn on "Enable Request Signing" only after they confirm.
6. Configure the SDK with the private key and the Key ID:

```go
pemKey, err := os.ReadFile("private.pem")
if err != nil {
    return err
}
client, err := paxos.NewClient(clientID, clientSecret,
    paxos.WithScopes(paxos.ScopeTransferReadTransfer),
    paxos.WithRequestSigning(string(pemKey), keyID),
)
```

`paxos.WithRequestSigning(privateKey, keyID)` takes both values as strings. The private key can be the PEM file contents as written by step 1 (an unencrypted PKCS#8 `-----BEGIN PRIVATE KEY-----` block for either key type, or an unencrypted SEC1 `-----BEGIN EC PRIVATE KEY-----` block for P-256) or, for Ed25519 only, 128 hex characters (the 32-byte seed followed by the matching 32-byte public key). Pass the PEM exactly as openssl writes it (multi-line, with the BEGIN and END lines); reformatted PEM, such as a single line with spaces, may be rejected. Load it from a file or a secret store (keep its newlines if it comes from an environment variable); do not put key literals in source. Encrypted PEM is rejected: decrypt it first. Whatever the format, an invalid key makes `NewClient` return an error.

The SDK supports Ed25519 (EdDSA) and P-256 (ES256) keys and picks the algorithm from the key; other key types and curves are rejected.

Only API requests are signed; OAuth token requests are not. Each signed request carries a JWS in the `Paxos-Signature` header covering the method, the full transmitted path (including the `/v2` prefix), the exact serialized query string and the body bytes. Signing runs after request hooks, so it covers the final request.

Once signing is enforced for the credential, unsigned or mis-signed requests are rejected with HTTP 401, surfaced as `*paxos.Error` for which `paxos.IsUnauthorized(err)` is true; `paxos.RequestID(err)` returns the request ID. Quote the request ID when contacting Support.

## Debug Logging

Enable debug output with the `PAXOS_LOG_LEVEL` environment variable:

```bash
PAXOS_LOG_LEVEL=debug go run main.go
```

Or provide a custom `*slog.Logger`:

```go
import "log/slog"

logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
client, err := paxos.NewClient(clientID, clientSecret,
    paxos.WithScopes(paxos.ScopeTransferReadTransfer),
	paxos.WithLogger(logger),
)
```

The explicit logger takes precedence over `PAXOS_LOG_LEVEL`.

Debug logging never writes request or response bodies. It logs the method, path (no query string), status, request ID, elapsed time and retry attempt.

## Response Metadata

Capture metadata alongside typed service results:

```go
var meta paxos.ResponseMetadata
transfer, err := client.Transfers.GetTransfer(ctx, "txn_123", paxos.WithResponseMetadata(&meta))
fmt.Println(meta.RequestID, meta.StatusCode)
if err != nil { return err }
fmt.Println(transfer.ID)
```

The holder records the terminal HTTP response, including API errors; it is empty
if no terminal API response exists. Use a fresh holder for each concurrent call.
Pagination overwrites it per page, so after an auto-paginated call it reflects the final page only; a response hook captures every page. See the [developer guide](https://docs.paxos.com/sdks).

## Request Bodies

Empty lists and empty map fields are omitted from request bodies; empty objects for nested model fields are sent as given. In Go, a required nested struct left entirely unset (its zero value) is omitted.

## Reference IDs

For operations that take a `RefID`, such as `CreateOrchestrationRequest`, the SDK fills an empty `RefID` with a UUID once per method call and reuses it across internal retries. Set it yourself to persist it before submission.

For a passthrough or BFF, forward the unmodified JSON of a call from `meta.RawBody`
(`[]byte`, set for every received response, empty for an empty body) instead of re-marshalling the typed model:

```go
var meta paxos.ResponseMetadata
_, err := client.Transfers.GetTransfer(ctx, "txn_123", paxos.WithResponseMetadata(&meta))
if err != nil { return err }
w.Write(meta.RawBody)
```

## License

[Apache 2.0](LICENSE)
