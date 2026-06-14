# Union Type Pattern

Model "one of several types" using `any` for open runtime unions, or `~A | ~B` constraints for generic compile-time unions.

## What It Is (and Isn't)

A **union type** allows a value to hold one of several types, but unlike a sum type the set is **not sealed** — the caller decides which types appear.

In TypeScript this is a first-class feature:

```typescript
type StringOrNumber = string | number;
```

Go has no native union type. It offers two different approximations depending on whether you need runtime or compile-time behavior:

| Scenario | Tool | Set | Exhaustive? |
|---|---|---|---|
| "Any type at runtime" (open set) | `any` + type switch | Open | No — `default` required |
| "One of these types in generics" | Type constraint with `\|` | Closed | Compiler-enforced |
| "Exactly one variant, all known" | Sealed interface | Closed | `panic("unreachable")` |

That last row is the **Sum Type Pattern** — see [sum-type-pattern](../sum-type-pattern/pattern.md) for that.

## Where You See It

**Runtime union (`any`)**
- JSON/YAML decoders that return `map[string]any`
- Dynamic configuration values (`string | int | bool | []string`)
- Expression evaluators where leaf values are different types

**Generic constraint union (`~A | ~B`)**
- Math helpers that work on all numeric types: `func Sum[N Number](vals []N) N`
- Sorting/comparison utilities that accept ordered types
- Serialization helpers that accept `string | []byte`

## Real Example

### 1. Runtime union — `any` + type switch (open set)

```go
// ConfigValue holds any scalar config value.
// The type switch dispatches on the actual runtime type.
type ConfigValue any

func describe(v ConfigValue) string {
    switch val := v.(type) {
    case string:
        return fmt.Sprintf("string(%q)", val)
    case int:
        return fmt.Sprintf("int(%d)", val)
    case bool:
        return fmt.Sprintf("bool(%v)", val)
    case []string:
        return fmt.Sprintf("[]string(len=%d)", len(val))
    default:
        // default is REQUIRED — the set is open, unknown types are possible
        return fmt.Sprintf("unknown(%T)", val)
    }
}
```

The `default` case is essential here: you cannot enumerate all possible types.

### 2. Generic type-set union — compile-time constraint (closed set)

```go
// Number constrains a type parameter to the common numeric types.
// This is a union at the type-parameter level — not a runtime value type.
type Number interface {
    ~int | ~int64 | ~float64
}

func Sum[N Number](vals []N) N {
    var total N
    for _, v := range vals {
        total += v
    }
    return total
}
```

The compiler rejects any type argument that doesn't satisfy `Number`. The set IS closed, but this constraint can only be used as a type parameter bound — you cannot declare `var x Number`.

## Sum Type vs Union Type — At a Glance

```
Sum type (sealed interface):
  PaymentStatus = Pending | Completed | Failed | Refunded
  → closed set, all variants in one package
  → type switch with panic("unreachable"), no default needed

Union type — runtime (any):
  ConfigValue = string | int | bool | ... | anything
  → open set, any type can appear
  → type switch with default required

Union type — generic constraint (~A | ~B):
  Number = ~int | ~int64 | ~float64
  → closed set, but only usable as a type parameter bound
  → compiler rejects non-matching type arguments
```

## Gotchas

**`any` loses compile-time safety**: The compiler won't warn if you forget a type in the switch. Add tests that cover each expected variant.

**Generic constraint union ≠ value type**: You cannot do `var x Number = 42`. The `|` syntax in an interface only works as a type parameter constraint — it does not create a value that holds "one of these types" at runtime.

**Prefer sum types when the set is truly fixed**: If you control all variants and they live in one package, a sealed interface gives more safety than `any`. See [sum-type-pattern](../sum-type-pattern/pattern.md).
