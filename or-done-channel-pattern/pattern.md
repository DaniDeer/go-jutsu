# Or-Done Channel Pattern

Wrap channel reads with done signal. Prevent goroutine leaks. Unify cancellation.

## What It Is (and Isn't)

Combine data channel with done channel. Single select. Easier cancellation handling.

Not required for all channels. Not automatic cleanup. Pattern for composable pipelines.

## Where You See It

**Pipeline stage:**
```go
for val := range orDone(done, inputChan) {
    // Process val
    // Automatically exits if done closed
}
```

**instead of:**
```go
for {
    select {
    case <-done:
        return
    case val := <-inputChan:
        // process
    }
}
```

## Real Example

```go
func orDone[T any](done <-chan struct{}, c <-chan T) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for {
            select {
            case <-done:
                return
            case v, ok := <-c:
                if !ok {
                    return
                }
                select {
                case out <- v:
                case <-done:
                    return
                }
            }
        }
    }()
    return out
}

// Usage
for val := range orDone(ctx.Done(), dataChan) {
    process(val)
}
```

## Gotchas

**Two selects needed:**
```go
// Receive
select {
case <-done:
    return
case v := <-c:
}

// Send
select {
case out <- v:
case <-done:
    return
}
```

**Goroutine cost:**
- Adds goroutine per wrapped channel
- Use for long-running pipelines

