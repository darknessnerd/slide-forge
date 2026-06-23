# Conventions

## Naming

- Files: `snake_case`
- Packages: lowercase single word — name by what it contains, not what it does (`user` not `userservice`)
- Exported types/funcs: `PascalCase`
- Unexported: `camelCase`
- Constants: `ALL_CAPS` only for truly global invariants
- Interfaces: noun or noun phrase (`Repository`, `Notifier`), not `IRepository` or `RepositoryInterface`
- Constructors: `New<Type>` returning the interface, not the concrete type

> See `.claude/skills/solid-principles.md` for package naming rules derived from SRP.

## Interfaces

**Interfaces belong in the consuming package, not the implementing package.** The package that uses a value of an interface type defines the interface — not the package that implements it. This is confirmed Go idiom from [go.dev/wiki/CodeReviewComments](https://go.dev/wiki/CodeReviewComments#interfaces): "Go interfaces generally belong in the package that uses values of the interface type, not the package that implements those values."

**Do not pre-define interfaces for mocking.** Define interfaces only when you have two or more implementations, or when the consuming package needs to isolate itself from a specific dependency. Per go.dev/wiki/CodeReviewComments: "Do not define interfaces before they are used."

**Never store `context.Context` in a struct.** Pass it as the first parameter to every function that needs it. Per pkg.go.dev/context: "Do not store Contexts inside a struct type; instead, pass a Context explicitly to each function that needs it."

```go
// WRONG — interface defined in the implementing package:
// domain/user.go
package domain
type UserStore interface {   // ← defined here, but domain doesn't use it
    Find(id string) (User, error)
}

// repository/postgres.go
package repository
type PostgresStore struct{ ... }
func (p *PostgresStore) Find(id string) (domain.User, error) { ... }

// service/user.go — now imports domain just to name the interface
package service
import "domain"
type UserService struct{ store domain.UserStore }


// RIGHT — interface defined in the consuming package (service):
// domain/user.go
package domain
type User struct{ ID, Name string }   // ← only value types here

// service/user.go
package service
// Defined here, by the consumer that needs it:
type userStore interface {
    Find(id string) (domain.User, error)
}
type UserService struct{ store userStore }

// repository/postgres.go — satisfies the interface implicitly; no domain import needed
package repository
type PostgresStore struct{ ... }
func (p *PostgresStore) Find(id string) (domain.User, error) { ... }
```

## File Layout

<!-- Team: describe your actual module/package structure here -->
<!-- Scaffold default below — adjust if you use a different layout -->

```
cmd/
  main.go              ← entry point, wires dependencies
internal/
  domain/              ← value types only (structs, enums, sentinel errors — no interfaces)
  service/             ← business logic; defines interfaces it needs from storage
  repository/          ← data access; implements service interfaces implicitly
  handler/             ← HTTP/gRPC; defines interfaces it needs from service layer
```

One file per logical concern. No `utils.go` or `helpers.go` — if it needs a file, it needs a package.

## Error Handling

Defaults until overridden:

- Wrap errors with context: `fmt.Errorf("user.Get: %w", err)` — include the call site
- Return errors to callers; log only at the top boundary (handler layer)
- Sentinel errors in `domain` package: `var ErrNotFound = errors.New("not found")`
- Never `panic` outside `init()` or package-level setup
- Never swallow errors with `_ = someFunc()`

## Imports

<!-- Team: confirm grouping order or override -->

Standard Go import grouping — enforced by `goimports`:

```go
import (
    // 1. stdlib
    "context"
    "fmt"

    // 2. internal packages
    "github.com/your-org/your-repo/internal/domain"

    // 3. external dependencies
    "github.com/some/library"
)
```

Blank line between each group. `goimports` or `gofumpt` handles this automatically.

## Forbidden Patterns

<!-- Team: add domain-specific forbidden patterns below -->

- No `panic` outside `init()` or package-level var initialization
- No global mutable state — pass dependencies via constructors
- No `interface{}` / `any` where a concrete type or typed interface suffices
- No `init()` functions that perform I/O or have side effects
- No `// nolint` on security-related linter warnings without team approval (see `04-security.md`)
- No type switches on concrete types to vary behavior — use interface methods instead (OCP violation)
- Constructors must accept interfaces, not concrete types (`func New(repo domain.UserReader)` not `func New(db *PostgresRepo)`)
