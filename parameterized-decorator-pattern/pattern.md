# Parameterized Decorator Pattern

A **decorator factory** — a function that accepts configuration and returns a decorator — enabling parameterized, composable behavior wrapping.

## What It Is (and Isn't)

A plain decorator takes a function and returns a wrapped version:

```go
func toUppercase(next Processor) Processor { ... }  // no config
```

A **parameterized decorator** (decorator factory) takes config *first*, returns a decorator:

```go
func truncate(n int) func(Processor) Processor { ... }  // config → decorator
```

This is the Go equivalent of Python's `@get_truncate(9)` syntax — the call `get_truncate(9)` *produces* the decorator that Python then applies.

Not the same as [Decorator Pattern](../decorator-pattern/) — that covers plain (zero-parameter) decorators and interface decorators. Not the same as [Closure Currying Pattern](../go-closure-currying-pattern/) — that covers data currying (`add(x)(y)`). This pattern is specifically about *configuring decorator behavior at setup time*.

## Where You See It

- **HTTP middleware with config**: `withTimeout(5*time.Second)`, `withRateLimit(100)`, `cors(allowedOrigins)`
- **Logging with level**: `withLogLevel(slog.LevelDebug)` wraps a handler and filters by level
- **Retry with count**: `withRetry(3)` wraps a handler and retries on error
- **Auth with key**: `withAuth(apiKey)` wraps a handler and validates the token

## Real Example

```go
// In Go, HTTP middleware is the canonical home for parameterized decorators.

type Middleware func(http.Handler) http.Handler

// Plain decorator — no config needed
func withLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// Decorator factory — config captured at setup time
func withRateLimit(reqPerMin int) Middleware {
    limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(reqPerMin)), 1)
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "too many requests", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// Stack: outermost runs first — logging wraps rate limiting wraps the mux
handler := withLogging(withRateLimit(100)(mux))
```

## Python → Go Translation

```python
# Python
@to_uppercase           # plain decorator
@get_truncate(9)        # factory: call produces the decorator
def print_input(s: str) -> None:
    print(s)
```

```go
// Go — explicit, no syntax sugar
type Processor func(string)

func toUppercase(next Processor) Processor {           // plain
    return func(s string) { next(strings.ToUpper(s)) }
}

func truncate(n int) func(Processor) Processor {       // factory
    return func(next Processor) Processor {
        return func(s string) { next(s[:min(n, len(s))]) }
    }
}

printInput := toUppercase(truncate(9)(func(s string) { fmt.Println(s) }))
printInput("Keep Calm and Carry On")  // → "KEEP CALM"
```

Use a `chain` helper when stacking more than two decorators:

```go
func chain(p Processor, decorators ...func(Processor) Processor) Processor {
    for i := len(decorators) - 1; i >= 0; i-- {
        p = decorators[i](p)
    }
    return p
}

// Equivalent to toUppercase(truncate(9)(printLine))
p := chain(printLine, toUppercase, truncate(9))
```

## Gotchas

**Stacking order matters**: decorators apply bottom-up (innermost first). `toUppercase(truncate(9)(f))` truncates first, then uppercases. Reverse the order and you truncate already-uppercased text — same result here, but order is critical when decorators depend on each other's output.

**Signature lock-in**: factory-produced decorators must match the exact function signature. Use generics for reusable factories across types.

**Config captured at factory call time**: the limiter in `withRateLimit(100)` is created once and shared across all requests — intentional here, but beware when config holds mutable state.
