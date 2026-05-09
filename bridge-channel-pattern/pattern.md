# Bridge Channel Pattern

Flatten channel of channels into single channel. Consume from sequence of channels.

## What It Is (and Isn't)

Transform `chan chan T` into `chan T`. Sequential channel consumption. Simplify nested channels.

Not merge (broadcasting multiple channels simultaneously). Sequential, one at a time.

## Where You See It

**Batched results:**
```go
// Generator produces channels
channelStream := generateChannels()
results := bridge(done, channelStream)
for val := range results {
    // Flat stream of values
}
```

**Pagination:**
```go
pages := fetchPages()  // Each page is channel
allItems := bridge(done, pages)
```

## Real Example

```go
func bridge[T any](done <-chan struct{}, chanStream <-chan <-chan T) <-chan T {
    out := make(chan T)
    
    go func() {
        defer close(out)
        
        for {
            var stream <-chan T
            select {
            case <-done:
                return
            case maybeStream, ok := <-chanStream:
                if !ok {
                    return
                }
                stream = maybeStream
            }
            
            // Drain current stream
            for val := range orDone(done, stream) {
                select {
                case out <- val:
                case <-done:
                    return
                }
            }
        }
    }()
    
    return out
}

// Usage with pagination
func FetchAllPages() <-chan Item {
    pageStream := make(chan <-chan Item)
    
    go func() {
        defer close(pageStream)
        for pageNum := 0; ; pageNum++ {
            page := fetchPage(pageNum)
            pageStream <- page
            if isLastPage(pageNum) {
                break
            }
        }
    }()
    
    return bridge(ctx.Done(), pageStream)
}
```

## Gotchas

**Order matters:**
```go
// Consumes channels sequentially
// Later channels wait for earlier ones
```

**Cancellation propagation:**
```go
// done signal must stop both outer and inner channels
for val := range orDone(done, stream) {  // Critical
```

