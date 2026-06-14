# Sum Type Pattern

Model "exactly one of these types, and nothing else" using a sealed interface.

## What It Is (and Isn't)

A **sum type** (also: discriminated union, tagged union) is a type that can hold one of a fixed set of variants — and nothing outside that set.

In Rust or Haskell this is built-in:

```rust
enum PaymentStatus { Pending, Completed, Failed(String), Refunded }
```

Go has no native sum type. The closest approximation is a **sealed interface**: adding an unexported method prevents any type outside the package from satisfying the interface.

Not the same as a plain `type switch` on `interface{}` — that's an *open* set where `default` is essential. A sealed interface is a *closed* set; `default` becomes `panic("unreachable")`.

## Where You See It

- **Domain state machines**: payment status, order lifecycle, request outcome
- **AST nodes**: `Expr = Number | Add | Mul | ...`
- **Event types**: event bus where each event variant carries different data
- **Result/Option types**: `Ok(value) | Err(err)` as an explicit return type

## Real Example

```go
// PaymentStatus is sealed: only types in this package can implement it.
type PaymentStatus interface {
    paymentStatusTag() // unexported — external packages can't add variants
}

type Pending   struct{}
type Completed struct{ TxID string }
type Failed    struct{ Reason string }
type Refunded  struct{ Amount float64 }

func (Pending) paymentStatusTag()   {}
func (Completed) paymentStatusTag() {}
func (Failed) paymentStatusTag()    {}
func (Refunded) paymentStatusTag()  {}

func describe(s PaymentStatus) string {
    switch v := s.(type) {
    case Pending:
        return "awaiting payment"
    case Completed:
        return "paid, tx=" + v.TxID
    case Failed:
        return "failed: " + v.Reason
    case Refunded:
        return fmt.Sprintf("refunded %.2f", v.Amount)
    }
    panic("unreachable: sealed interface guarantees exhaustive switch")
}
```

## Bonus: Result[T] with Generics

A generic `Result` type is a sum type pattern that replaces `(T, error)` returns with an explicit value:

```go
type Result[T any] struct {
    val T
    err error
    ok  bool
}

func Ok[T any](v T) Result[T]        { return Result[T]{val: v, ok: true} }
func Err[T any](e error) Result[T]   { return Result[T]{err: e} }

func (r Result[T]) Unwrap() (T, error) { return r.val, r.err }
func (r Result[T]) IsOk() bool         { return r.ok }
```

## Sum Type vs Union Type

These are related but different concepts:

| | Sum type | Union type (runtime) | Union type (constraint) |
|---|---|---|---|
| Set | Closed | Open | Closed |
| Go tool | Sealed interface | `any` + type switch | `~A \| ~B` constraint |
| `default` in switch | `panic("unreachable")` | Required | N/A |
| External variants | Impossible | Always possible | N/A |

Use a sealed interface (sum type) when **you own all variants** and they live in one package. Use `any` (runtime union) when the set is open or caller-controlled. Use generic type-set unions for numeric/comparable type parameters.

If your variants carry **no data** (just a label), consider a simpler iota enum instead — see [iota-enum-pattern](../iota-enum-pattern/pattern.md).

See [union-type-pattern](../union-type-pattern/pattern.md) for the full comparison with runnable examples.

## Gotchas

**No compiler-enforced exhaustiveness**: Go won't warn if you miss a case. The `panic("unreachable")` at the end is the convention — it will surface at runtime if a new variant is added without updating the switch.

**Sealed = no external mocking**: Tests in other packages can't create a mock variant. Solution: expose a test-only type inside the same package, or use a small test file in `package foo` (not `package foo_test`).

**One package only**: all variants must live in the same package as the interface. This is a deliberate constraint, not a bug — it enforces the closed-world assumption that makes sum types useful.
