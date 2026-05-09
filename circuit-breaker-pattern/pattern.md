# Circuit Breaker Pattern

Prevent cascading failures. Three states: Closed (normal), Open (failing), Half-Open (testing recovery).

## What It Is (and Isn't)

Fail fast on broken dependencies. Auto-recovery testing. Stops hammering failing service.

Not retry logic. Not rate limiting. Pattern for service resilience.

## Where You See It

**HTTP client:**

```go
cb := NewCircuitBreaker(5, time.Minute)
if cb.IsOpen() {
    return ErrCircuitOpen
}
err := cb.Call(func() error {
    return httpGet(url)
})
```

**Database calls:**

```go
if breaker.AllowRequest() {
    err := db.Query(...)
    breaker.RecordResult(err)
}
```

## Real Example

```go
type State int
const (
    Closed State = iota
    Open
    HalfOpen
)

type CircuitBreaker struct {
    state        State
    failures     int
    threshold    int
    timeout      time.Duration
    lastFailTime time.Time
    mu           sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.state == Open {
        if time.Since(cb.lastFailTime) > cb.timeout {
            cb.state = HalfOpen
        } else {
            return ErrCircuitOpen
        }
    }

    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()
        if cb.failures >= cb.threshold {
            cb.state = Open
        }
    } else {
        cb.state = Closed
        cb.failures = 0
    }
    return err
}
```

## Gotchas

**Don't confuse with retry:**

- Retry: try again immediately
- Circuit breaker: stop trying if repeatedly failing

**Threshold tuning:**

```go
threshold := 5       // Failures before opening
timeout := 1*time.Minute  // Wait before testing recovery
```
