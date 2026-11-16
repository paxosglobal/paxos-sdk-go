# Paxos Go SDK

Guidance for coding agents integrating the Paxos Go SDK. Docs: https://docs.paxos.com/sdks (SDK guides) and https://docs.paxos.com (API reference).

## Install and import

Requires Go 1.26+.

```bash
go get github.com/paxosglobal/paxos-sdk-go
```

```go
import paxos "github.com/paxosglobal/paxos-sdk-go"
```

## Quick start

```go
client, err := paxos.NewClient(os.Getenv("PAXOS_CLIENT_ID"), os.Getenv("PAXOS_CLIENT_SECRET"),
    paxos.WithSandbox(), paxos.WithScopes(paxos.ScopeTransferReadTransfer))
if err != nil { return err }
transfer, err := client.Transfers.GetTransfer(ctx, "txn_123")
```

## Client construction, scopes, sandbox vs production, environment variables

```go
client, err := paxos.NewClient(clientID, clientSecret,
    paxos.WithSandbox(),                              // sandbox environment (omit for production)
    paxos.WithScopes(paxos.ScopeTransferReadTransfer), // OAuth2 scopes (required)
    paxos.WithHTTPClient(customClient),               // custom http.Client
    paxos.WithLogger(slogger),                        // structured debug logging
    paxos.WithRequestSigning(privateKey, keyID),      // JWS signing (PEM or hex)
    paxos.WithRequestHook(func(r *http.Request) error { return nil }),
    paxos.WithResponseHook(func(r *http.Request, resp *http.Response) error { return nil }),
)
```

- At least one scope is required: `NewClient` returns an error otherwise. `WithScopes` takes the typed `paxos.Scope` (use the `paxos.Scope...` constants; untyped string constants also work, but a `string` variable needs `paxos.Scope(s)`). Methods state the scope they need in their doc comment; a few public endpoints need none, but the client still requires one.
- Credentials do not select the environment: pass `WithSandbox()` explicitly. Sandbox and production credentials differ. `WithSandbox()` cannot be combined with `WithBaseURL` or `WithTokenURL` (set `WithBaseURL` and `WithTokenURL` together, or neither; `NewClient` returns an error if only one is set; use the pair without sandbox for a local or test server).
- Create one client per application component and share it across goroutines; this reuses OAuth tokens and connections. Auth (token acquisition, caching, refresh) is automatic. Consume an iterator from one goroutine only. There is no client `Close`; manage the lifetime of any `http.Client` you supply.
- The SDK does not read credentials itself; pass them in. The conventional variable names are:

| Variable | Read by | Description |
|---|---|---|
| `PAXOS_CLIENT_ID` | your code | OAuth2 client ID |
| `PAXOS_CLIENT_SECRET` | your code | OAuth2 client secret |
| `PAXOS_LOG_LEVEL` | SDK | `debug` to enable debug logging |

## Calling services

Services are fields on the client: `client.<Service>.<Method>(ctx, ...)`. Methods take `context.Context` first and optional `...RequestOption` last. Empty lists and empty map fields are omitted from request bodies; empty objects for nested model fields are sent as given. A query parameter set to an empty string (`paxos.String("")`) is sent as `key=`; leave it `nil` to omit it (an empty `PageCursor` is the exception and is not sent). In Go, a required nested struct left entirely unset (its zero value) is omitted.

```go
transfer, err := client.Transfers.GetTransfer(ctx, "txn_123")                                  // (*Transfer, error)
withdrawal, err := client.CryptoWithdrawals.CreateCryptoWithdrawal(ctx, &paxos.CreateCryptoWithdrawalRequest{...})
iter := client.Transfers.ListTransfers(ctx, &paxos.ListTransfersParams{Limit: paxos.Int32(50)}) // *Iterator[Transfer]
```

Optional fields are pointers; use the helpers `paxos.String("v")`, `Int`, `Int32`, `Int64`, `Bool` and generic `paxos.Ptr(v)`. There is no `Float64` helper; use `paxos.Decimal` for amounts.

Response enum fields (for example `transfer.Status`) are string types that keep values added by the API after this SDK version. Compare against constants such as `paxos.TransferStatusCOMPLETED` and always handle an unrecognized value; never treat an unknown status as success or failure.

