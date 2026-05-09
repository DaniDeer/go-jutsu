# Worker Pool Pattern

Fixed number of goroutines processing jobs from shared channel. Control concurrency, reuse workers, aggregate results.

## What It Is (and Isn't)

Pool of workers consuming from job queue. Bounded parallelism. Efficient resource usage.

Not unlimited goroutines. Not one goroutine per job. Prevents resource exhaustion.

## Where You See It

**Basic structure:**

```go
jobs := make(chan Job, 100)
results := make(chan Result, 100)

for w := 0; w < numWorkers; w++ {
    go worker(jobs, results)
}

// Send jobs
go func() {
    for _, job := range allJobs {
        jobs <- job
    }
    close(jobs)
}()

// Collect results
```

**HTTP server pattern:**

```go
// net/http uses worker pool internally
http.ListenAndServe(":8080", handler)
```

## Real Example

```go
func ProcessURLs(urls []string, workers int) []Result {
    jobs := make(chan string, len(urls))
    results := make(chan Result, len(urls))

    // Start workers
    for w := 0; w < workers; w++ {
        go func() {
            for url := range jobs {
                results <- fetch(url)
            }
        }()
    }

    // Send jobs
    for _, url := range urls {
        jobs <- url
    }
    close(jobs)

    // Collect
    var out []Result
    for i := 0; i < len(urls); i++ {
        out = append(out, <-results)
    }
    return out
}
```

## Gotchas

**Close jobs channel:**

```go
close(jobs)  // Workers exit when range ends
```

**WaitGroup for completion:**

```go
var wg sync.WaitGroup
wg.Add(numWorkers)
go func() {
    wg.Wait()
    close(results)  // After all workers done
}()
```

**Buffer channels appropriately:**

```go
jobs := make(chan T, queueSize)  // Prevent sender blocking
```
