# Functional Options Pattern

Configure structs with optional parameters using functions. Extensible, clean APIs.

## What It Is (and Isn't)

Pattern by Dave Cheney. Functions return functions that configure struct. No builder pattern needed.

Not method chaining. Not config structs. Idiomatic Go for optional config.

## Where You See It

**Client configuration:**

```go
client := NewClient(
    WithTimeout(30*time.Second),
    WithRetries(3),
    WithLogger(log.Default()),
)
```

**Server options:**

```go
server := NewServer(
    addr,
    WithTLS(cert, key),
    WithMaxConns(100),
)
```

## Real Example

```go
type Server struct {
    addr    string
    timeout time.Duration
    maxConns int
    tls     *TLSConfig
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithMaxConns(n int) Option {
    return func(s *Server) {
        s.maxConns = n
    }
}

func WithTLS(cert, key string) Option {
    return func(s *Server) {
        s.tls = &TLSConfig{cert, key}
    }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{
        addr:     addr,
        timeout:  30 * time.Second,  // Default
        maxConns: 100,                // Default
    }

    for _, opt := range opts {
        opt(s)
    }

    return s
}

// Usage
server := NewServer(":8080",
    WithTimeout(60*time.Second),
    WithMaxConns(200),
)
```

## Gotchas

**Variadic must be last:**

```go
// GOOD
func New(required string, opts ...Option) *Thing

// BAD
func New(opts ...Option, required string) *Thing  // won't compile
```

**Options can conflict:**

```go
// Last wins
s := NewServer(":8080",
    WithTimeout(10*time.Second),
    WithTimeout(20*time.Second),  // Overwrites
)
```

**Validation:**

```go
func NewServer(addr string, opts ...Option) (*Server, error) {
    s := &Server{addr: addr}
    for _, opt := range opts {
        opt(s)
    }

    if s.maxConns < 1 {
        return nil, errors.New("maxConns must be positive")
    }

    return s, nil
}
```

## Benefits

- Backward compatible (add options without breaking API)
- Self-documenting (`WithTimeout` vs `timeout time.Duration`)
- Optional parameters without nil checks
- Flexible ordering
- No config struct pollution
