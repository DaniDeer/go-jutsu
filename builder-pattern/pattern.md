# Builder Pattern

Construct complex objects step-by-step. Mutable builder, immutable result. Alternative to long constructors.

## What It Is (and Isn't)

Fluent API for object construction. Optional fields without nil pointers. Clear initialization.

Not required for simple structs. Not Go idiomatic (prefer functional options for libraries).

## Where You See It

**HTTP request builder:**

```go
req := NewRequestBuilder().
    Method("POST").
    URL("/api/users").
    Header("Auth", token).
    Body(data).
    Build()
```

**Query builder:**

```go
query := NewQueryBuilder().
    Select("name", "age").
    From("users").
    Where("age > ?", 18).
    Build()
```

## Real Example

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
    logger  Logger
}

type ServerBuilder struct {
    server Server
}

func NewServerBuilder() *ServerBuilder {
    return &ServerBuilder{
        server: Server{
            host:    "localhost",
            port:    8080,
            timeout: 30 * time.Second,
        },
    }
}

func (b *ServerBuilder) Host(host string) *ServerBuilder {
    b.server.host = host
    return b
}

func (b *ServerBuilder) Build() Server {
    return b.server  // Return copy (immutable)
}
```

## Gotchas

**Return pointer for chaining:**

```go
func (b *Builder) Field(v T) *Builder {  // *Builder not Builder
    return b
}
```

**Validation in Build():**

```go
func (b *Builder) Build() (Server, error) {
    if b.server.port == 0 {
        return Server{}, errors.New("port required")
    }
    return b.server, nil
}
```
