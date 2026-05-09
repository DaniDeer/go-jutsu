# Rate Limiting Pattern

Control request rate. Token bucket algorithm. Use `golang.org/x/time/rate` or channels.

## What It Is (and Isn't)

Throttle operations to max rate. Prevent resource exhaustion. Smooth traffic spikes.

Not circuit breaking. Not backpressure. Pattern for protecting resources.

## Where You See It

**API rate limiter:**

```go
limiter := rate.NewLimiter(10, 1)  // 10 req/sec, burst 1
if !limiter.Allow() {
    return ErrRateLimited
}
```

**Channel-based:**

```go
ticker := time.NewTicker(100 * time.Millisecond)
for range ticker.C {
    processNext()  // Max 10/sec
}
```

## Real Example

```go
import "golang.org/x/time/rate"

type RateLimitedAPI struct {
    limiter *rate.Limiter
}

func NewAPI(rps int) *RateLimitedAPI {
    return &RateLimitedAPI{
        limiter: rate.NewLimiter(rate.Limit(rps), rps),
    }
}

func (api *RateLimitedAPI) Call(ctx context.Context) error {
    if err := api.limiter.Wait(ctx); err != nil {
        return err
    }
    // Make API call
    return nil
}
```

## Gotchas

**Allow vs Wait:**

```go
limiter.Allow()  // Returns false immediately if no tokens
limiter.Wait(ctx)  // Blocks until token available
```

**Burst parameter:**

```go
rate.NewLimiter(10, 5)  // 10/sec, allow burst of 5
```

**Per-client limiting:**

```go
var limiters = make(map[string]*rate.Limiter)
mu.Lock()
limiter := limiters[clientID]
mu.Unlock()
```
