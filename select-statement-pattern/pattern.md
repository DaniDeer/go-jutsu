# Select Statement Pattern

Non-deterministic channel multiplexing. Choose from multiple channel operations.

## What It Is (and Isn't)

`select` waits on multiple channels simultaneously. First ready case wins. Non-deterministic if multiple ready.

Not if-else. Not switch. Specifically for channel operations. Unique to Go.

## Where You See It

**Timeout pattern:**

```go
select {
case result := <-ch:
    return result
case <-time.After(5 * time.Second):
    return ErrTimeout
}
```

**Non-blocking channel ops:**

```go
select {
case msg := <-ch:
    process(msg)
default:
    // Doesn't block
}
```

**Multiple channels:**

```go
select {
case msg := <-ch1:
    handle1(msg)
case msg := <-ch2:
    handle2(msg)
case <-quit:
    return
}
```

## Real Example

```go
func worker(jobs <-chan Job, results chan<- Result, done <-chan struct{}) {
    for {
        select {
        case job := <-jobs:
            results <- process(job)

        case <-done:
            log.Println("Worker stopping")
            return

        case <-time.After(30 * time.Second):
            log.Println("Worker idle for 30s")
        }
    }
}

func fanIn(ch1, ch2 <-chan string) <-chan string {
    out := make(chan string)
    go func() {
        for {
            select {
            case msg := <-ch1:
                out <- msg
            case msg := <-ch2:
                out <- msg
            }
        }
    }()
    return out
}
```

## Gotchas

**Random selection when multiple ready:**

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)
ch1 <- 1
ch2 <- 2

select {
case v := <-ch1:
    fmt.Println(v)   // Might print
case v := <-ch2:
    fmt.Println(v)   // Or might print
}
// Non-deterministic!
```

**Default makes it non-blocking:**

```go
select {
case msg := <-ch:
    process(msg)
default:
    // Runs immediately if ch not ready
}
```

**Empty select blocks forever:**

```go
select {}  // Blocks forever (deadlock)
```

**Nil channel blocks forever:**

```go
var ch chan int  // nil
select {
case <-ch:  // Never executes
    // ...
}
```

**time.After creates goroutine:**

```go
// BAD: creates goroutine every iteration
for {
    select {
    case <-time.After(1 * time.Second):  // Leak!
    }
}

// GOOD: reuse ticker
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for {
    select {
    case <-ticker.C:
        // ...
    }
}
```

**Send and receive:**

```go
select {
case msg := <-inbox:   // Receive
    process(msg)
case outbox <- result:  // Send
    log.Println("sent")
}
```

## Common Patterns

**Timeout:**

```go
select {
case result := <-ch:
    return result, nil
case <-time.After(timeout):
    return nil, ErrTimeout
}
```

**Context cancellation:**

```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-ch:
    return result, nil
}
```

**Priority select (bias):**

```go
// Check high priority first
select {
case msg := <-highPriority:
    process(msg)
default:
    // Then check normal
    select {
    case msg := <-normal:
        process(msg)
    default:
        // Neither ready
    }
}
```
