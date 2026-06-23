# Testing

## Framework

- **Assertions:** `github.com/stretchr/testify/assert` and `testify/require` (require stops the test on first failure; use it for preconditions)
- **Mocks:** `github.com/stretchr/testify/mock` for hand-written mocks; `github.com/uber-go/mock/mockgen` for generated mocks on large interfaces
- **Integration deps:** `github.com/testcontainers/testcontainers-go` — spin up real Postgres/Redis in CI, never fake them
- **Test runner:** standard `go test` — no third-party runner needed

Table-driven subtests via `t.Run()` are required for any function with more than one input variation. Use `t.Parallel()` in unit tests unless the test touches shared state.

## Coverage Expectations

- **domain/** and **service/**: 100% of exported functions must have tests
- **handler/**: every route must have at minimum happy path + 400 + one error response (404 or 500)
- **repository/**: integration tests only — no unit-test coverage target; real DB via testcontainers
- **Overall floor:** 80% line coverage enforced in CI (`go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out`)
- Error paths are not optional — every `if err != nil` branch that returns a non-nil error needs a test case

## Test Layout

Tests live alongside source in the same package (white-box) unless testing the public API surface (black-box, `package foo_test`).

```
internal/
  service/
    user.go
    user_test.go        ← white-box: same package, tests unexported helpers too
  handler/
    user_test.go        ← black-box: package handler_test, tests public surface only
test/
  integration/          ← integration and e2e tests, separate from unit tests
```

One test file per source file. Integration tests in `/test/integration/`, tagged with `//go:build integration`.

## What Must Have Tests

- All exported functions in `domain/` and `service/`
- All error paths — every `if err != nil` branch that returns a non-nil error
- All HTTP handler routes — at minimum: happy path + 400 + 404/500
- Any function touching auth, tokens, or credentials

## What Must NOT Be Mocked

- **Database / SQL layer** — repository tests must use a real DB via testcontainers. Mock DB = mock bug surface. Prior incidents where mock/prod divergence masked broken migrations justify this rule permanently.
- **`context.Context`** — never replace with `context.Background()` in tests that need deadline or cancellation behaviour; use `context.WithTimeout` with a realistic value
- **Sentinel errors from `domain/`** — never redefine them in test files; import and assert against the real `domain.ErrNotFound` etc.
- **The Logger interface** — pass a real `slog.New(slog.NewTextHandler(io.Discard, nil))` in tests; never a nil or stub that silently drops output and hides panics

## TDD Workflow (Red → Green → Refactor)

Follow this loop for every new exported function:

1. **Red** — write the test first. It must fail (`go test` exits non-zero) before writing any implementation.
2. **Green** — write the minimal implementation that makes it pass. No extra logic.
3. **Refactor** — clean up without breaking tests. Race detector must still pass (`go test -race`).

For agent skills (Claude Code `.claude/skills/`), use the skill-creator plugin's evaluation loop:
- Write failing test cases (no-skill baseline) → author skill → confirm pass rate rises in `benchmark.json`
- Run description-tuning to verify trigger accuracy (should-trigger + should-not-trigger prompts) before committing

## Grading Agent Output in Tests

When a test must assert on LLM-generated text (e.g., testing a skill's output quality):

1. **Prefer deterministic checks first** — exact string match, JSON schema validation, substring presence. Fastest, most reliable, most scalable.
2. **LLM-as-judge only for complex judgment** — use a *different* model from the one that generated the output. Instruct the judge to reason first, then emit a score; discard the reasoning chain in the assertion.
3. **Volume over quality** — 20 automated cases with imperfect signal beats 5 hand-graded golden cases.

## Running Tests

```bash
# unit tests only
go test ./...

# with race detector (required before merge)
go test -race ./...

# integration tests (requires running dependencies)
go test -tags=integration ./test/integration/...

# coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```
