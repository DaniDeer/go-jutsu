# Go Closure Pattern (Manual "Currying")

## Not True Currying

Go doesn't have automatic currying like Haskell/ML. This is **manual closure pattern** that looks similar.

Real currying transforms:

```haskell
-- Haskell: automatic
add x y = x + y
add 5     -- returns function waiting for y
```

Go requires explicit returns:

```go
// Go: manual
func add(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}
add(5) // returns function waiting for y
```

## Where You See It

**Middleware** (most common):

```go
func cors(allowedOrigins []string) func(http.Handler) http.Handler
func rateLimit(reqPerMin int) func(http.Handler) http.Handler
```

**Option builders:**

```go
func WithTimeout(d time.Duration) func(*Client)
func WithRetries(n int) func(*Client)
```

**Factory functions:**

```go
func NewValidator(rules []Rule) func(string) error
func MakeFilter(predicate func(int) bool) func([]int) []int
```

## Power Move

Chain multiple middleware:

```go
Handler: cors(origins)(rateLimit(100)(auth(apiKey)(mux)))
```

Or cleaner with helpers:

```go
Handler: chain(mux,
    cors(origins),
    rateLimit(100),
    auth(apiKey),
)
```

## Real Example

```go
func requestLogger(logger *log.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            next.ServeHTTP(w, r)
            logger.Printf("Request: %s %s", r.Method, r.URL.Path)
        })
    }
}

// Usage: configure at setup time, execute at request time
srv := &http.Server{
    Handler: requestLogger(logger)(mux),
}
```

Pattern lets you configure behavior at setup time, execute at request time. Separation of concerns.
