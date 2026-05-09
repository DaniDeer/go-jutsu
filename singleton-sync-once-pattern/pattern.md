# Singleton with sync.Once

Lazy initialization of singleton instance. Thread-safe, executed exactly once.

## What It Is (and Isn't)

One instance globally. `sync.Once.Do()` ensures single execution. Lazy initialization.

Not eager loading. Not package-level var (that's immediate init). Pattern for expensive setup.

## Where You See It

**Database connection:**

```go
var (
    db   *sql.DB
    once sync.Once
)

func GetDB() *sql.DB {
    once.Do(func() {
        db, _ = sql.Open("postgres", dsn)
    })
    return db
}
```

**Config loading:**

```go
var (
    config Config
    once   sync.Once
)

func GetConfig() Config {
    once.Do(func() {
        config = loadConfig()
    })
    return config
}
```

## Real Example

```go
type Logger struct {
    file *os.File
}

var (
    logger *Logger
    once   sync.Once
)

func GetLogger() *Logger {
    once.Do(func() {
        file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND, 0644)
        if err != nil {
            panic(err)
        }
        logger = &Logger{file: file}
    })
    return logger
}
```

## Gotchas

**Error handling:**

```go
var (
    instance *T
    initErr  error
    once     sync.Once
)

func Get() (*T, error) {
    once.Do(func() {
        instance, initErr = initialize()
    })
    return instance, initErr
}
```

**sync.Once never resets:**

```go
once.Do(fn1)  // Runs
once.Do(fn2)  // NEVER runs
```
