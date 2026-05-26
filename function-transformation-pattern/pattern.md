# Function Transformation Pattern

A function takes another function as input and returns a *new*, enhanced version of it.

## What It Is (and Isn't)

Transforms behavior — wraps a function to add cross-cutting concerns without touching its logic.

Not the same as function composition (`f(g(x))`). Composition chains *outputs*. Transformation wraps *behavior* around an existing function.

Not a closure capturing data — transformers capture *functions*.

## Where You See It

- **Retry logic**: wrap any fallible func to auto-retry on failure
- **Timing/tracing**: wrap any func to measure and log duration
- **Memoization**: wrap a pure func to cache results by input
- **Auth guards**: wrap a handler to check permissions first

## Real Example

```go
// withRetry wraps fn to retry up to n times on error
func withRetry(n int, fn func() error) func() error {
    return func() error {
        var err error
        for range n {
            if err = fn(); err == nil {
                return nil
            }
        }
        return err
    }
}

// withTiming wraps fn to print execution duration
func withTiming(name string, fn func()) func() {
    return func() {
        start := time.Now()
        fn()
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}

// Original logic is untouched — behavior is layered on top
fetchWithRetry := withRetry(3, fetchData)
timedFetch     := withTiming("fetch", fetchData)
```

## Gotchas

**Signature coupling**: the transformer must match the exact function signature. For reusability across types, use generics (Go 1.18+):

```go
func memoize[K comparable, V any](fn func(K) V) func(K) V {
    cache := make(map[K]V)
    return func(k K) V {
        if v, ok := cache[k]; ok {
            return v
        }
        v := fn(k)
        cache[k] = v
        return v
    }
}
```

**Not goroutine-safe by default**: the memoize cache above has no locking. Add `sync.Mutex` if called concurrently.
