# Closure State Pattern

A closure captures a variable from its enclosing scope, giving the returned function *persistent memory* between calls.

## What It Is (and Isn't)

The returned function reads and/or mutates variables from the outer scope. Those variables live as long as the closure does.

Not a struct with methods — closures are lighter and anonymous, can't be extended. Use a struct when state is complex or needs to be serialized.

Not global state — each call to the outer function creates *independent* state. Two closures from the same factory don't share a counter.

## Where You See It

- **ID generators**: auto-incrementing unique identifiers scoped to a component
- **Request counters**: per-server stats middleware that doesn't need a global
- **In-memory caches**: key-value store tied to a specific lifecycle
- **Accumulators**: build up a collection over multiple calls

## Real Example

```go
// idGenerator returns a function that produces unique prefixed IDs
func idGenerator(prefix string) func() string {
    n := 0                          // captured state
    return func() string {
        n++
        return fmt.Sprintf("%s-%04d", prefix, n)
    }
}

newUserID  := idGenerator("user")
newOrderID := idGenerator("order")

fmt.Println(newUserID())  // user-0001
fmt.Println(newUserID())  // user-0002
fmt.Println(newOrderID()) // order-0001  ← independent state
```

## Gotchas

**Not goroutine-safe**: the captured variable is shared across concurrent callers. Protect with a mutex:

```go
func safeCounter() func() int {
    var mu sync.Mutex
    n := 0
    return func() int {
        mu.Lock()
        defer mu.Unlock()
        n++
        return n
    }
}
```

**Shared vs independent state**: call the outer function *once* if you want shared state, pass the result around. Call it multiple times if you want independent counters.
