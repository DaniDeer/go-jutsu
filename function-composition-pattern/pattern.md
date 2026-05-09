# Function Composition Pattern

Combine functions to create new functions. `f(g(x))` becomes `Compose(f, g)(x)`.

## What It Is (and Isn't)

Mathematical function composition. Build complex operations from simple ones.

Not method chaining. Not OOP. Pure functional programming in Go.

## Where You See It

**Basic composition:**

```go
func Compose[A, B, C any](f func(B) C, g func(A) B) func(A) C {
    return func(x A) C {
        return f(g(x))
    }
}

add5 := func(x int) int { return x + 5 }
double := func(x int) int { return x * 2 }

add5ThenDouble := Compose(double, add5)
result := add5ThenDouble(3)  // (3 + 5) * 2 = 16
```

**Pipeline builder:**

```go
type Transform[T any] func(T) T

func (t Transform[T]) Then(next Transform[T]) Transform[T] {
    return func(x T) T {
        return next(t(x))
    }
}

pipeline := Transform[int](add5).
    Then(double).
    Then(square)

result := pipeline(2)  // ((2+5)*2)² = 196
```

**Middleware chaining:**

```go
type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
    return func(h http.Handler) http.Handler {
        for i := len(middlewares) - 1; i >= 0; i-- {
            h = middlewares[i](h)
        }
        return h
    }
}
```

## Real Example

```go
// Text processing pipeline
type TextProcessor func(string) string

func (tp TextProcessor) Then(next TextProcessor) TextProcessor {
    return func(s string) string {
        return next(tp(s))
    }
}

// Individual processors
toLower := TextProcessor(strings.ToLower)
trim := TextProcessor(strings.TrimSpace)
removePunctuation := TextProcessor(func(s string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsPunct(r) {
            return -1
        }
        return r
    }, s)
})

// Compose pipeline
normalize := toLower.
    Then(trim).
    Then(removePunctuation)

result := normalize("  Hello, World!  ")
// "hello world"
```

## Gotchas

**Type constraints with generics:**

```go
// Works for any types
func Compose[A, B, C any](f func(B) C, g func(A) B) func(A) C

// But limited chaining without same type
add := func(x int) int { ... }
toString := func(x int) string { ... }
// Can't easily chain add -> toString -> anotherInt
```

**Order matters:**

```go
// f(g(x)) - g runs first
Compose(f, g)

// Pipe notation (reverse)
func Pipe[A, B, C any](g func(A) B, f func(B) C) func(A) C {
    return Compose(f, g)
}
Pipe(g, f)  // More intuitive left-to-right
```

**Error handling:**

```go
// Composition breaks with errors
 type Result[T any] func() (T, error)

func ComposeWithError[A, B, C any](
    f func(B) (C, error),
    g func(A) (B, error),
) func(A) (C, error) {
    return func(x A) (C, error) {
        b, err := g(x)
        if err != nil {
            var zero C
            return zero, err
        }
        return f(b)
    }
}
```

## Benefits

- Reusable building blocks
- Declarative code
- Easy testing (test pieces separately)
- Functional style in Go
- Self-documenting pipelines
