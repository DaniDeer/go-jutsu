# Factory Pattern

Registry of constructors. Create objects by name/type. Extensible object creation.

## What It Is (and Isn't)

Map type names to constructors. Plugin-style creation. Runtime type selection.

Not builder pattern. Not always needed (simple constructors work). Pattern for extensibility.

## Where You See It

**Plugin system:**
```go
var handlers = map[string]func() Handler{}

func RegisterHandler(name string, fn func() Handler) {
    handlers[name] = fn
}

func CreateHandler(name string) Handler {
    return handlers[name]()
}
```

**database/sql:**
```go
sql.Register("postgres", &pq.Driver{})  // Factory registration
db, _ := sql.Open("postgres", dsn)      // Factory usage
```

## Real Example

```go
type Logger interface {
    Log(msg string)
}

type LoggerFactory func(config Config) Logger

var loggerRegistry = make(map[string]LoggerFactory)

func RegisterLogger(name string, factory LoggerFactory) {
    loggerRegistry[name] = factory
}

func NewLogger(name string, cfg Config) (Logger, error) {
    factory, ok := loggerRegistry[name]
    if !ok {
        return nil, fmt.Errorf("unknown logger: %s", name)
    }
    return factory(cfg), nil
}

// Usage
func init() {
    RegisterLogger("stdout", func(cfg Config) Logger {
        return &StdoutLogger{}
    })
    RegisterLogger("file", func(cfg Config) Logger {
        return &FileLogger{path: cfg.Path}
    })
}
```

## Gotchas

**init() for registration:**
```go
func init() {
    RegisterType("mytype", factory)
}
```

**Thread-safe registry:**
```go
var (
    mu       sync.RWMutex
    registry = make(map[string]Factory)
)
```

