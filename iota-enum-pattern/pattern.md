# Iota Enum Pattern

Go's idiomatic enum: named integer constants with `iota`, a `String()` method, and a `Valid()` guard.

## What It Is (and Isn't)

Go has no `enum` keyword. The idiomatic substitute is a **named integer type** paired with a `const` block using `iota`:

```go
type OrderStatus int

const (
    Pending  OrderStatus = iota // 0
    Approved                    // 1
    Shipped                     // 2
    Delivered                   // 3
    Cancelled                   // 4
)
```

This gives you a distinct type (`OrderStatus`) with named constants — the compiler will warn if you assign a plain `int` without a conversion. But it's **not truly sealed**: any integer can be cast to `OrderStatus`, so `OrderStatus(99)` compiles and runs silently.

## Iota Enum vs Sealed Interface (Sum Type)

| | Iota enum | Sealed interface (sum type) |
|---|---|---|
| Variants carry data? | No — all variants are just `int` | Yes — each variant is a distinct struct |
| Truly closed? | No — `Status(99)` is valid | Yes — unexported method blocks external types |
| Exhaustiveness check | None (use `Valid()` + `default: panic`) | `panic("unreachable")` at switch end |
| Memory | One `int` per value | Interface header + concrete struct |
| Use case | Flags, codes, states with no payload | State machines where variants carry different data |

**Rule of thumb**: if all your variants have the same shape (just a label, no payload), use iota. If variants differ structurally (e.g., `Failed` has a reason, `Completed` has a tx ID), use a sealed interface.

## Where You See It

- HTTP methods, status codes, log levels
- Order/payment/request lifecycle states (data-free variants)
- CLI flags and configuration modes
- Anything that maps naturally to a C-style enum

## Real Example

```go
type OrderStatus int

const (
    Pending   OrderStatus = iota
    Approved
    Shipped
    Delivered
    Cancelled
    orderStatusSentinel // unexported: marks the end of valid range
)

// String makes OrderStatus print as a name, not a number.
func (s OrderStatus) String() string {
    names := [...]string{"Pending", "Approved", "Shipped", "Delivered", "Cancelled"}
    if !s.Valid() {
        return fmt.Sprintf("OrderStatus(%d)", int(s))
    }
    return names[s]
}

// Valid guards against out-of-range integer casts.
// Without this, OrderStatus(99).String() would panic on the array index.
func (s OrderStatus) Valid() bool {
    return s >= Pending && s < orderStatusSentinel
}
```

### The open-integer gotcha

```go
var s OrderStatus = OrderStatus(99) // compiles fine — no error
fmt.Println(s.Valid())              // false — caught at runtime, not compile time
```

This is the key difference from a sealed interface: the compiler doesn't prevent invalid values.

## Generating String() Automatically

For large enums, avoid writing `String()` by hand — use the standard `stringer` tool:

```go
//go:generate stringer -type=OrderStatus
```

Running `go generate` creates `orderstatus_string.go` with a generated `String()` method. Install with:

```bash
go install golang.org/x/tools/cmd/stringer@latest
```

## Gotchas

**Iota resets per const block**: each new `const (...)` block restarts `iota` at 0. Don't mix enum variants across blocks.

**Don't start at 0 if 0 means "unset"**: use `iota + 1` so the zero value is clearly invalid:

```go
const (
    _ OrderStatus = iota // discard 0
    Pending              // 1
    Approved             // 2
    ...
)
```

**Avoid bitmask collision**: `iota` is sequential, not powers of two. For bitmasks, use `1 << iota`:

```go
const (
    Read    Permission = 1 << iota // 1
    Write                          // 2
    Execute                        // 4
)
```

## Coming from Python?

Python's `enum.Enum` maps cleanly to Go's iota pattern — **not** to sum types. Python enums are data-free named labels, just like iota constants.

```python
from enum import Enum
class Color(Enum):
    RED   = 1
    GREEN = 2
    BLUE  = 3
# or: Color = Enum('Color', ['RED', 'GREEN', 'BLUE'])
```

| Python `enum.Enum` feature | Go equivalent |
|---|---|
| `Color.RED` | `Red` constant |
| `Color.RED.name` → `"RED"` | `Red.String()` → `"Red"` (via `String()` method) |
| `Color.RED.value` → `1` | `int(Red)` → `0` (underlying integer) |
| `list(Color)` — iterate all members | manual slice: `[]Color{Red, Green, Blue}` |
| `Enum('Color', [...])` functional constructor | no equivalent — use `const` block |
| Can't construct invalid member | `Color(99)` compiles — use `Valid()` to guard |

The one Python Enum feature Go lacks is **built-in iteration**. The idiomatic workaround is an explicit slice or a `Values()` function you maintain alongside the constants:

```go
type Color int

const (
    Red Color = iota
    Green
    Blue
    colorSentinel
)

// Values returns all valid Color values — maintain this when adding variants.
func (Color) Values() []Color {
    all := make([]Color, 0, colorSentinel)
    for c := Red; c < colorSentinel; c++ {
        all = append(all, c)
    }
    return all
}
```

### Is there a Go enum package?

No standard library package. The official tooling is:
- **`stringer`** (`golang.org/x/tools/cmd/stringer`) — generates `String()` only, via `go generate`

Third-party codegen tools add more Python-like features (iteration, JSON marshaling, SQL scanning):
- `github.com/abice/go-enum` — generates from `// ENUM(...)` comments
- `github.com/nikolaydubina/go-enum` — similar codegen approach

These are code generators, not runtime packages — they produce `.go` files you commit alongside your source. None are part of the standard library or officially endorsed.

## See Also

- [Sum Type Pattern](../sum-type-pattern/pattern.md) — when variants carry different data (truly sealed)
- [Union Type Pattern](../union-type-pattern/pattern.md) — open runtime unions and generic constraints
