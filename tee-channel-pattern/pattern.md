# Tee Channel Pattern

Duplicate values from one channel to multiple outputs. Like Unix `tee` command.

## What It Is (and Isn't)

Split channel into N outputs. All receivers get same values. Fan-out pattern.

Not fan-out/fan-in. Not merge. One-to-many duplication.

## Where You See It

**Logging + processing:**
```go
data, log := tee(done, input)
go logger(log)
process(data)
```

**Multiple consumers:**
```go
out1, out2 := tee(done, source)
go consumer1(out1)
go consumer2(out2)
```

## Real Example

```go
func tee[T any](done <-chan struct{}, in <-chan T) (<-chan T, <-chan T) {
    out1 := make(chan T)
    out2 := make(chan T)
    
    go func() {
        defer close(out1)
        defer close(out2)
        
        for val := range orDone(done, in) {
            var out1, out2 = out1, out2  // Shadow for select
            for i := 0; i < 2; i++ {
                select {
                case <-done:
                    return
                case out1 <- val:
                    out1 = nil  // Disable this case
                case out2 <- val:
                    out2 = nil
                }
            }
        }
    }()
    
    return out1, out2
}
```

## Gotchas

**Slowest consumer wins:**
```go
// If out1 blocks, out2 also blocks
// Both must consume at similar rates
```

**Unbuffered channels:**
```go
// Sends blocked until both receivers ready
// Consider buffering if consumers differ in speed
```

**Nil channel trick:**
```go
out1 <- val
out1 = nil  // Disable in select
```