The API omits zero values from responses, so a field can be absent even when the spec marks it required. Optional fields are `nil` pointers, required scalars are zero values (`""`, `0`) and a required timestamp is `0001-01-01T00:00:00Z`. Collections differ by model: response-only models use plain slices and maps, which are `nil` when absent and safe to range over; models shared with requests (for example `Account.Members` is `*[]AccountMember` and `Account.Metadata` is `*map[string]string`) use pointers, so nil-check before dereferencing. Decimals or timestamps may also arrive as `"0"` or the epoch; check the field's meaning before treating a zero as missing.

## Pagination

```go
iter := client.Transfers.ListTransfers(ctx, &paxos.ListTransfersParams{Limit: paxos.Int32(100)})
for iter.Next() {
    t := iter.Current()
}
if err := iter.Err(); err != nil { return err }
```

Always check `iter.Err()` after the loop. Iteration always terminates, even if the server returns a non-advancing cursor or an empty page alongside one. `iter.Collect() ([]T, error)` returns the items fetched so far with the error; `iter.All()` is an `iter.Seq2[T, error]`; both release the iterator. Call `iter.Close()` if you break out early.

`iter.Pages()` yields `paxos.Page[T]` with `Items`, `NextCursor` (the cursor to resume from) and `Metadata`. Do not mix it with `Next`/`Collect`/`All` on a live iterator.

`Limit` is a page-size hint, not a total cap, and is re-sent unchanged on every page. Bound the whole iteration with `WithMaxItems` and/or `WithMaxDuration`:

```go
iter := client.Transfers.ListTransfers(ctx,
    &paxos.ListTransfersParams{Limit: paxos.Int32(100), ProfileIDs: &[]string{"p_123"}},
    paxos.WithMaxItems(250),               // stop after 250 items total
    paxos.WithMaxDuration(10*time.Second), // or after 10s, whichever comes first
)
```

Exceeding `WithMaxDuration` surfaces from `iter.Err()` as an error wrapping `context.DeadlineExceeded` (`paxos.IsTimeoutError(err)`), not a silent truncation. The duration starts on the first `Next`, counts pauses, and is not disabled by disabling the call budget.

## Errors

```go
var apiErr *paxos.Error
if errors.As(err, &apiErr) {
    log.Printf("status=%d request_id=%s detail=%s", apiErr.Status, apiErr.RequestID, apiErr.Detail)
}
```

- Accessors `paxos.RequestID(err)`, `StatusCode`, `Detail` and `Title` work on any SDK error (empty or 0 when absent).
- Helpers: `paxos.IsNotFound(err)`, `IsRateLimited`, `IsUnauthorized`, `IsForbidden`, `IsValidationError`, `IsConflict`, `IsUnprocessable`, `IsServerError`, `IsTimeoutError`, `IsCanceledError`, `IsNetworkError`, `IsAuthenticationError`.
- HTTP status selects the error class, even when the error body is malformed or has a conflicting `status`. The request ID comes from a valid body `request_id`, else the `Paxos-Request-Id` response header, else the outgoing request's ID.
- `*paxos.NetworkError` (transport failure, cause preserved), `*paxos.TimeoutError` (wraps `context.DeadlineExceeded`) and `*paxos.CanceledError` (caller cancellation; wraps `context.Canceled`) are not `*paxos.Error` and have no status. They expose `RequestID` of the most recent attempt built (empty if none); the request may still have reached Paxos, so quote the ID to support.
- A 2xx body that cannot be decoded (including a JSON `null` body, or a bad item during pagination) returns `*paxos.Error` (`Malformed response`, real status and request ID) with a sanitized cause that contains no body text. It is not retried.
- OAuth HTTP failures and malformed token replies return `*paxos.AuthenticationError` (`Status`, `RequestID`, allowlisted OAuth `Code`; no raw token bodies), detected with `IsAuthenticationError`. This is distinct from an endpoint 401 (`IsUnauthorized`), but it is also an API error: `errors.As(err, &apiErr)` with `*paxos.Error` matches it (status, request ID, title `OAuth authentication failed`, code as detail). Status helpers do not match it. Transport failures reaching the token endpoint are `*paxos.NetworkError` or `*paxos.TimeoutError`.
- Errors on a later pagination page reach `iter.Err()` with their type, status and request ID; they never look like normal completion.

## Money and decimals

