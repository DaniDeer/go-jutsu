# Empty Struct Signal Pattern

Zero-size type `struct{}` for signaling without memory allocation. Only in Go.

## What It Is (and Isn't)

`struct{}` is zero bytes. Compiler optimizes - takes no memory. Used for signaling, not data.

Not null/nil. Not boolean. Actual type with zero size. Go compiler magic.

## Where You See It

**Channel signals:**

```go
done := make(chan struct{})
go func() {
    work()
    done <- struct{}{}  // Signal complete (0 bytes sent)
}()
<-done  // Wait for signal
```

**Set implementation:**

```go
set := make(map[string]struct{})
set["key"] = struct{}{}  // Only keys matter, value takes 0 bytes
if _, exists := set["key"]; exists { ... }
```

**Method receivers without state:**

```go
type Logger struct{}  // Zero-size type

func (Logger) Log(msg string) {
    fmt.Println(msg)  // Stateless logging
}
```

**Semaphore/worker pool:**

```go
sem := make(chan struct{}, 5)  // Limit to 5 concurrent
go func() {
    sem <- struct{}{}  // Acquire
    defer func() { <-sem }()  // Release
    work()
}()
```

## Real Example

```go
// Coordinate goroutines with zero allocation
func processAll(tasks []Task) {
    done := make(chan struct{})
    errors := make(chan error, len(tasks))

    for _, task := range tasks {
        go func(t Task) {
            if err := t.Process(); err != nil {
                errors <- err
            }
            done <- struct{}{}  // Signal completion
        }(task)
    }

    // Wait for all
    for range tasks {
        <-done
    }
    close(errors)

    // Check errors
    for err := range errors {
        log.Println(err)
    }
}
```

## Gotchas

**Empty struct literal looks weird:**

```go
ch <- struct{}{}  // This is correct
//    ^^^^^^^^ type
//            ^^ value (empty composite literal)
```

**Don't use for data:**

```go
// BAD - confusing
type Config struct{}  // Looks like config but stores nothing

// GOOD - explicit about signaling
type signal struct{}
```

**Not same as nil channel:**

```go
var nilChan chan struct{}       // nil, blocks forever
doneChan := make(chan struct{}) // Real channel, works
```

**Memory:** Zero-size slice still has header:

```go
s := make([]struct{}, 1000000)  // Elements free, but slice header exists
m := make(map[int]struct{})     // Keys cost memory, values don't
```

## Why It Matters

Memory efficient signaling:

```go
// 8 bytes per signal (on 64-bit)
doneBool := make(chan bool)

// 0 bytes per signal
doneSignal := make(chan struct{})
```

For millions of operations, this adds up. Idiomatic Go.
