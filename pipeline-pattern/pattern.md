# Pipeline Pattern

Channel-based stream processing. Stages connected via channels. Concurrent and composable.

## What It Is (and Isn't)

Functions that receive/send on channels. Chain stages for data flow. Go concurrency pattern.

Not batch processing. Not sequential loops. Concurrent by design.

## Where You See It

**Basic pipeline:**

```go
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// Usage
nums := generator(1, 2, 3, 4)
squares := square(nums)
for n := range squares {
    fmt.Println(n)  // 1, 4, 9, 16
}
```

**Fan-out (parallel processing):**

```go
func fanOut(in <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        outputs[i] = worker(in)
    }
    return outputs
}
```

**Fan-in (merge channels):**

```go
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for n := range c {
                out <- n
            }
        }(ch)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

## Real Example

```go
// Image processing pipeline
func loadImages(paths []string) <-chan Image {
    out := make(chan Image)
    go func() {
        defer close(out)
        for _, path := range paths {
            img := load(path)
            out <- img
        }
    }()
    return out
}

func resize(in <-chan Image) <-chan Image {
    out := make(chan Image)
    go func() {
        defer close(out)
        for img := range in {
            out <- img.Resize(800, 600)
        }
    }()
    return out
}

func compress(in <-chan Image) <-chan Image {
    out := make(chan Image)
    go func() {
        defer close(out)
        for img := range in {
            out <- img.CompressJPEG(85)
        }
    }()
    return out
}

// Build pipeline
images := loadImages(paths)
resized := resize(images)
compressed := compress(resized)

for img := range compressed {
    save(img)
}
```

## Gotchas

**Always close channels:**

```go
out := make(chan int)
go func() {
    defer close(out)  // Critical!
    for _, n := range nums {
        out <- n
    }
}()
```

**Done signal for cancellation:**

```go
func stage(done <-chan struct{}, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n:
            case <-done:
                return
            }
        }
    }()
    return out
}
```

**Buffered channels for throughput:**

```go
out := make(chan int, 100)  // Buffer reduces blocking
```

## Benefits

- Concurrent processing
- Composable stages
- Memory efficient (streaming)
- Easy to parallelize (fan-out)
- Cancellable (with context)
