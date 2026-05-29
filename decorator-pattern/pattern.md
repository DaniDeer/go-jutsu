# Decorator Pattern

Wrap a function or interface to add behavior — without changing the original.

## What It Is (and Isn't)

Two forms share the same idea: intercept a call, do something extra, delegate to the original.

**Function decorator** — wraps a `func` value. The Go equivalent of Python's `@decorator` syntax:

```python
# Python
@log_duration
def fetch(url): ...
```

```go
// Go — explicit, no syntax sugar
fetch = logDuration("fetch", fetch)
```

**Interface decorator** — wraps a struct that implements an interface, adding behavior to every method transparently. Classic GoF Decorator pattern; `bufio.NewReader` is the standard library example.

Not the same as [Function Transformation Pattern](../function-transformation-pattern/) — that covers
specific transformations (retry, memoize). The Decorator Pattern is the broader concept either form
belongs to.

Not middleware (though middleware *is* a decorator applied to `http.Handler`).

## Where You See It

- **Function decorators**: logging, timing, auth checks, rate-limiting wrappers around any `func`
- **Interface decorators**: `bufio.NewReader` wraps `io.Reader`; `gzip.NewWriter` wraps `io.Writer`; `tls.Client` wraps `net.Conn`
- **HTTP middleware**: `loggingMiddleware(http.Handler) http.Handler` — industry standard pattern

## Real Examples

### Function Decorator

```go
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// withLogging decorates any HandlerFunc with request logging
func withLogging(next HandlerFunc) HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    }
}

// withAuth decorates any HandlerFunc with a token check
func withAuth(next HandlerFunc) HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") == "" {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}

// Stack decorators — outermost runs first
handler := withLogging(withAuth(serveAPI))
```

### Interface Decorator

```go
// countingReader wraps io.Reader and tracks total bytes read
type countingReader struct {
    r     io.Reader
    total int64
}

func (c *countingReader) Read(p []byte) (n int, err error) {
    n, err = c.r.Read(p)
    c.total += int64(n)
    return
}

// NewCountingReader decorates any io.Reader
func NewCountingReader(r io.Reader) *countingReader {
    return &countingReader{r: r}
}

// Usage: drop-in replacement anywhere io.Reader is accepted
cr := NewCountingReader(resp.Body)
io.Copy(dst, cr)
fmt.Printf("read %d bytes\n", cr.total)
```

## Gotchas

**Signature lock-in** (function decorator): the wrapper must exactly match the wrapped function's
signature. Use interfaces or generics when you need a decorator that works across many types.

**Transparency breaks on concrete return types** (interface decorator): if the wrapped type exposes
methods beyond the interface, callers lose access to them after wrapping. Return the concrete wrapper
type when callers need those extra methods (e.g., `countingReader.total`).

**Stack order matters**: `withLogging(withAuth(h))` logs *then* checks auth (auth may short-circuit).
`withAuth(withLogging(h))` checks auth first, only logs authenticated requests.
