# Security

> Claude must follow these rules without exception.

## Secrets

- Never hardcode credentials, tokens, or keys
- Never write to `.env` files (blocked by hook + deny list)
- Never log values from auth headers, passwords, or tokens
- All secrets via env vars or secret manager — see `.mcp.json` for the pattern
- If a column name contains `token`, `secret`, `password`, `key`, or `credential` — flag it in review

## Auth Model

- **Validation layer:** HTTP middleware in `internal/handler/` only — never in `internal/service/` or `internal/domain/`
- **Algorithm:** RS256 (asymmetric) for any multi-service, OAuth, or third-party scenario — public key can be freely distributed. HS256 only in closed, fully-trusted single-service deployments where both signer and verifier are equally protected.
- **Library:** `github.com/golang-jwt/jwt/v5` (JWT/JWS only, simpler API) or `github.com/lestrrat-go/jwx/v2` (full JOSE suite, JWE, JWKS tooling). **Always** restrict accepted algorithms — never omit `jwt.WithValidMethods()` / `jwt.WithKey()`. Omitting it leaves the parser open to algorithm-confusion attacks (RS256 → HS256 downgrade).
- **Key rotation:** Use `jwk.Cache` (`lestrrat-go/jwx`) for automatic background JWKS refresh from a stable provider (Google, AWS, Azure). For `golang-jwt/jwt` users, use `github.com/MicahParks/keyfunc` to bridge JWKS endpoints into `jwt.Keyfunc` with on-demand refresh on unknown `kid`. Retrieved JWKS objects are read-only shared state — never mutate.
- **Claims propagation:** After validation, store parsed claims in `context.Context` using an **unexported package-local key type** — never a plain string or exported type. Expose only `NewContext(ctx, claims)` and `FromContext(ctx)` accessors. Service layer consumes context values; it never receives raw tokens.

```go
// internal/handler/auth.go — canonical pattern
type contextKey struct{}

func NewContext(ctx context.Context, claims Claims) context.Context {
    return context.WithValue(ctx, contextKey{}, claims)
}

func FromContext(ctx context.Context) (Claims, bool) {
    c, ok := ctx.Value(contextKey{}).(Claims)
    return c, ok
}
```

<!-- Team: specify signing key source (HSM, KMS, identity provider JWKS URL) and token lifetime (access token TTL, refresh token strategy) -->

## What Claude Must Never Do

- Commit `.env` or any file containing a secret
- Generate code that stores plaintext passwords
- Disable TLS verification (`InsecureSkipVerify: true`)
- Introduce SQL string concatenation — parameterized queries only
- Add `// nolint` to security-related linter warnings without team approval
- Call `exec.Command` with user-supplied input — allowlist commands, never interpolate user data
- Write outbound HTTP calls without URL validation — validate or allowlist the host first

## Vulnerability Classes to Avoid

| Class | Rule |
|-------|------|
| SQL injection | Parameterized queries always — `db.QueryContext(ctx, query, args...)`, never `fmt.Sprintf` into SQL |
| XSS | Escape all user-controlled output; in Go templates use `html/template`, never `text/template` for HTML |
| SSRF | Validate or allowlist outbound URLs; never fetch a URL constructed from user input without host check |
| Command injection | Never `exec.Command(userInput)` — allowlist commands and args separately |
| Path traversal | Never use user input in file paths; use `filepath.Clean` and validate stays within allowed root |
| Open redirect | Validate redirect targets against an allowlist of trusted hosts |

## Dependency Security

- Run `go mod tidy` before commit — remove unused dependencies
- Pin dependencies to specific versions in `go.sum`
- Flag any dependency added without a clear justification in the PR description
