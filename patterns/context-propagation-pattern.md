# Context Propagation Pattern

Request-scoped cancellation, deadlines, and values propagate through call stack.

## What It Is (and Isn't)

`context.Context` carries cancellation signals, deadlines, request-scoped values. First param by convention.

Not for passing business logic data. For request lifetime info only. Standard library pattern.

## Where You See It

**HTTP handlers:**

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // Request context
    result, err := queryDB(ctx, "SELECT...")
}
```

**Cancellation:**

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go worker(ctx)
// cancel() stops worker
```

**Timeout:**

```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

result, err := slowOperation(ctx)
```

## Real Example

```go
func fetchData(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}

func handler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    data, err := fetchData(ctx, "https://api.example.com")
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            http.Error(w, "timeout", http.StatusGatewayTimeout)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write(data)
}
```

## Gotchas

**Always check ctx.Done():**

```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-ch:
    return process(result)
}
```

**Context values (avoid):**

```go
// BAD: passing business data
ctx = context.WithValue(ctx, "user", user)

// GOOD: pass explicitly
func doThing(ctx context.Context, user User) {...}
```

**Cancel parent cancels children:**

```go
parent, cancel := context.WithCancel(context.Background())
child, _ := context.WithTimeout(parent, time.Hour)

cancel()  // Cancels both parent and child
```

**Always defer cancel():**

```go
ctx, cancel := context.WithCancel(ctx)
defer cancel()  // Prevent leak
```
