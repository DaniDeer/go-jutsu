# Bounded Parallelism Pattern

Limit concurrent goroutines using semaphore. Buffered channel as token bucket.

## What It Is (and Isn't)

Max N operations simultaneously. Semaphore with channels. Prevents resource exhaustion.

Not sequential. Not unlimited goroutines. Controls concurrency level.

## Where You See It

**File processing:**

```go
sem := make(chan struct{}, maxConcurrent)
for _, file := range files {
    sem <- struct{}{}  // Acquire
    go func(f string) {
        defer func() { <-sem }()  // Release
        process(f)
    }(file)
}
```

**weighted semaphore:**

```go
import "golang.org/x/sync/semaphore"
sem := semaphore.NewWeighted(10)
sem.Acquire(ctx, 3)  // Acquire 3 tokens
defer sem.Release(3)
```

## Real Example

```go
func ProcessFiles(files []string, maxConcurrent int) {
    sem := make(chan struct{}, maxConcurrent)
    var wg sync.WaitGroup

    for _, file := range files {
        wg.Add(1)
        sem <- struct{}{}  // Block if maxConcurrent reached

        go func(f string) {
            defer wg.Done()
            defer func() { <-sem }()
            processFile(f)
        }(file)
    }

    wg.Wait()
}
```

## Gotchas

**Always release:**

```go
defer func() { <-sem }()  // Use defer
```

**WaitGroup required:**

```go
// Channel only limits concurrency, doesn't wait for completion
var wg sync.WaitGroup
```
