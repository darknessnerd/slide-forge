# Skill: solid-principles

**Trigger:** Claude is creating a new package, designing an interface, reviewing a component for structure, or the user asks about design quality.

Apply automatically during design and review. Go has no classes — map every check to Go idioms.

---

## The Five Principles — Go Mappings

### S — Single Responsibility

> A package/type does one thing. One reason to change.

**Go checks:**
- Package name is a noun that fits in 5 words or fewer (`user`, `order`, `invoice`, not `manager`, `handler_utils`, `misc`)
- If you remove the package name from a function name and it still makes sense, the function is in the wrong package (`user.GetUser` → bad; `user.Get` → good)
- A struct with more than ~5 fields that span unrelated concerns is a signal — consider splitting

**Smell:** package named `util`, `helper`, `common`, `misc` — these are SRP violations waiting to happen.

**Ask:** "What is the single thing this package/type is responsible for? What would make me change it?"

---

### O — Open/Closed

> Open for extension, closed for modification. Add behavior without editing existing code.

**Go checks:**
- Behavior variation via interfaces, not switch/if on type strings or enum flags
- New behavior = new type that satisfies an interface, not a new case in existing code
- `switch req.Type { case "A": ... case "B": ... }` in core logic = OCP violation — extract to a registry or strategy map

**Go pattern:**

```go
// Closed to modification:
type Notifier interface {
    Notify(ctx context.Context, msg Message) error
}

// Open to extension — add EmailNotifier, SMSNotifier without touching existing code:
type EmailNotifier struct{ ... }
func (e EmailNotifier) Notify(ctx context.Context, msg Message) error { ... }
```

**Ask:** "To add a new behavior, do I edit existing code or add a new type?"

---

### L — Liskov Substitution

> Any implementation of an interface must be usable wherever the interface is expected — no surprises.

**Go checks:**
- An implementation must not panic or error in cases where the interface contract says it won't
- An implementation must not require the caller to type-assert to get real behavior
- An implementation that ignores inputs silently (no-op) while the contract implies action is a violation

**Smell:**

```go
// Violation — caller must know the concrete type to get real behavior:
func Process(n Notifier) {
    if email, ok := n.(*EmailNotifier); ok {
        email.SetPriority(High) // only works for one concrete type
    }
    n.Notify(ctx, msg)
}
```

**Ask:** "Can I swap any implementation of this interface and have the caller work correctly without knowing which one it got?"

---

### I — Interface Segregation

> Small, focused interfaces. Callers depend only on what they use.

**Go idiomatic rule:** prefer 1–3 method interfaces. The standard library models this — `io.Reader`, `io.Writer`, `io.Closer`.

**Go checks:**
- If a caller only uses 2 of 8 methods on an interface, split the interface
- Passing `*sql.DB` to a function that only queries = ISP violation; pass `interface { QueryContext(...) }` instead
- Embed interfaces to compose without forcing implementors to satisfy everything

```go
// Too broad — callers that only read are forced to depend on write methods:
type UserStore interface {
    Find(id string) (User, error)
    Save(u User) error
    Delete(id string) error
    List() ([]User, error)
    Count() int
}

// Segregated:
type UserReader interface {
    Find(id string) (User, error)
    List() ([]User, error)
}

type UserWriter interface {
    Save(u User) error
    Delete(id string) error
}
```

**Ask:** "Does every caller of this interface use every method? If not, split."

---

### D — Dependency Inversion

> High-level policy must not depend on low-level detail. Both depend on abstractions.

**Go checks:**
- Constructors accept interfaces, not concrete types (`func New(repo userStore)` not `func New(db *PostgresRepo)`)
- Business logic packages must not import infrastructure packages (`service` must not import `postgres`, `redis`, `kafka`)
- **The consuming layer defines the interface it needs** — not a shared `domain/` package. `service/` defines `type userStore interface {...}` for what it needs from storage. `handler/` defines `type userService interface {...}` for what it needs from business logic. `domain/` contains only value types (structs, enums, sentinel errors) — no interfaces.

**Import graph rule for this scaffold:**

```
cmd/
 └─ main.go          ← wires everything (only place allowed to import all layers)

internal/
 ├─ domain/          ← value types only (User, Order, ErrNotFound…). Zero external imports.
 │                      No interfaces — consumers define the interfaces they need.
 ├─ service/         ← imports domain; defines type userStore interface{} for storage
 ├─ repository/      ← imports domain; implicitly satisfies service.userStore
 └─ handler/         ← imports domain; defines type userService interface{} for business logic
```

**Why this matters:** Defining interfaces in `domain/` creates a shared abstraction that every layer couples to. Defining them in the consuming package means the interface is as small as the consumer actually needs (ISP), and `repository` can satisfy multiple different consumer interfaces without knowing about any of them.

**Smell:** `service` package has `import "database/sql"` — direct DB dependency in business logic.

**Also a smell:** `domain/` package contains interface types — this forces the interface definition to live in the wrong package and typically produces over-broad interfaces.

**Ask:** "If I swap Postgres for MongoDB, how many packages change?" Answer should be: only `repository`.

---

## Review Checklist

Run this when reviewing a new package or component design:

```
[ ] S — Package has one name, one job, one reason to change
[ ] O — New behavior added via new type, not new case
[ ] L — All interface implementations are substitutable without caller changes
[ ] I — Interfaces are ≤3 methods; callers depend only on what they use
[ ] D — Business logic imports only domain; infra imports only domain; main wires all
```

## Output Format

When reporting a SOLID finding:

```
**[Principle]** `file:line` — [what the violation is] → [specific fix]
```

Example:
```
**[DIP]** `service/user.go:12` — imports `repository/postgres` directly → define a `userStore` interface in `service/` and accept that in the constructor instead
```

Report findings with severity matching the `/review` command convention: `critical` / `warning` / `suggestion`.
