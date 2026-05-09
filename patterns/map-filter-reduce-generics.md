# Map/Filter/Reduce with Generics

Higher-order functions for slice transformations. Type-safe with Go 1.18+ generics.

## What It Is (and Isn't)

Functional programming patterns using generics. Transform collections without explicit loops.

Not built-in (unlike other languages). Not lazy by default. You implement or use libraries.

## Where You See It

**Map - transform each element:**

```go
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

doubled := Map([]int{1, 2, 3}, func(x int) int { return x * 2 })
// [2, 4, 6]
```

**Filter - keep matching elements:**

```go
func Filter[T any](slice []T, pred func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if pred(v) {
            result = append(result, v)
        }
    }
    return result
}

evens := Filter([]int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 })
// [2, 4]
```

**Reduce - aggregate to single value:**

```go
func Reduce[T, U any](slice []T, init U, fn func(U, T) U) U {
    acc := init
    for _, v := range slice {
        acc = fn(acc, v)
    }
    return acc
}

sum := Reduce([]int{1, 2, 3}, 0, func(acc, x int) int { return acc + x })
// 6
```

## Real Example

```go
type User struct {
    Name   string
    Age    int
    Active bool
}

users := []User{
    {Name: "Alice", Age: 30, Active: true},
    {Name: "Bob", Age: 25, Active: false},
    {Name: "Charlie", Age: 35, Active: true},
}

// Get names of active users
activeNames := Map(
    Filter(users, func(u User) bool { return u.Active }),
    func(u User) string { return u.Name },
)
// ["Alice", "Charlie"]

// Total age of active users
totalAge := Reduce(
    Filter(users, func(u User) bool { return u.Active }),
    0,
    func(sum int, u User) int { return sum + u.Age },
)
// 65
```

## Gotchas

**Not lazy:**

```go
// Creates intermediate slices
result := Map(Filter(bigSlice, pred), transform)

// More efficient: single loop
result := make([]T, 0)
for _, v := range bigSlice {
    if pred(v) {
        result = append(result, transform(v))
    }
}
```

**Generic constraints:**

```go
// Need comparable for map keys
func ToMap[K comparable, V any](slice []V, keyFn func(V) K) map[K]V

// Need constraints for operations
func Sum[T constraints.Integer](slice []T) T
```

**Allocations:**

```go
// Pre-allocate when size known
func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))  // Not len(slice), 0
    // ...
}
```

## Benefits

- Declarative code
- Reusable transformations
- Type-safe
- Composable operations
- Less boilerplate

## Chaining

```go
result := Reduce(
    Map(
        Filter(data, isValid),
        transform,
    ),
    initial,
    combine,
)
```