Use `paxos.Decimal` (methods: `String`, `StringFixed`, `Cmp`, `Equal`, `IsZero`, `Sign`, `Add`, `Sub`, `Mul`, `Neg`, `Abs`, `Div(o, precision)`, `Round`, `Truncate`, `IntPart`, `InexactFloat64`; also text, binary (gob) and `database/sql` encodings). Never use `float64` for money. Pass decimals directly in typed requests; the SDK writes the API's string format.

```go
d := paxos.MustDecimal("1.23")               // panics on error (constants)
d, err := paxos.NewDecimalFromString("1.23") // returns error
req := &paxos.CreateCryptoWithdrawalRequest{Amount: &d, /* ... */}
```

## Timestamps

Timestamps in models and params are `time.Time` (`*time.Time` when optional), sent as RFC 3339 at native precision, canonicalised to UTC `Z`. The API accepts any RFC 3339 offset and up to nanoseconds; stored precision varies by resource (microseconds for Postgres-backed resources, milliseconds for events). Prefer UTC values for query filters such as `CreatedAtLt`.

## Per-request options

Pass as variadic args to any service method:

```go
client.Transfers.GetTransfer(ctx, "txn_123",
    paxos.WithCallTimeout(60*time.Second),
    paxos.WithExtraHeaders(http.Header{"X-Trace-Id": {"abc"}}),
)
```

Cancellation: pass a context with a deadline or cancel function; context cancellation always wins over SDK budgets.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

## Retries, timeouts and ref_id/idempotency

- Default policy: up to 3 retries with backoff for HTTP 429, 500, 502, 503 and 504. Set `WithMaxRetries(n)` on the client or `WithRequestMaxRetries(n)` per call (integers 0 to 2147483646; 0 disables ordinary retries). One 401 token refresh is separate and does not count as a retry.
- HTTP 500, 502, 503 and 504 responses are retried for GET, HEAD, OPTIONS and DELETE, and for POST, PUT and PATCH requests on operations whose ref_id the SDK manages (create operations with a ref_id). Other POST, PUT and PATCH requests, including UpdateIdentity and UpdateAccount, are returned to the caller without a retry because the request may already have taken effect; check the resource's state before retrying such a call yourself.
- Network failures retry only for GET/HEAD/OPTIONS. An ambiguous network failure on a mutation is surfaced, never replayed. OAuth failures are surfaced without an extra retry loop. Once a rate-limit or server error reaches your code the SDK call is finished; nothing continues in the background.
- The Paxos API does not send `Retry-After` (or rate-limit headers), so the SDK does not read them; 429 retries use the SDK's own backoff. Do not build retry logic that waits for or parses `Retry-After`.
- Budgets: the default call budget is 30 seconds in total, covering OAuth, token waits, retries, backoff and reading the full response body. `http.Client.Timeout` is not used; the SDK copies a supplied client and applies its own budgets (custom transport phase limits still apply).

```go
client, err := paxos.NewClient(id, secret,
    paxos.WithScopes(paxos.ScopeTransferReadTransfer),
    paxos.WithDefaultCallTimeout(30*time.Second),
    paxos.WithDefaultCallAttemptTimeout(10*time.Second), // optional cap per HTTP attempt; unset by default
)
transfer, err := client.Transfers.GetTransfer(ctx, transferID,
    paxos.WithCallTimeout(time.Minute),
    paxos.WithoutCallAttemptTimeout(),
)
```

Request options override client settings; omit one to inherit. Durations must be positive. Disable explicitly with `WithoutDefaultCallTimeout`, `WithoutDefaultCallAttemptTimeout`, `WithoutCallTimeout` or `WithoutCallAttemptTimeout`.

- `ref_id`: for operations that accept one the SDK fills an omitted or empty value once per method call without mutating your request, and reuses it across internal retries. The field is `RefID *string` on most requests (`Account.RefID` for account creation) but a plain `RefID string` on `CreateOrchestrationRequest` and `CreateOrchestrationRuleRequest`. A second call without an ID gets a new one, so this does not give application-level deduplication.
- For jobs that may restart or be retried by a queue: generate an ID (a UUID is fine) when the job is created and persist it with the payload before the first API call; send it as `RefID`; after a timeout, transport error or exhausted retries keep the same ID and payload and reconcile via the operation's lookup or list API (for example `ListTransfers` with `RefIDs`) before resubmitting; treat a conflict as a reason to reconcile, not to mint a new ID; use a new ID only for a genuinely new operation. A timeout or cancellation does not prove the server rejected a mutation.

## Hooks and authentication errors

- `WithRequestHook` runs on every API attempt (including retries and the 401 refresh retry), before request signing, and not for OAuth token requests. It may add headers and override `Paxos-Request-Id`. Changing the method, URL or body, changing or removing an SDK-owned header (`Authorization`, `Content-Type`, `Accept`, `Paxos-SDK-Language`, `Paxos-SDK-Version`), or adding `Paxos-Signature` fails the call before anything is sent.
- `WithResponseHook` runs for every attempt that received a response, before status handling; bodies are buffered so the hook may read them. A hook error ends the call, preserves its cause for `errors.Is`/`errors.As`, prevents retry or token refresh, and takes precedence over the HTTP error. A repeated hook option replaces the earlier one.
- Token responses require HTTP 200 and a nonblank `access_token`; `token_type` must be Bearer if present; `expires_in` defaults to 3600 seconds if omitted. Violations return `*paxos.AuthenticationError`.

## Request signing

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

## Raw JSON and response metadata

Top-level response objects (list items and get/create/update results) expose `RawJSON() json.RawMessage`: that object's JSON as received, nil if built by the caller. Nested objects have none, and map-alias responses (delete and sandbox operations) have no accessor. Use it to forward unmodified JSON; decimals keep their original text and unknown fields are kept. Decoded models differ from hand-built ones under `reflect.DeepEqual`; `cmp.Equal` needs `cmpopts.IgnoreUnexported`.

Capture the API request ID, status and unmodified body of the terminal response (including errors) with a holder:

```go
var meta paxos.ResponseMetadata
transfer, err := client.Transfers.GetTransfer(ctx, "txn_123", paxos.WithResponseMetadata(&meta))
log.Printf("request_id=%s status=%d", meta.RequestID, meta.StatusCode) // meta.RawBody holds the body
```

The holder stays empty if no response headers arrived (timeout, network failure before a response). If the body read then fails or times out, `RequestID` and `StatusCode` are already set but `RawBody` is empty, so populated metadata does not mean a complete response. In an auto-paginated call it holds the final page only; a response hook sees every page. Use one holder per concurrent call. Keep the request ID for support; avoid logging identity or financial payloads.

## Logging

```bash
PAXOS_LOG_LEVEL=debug go run main.go
```

Programmatic (takes precedence over the env var):

```go
paxos.WithLogger(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
```

## Testing and mocking

The `paxostesting` package ships a local HTTP mock server and fixtures. It also serves the OAuth token endpoint.

```go
import "github.com/paxosglobal/paxos-sdk-go/paxostesting"

func TestGetTransfer(t *testing.T) {
    server := paxostesting.NewMockServer(t) // shuts down when the test ends
    server.On("GET", "/transfer/transfers/txn_123").Respond(200, paxostesting.TransferFixture("txn_123"))

    client, err := paxos.NewClient("id", "secret",
        paxos.WithScopes(paxos.ScopeTransferReadTransfer), server.ClientOption())
    if err != nil { t.Fatal(err) }
    transfer, err := client.Transfers.GetTransfer(context.Background(), "txn_123")
    if err != nil { t.Fatal(err) }
    if transfer.ID != "txn_123" { t.Errorf("unexpected transfer %q", transfer.ID) }
    if n := len(server.Calls()); n != 1 { t.Errorf("expected 1 recorded request, got %d", n) }
}
```

Further `MockServer` methods: `ReturnNotFound`, `ReturnRateLimit`, `ReturnServerError`, `ReturnValidationError` and `Reset`; `Route` has `RespondWithHeaders`. Other package-level fixtures: `ProfileFixture`, `OrderFixture`, `ListResponse` and `ListResponseWithKey` (for list methods that read a key other than `items`: addresses, events, payments, statements, issuer_quote_executions). Unmatched requests return 404.

## Language-specific notes

- Pointers for optionals: see the helpers in "Calling services".
- Network and timeout errors are not `*paxos.Error`: use `errors.As` on the right type or the `Is...` helpers rather than reading `apiErr.Status` unconditionally.
- Go has no async split; concurrency is goroutines sharing one client.

## Examples and docs

Runnable programs are in `examples/`: `quickstart`, `sandbox_deposit`, `create_crypto_withdrawal` and `list_transfers_pagination`. SDK guides: https://docs.paxos.com/sdks. API reference: https://docs.paxos.com.
