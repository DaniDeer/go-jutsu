# Zero Values Are Useful

Go initializes all variables. Zero values are intentionally useful. Design philosophy.

## What It Is (and Isn't)

Every type has zero value. Not null/undefined. Actual usable initial state.

Philosophy: zero value should be useful without initialization. Not all languages do this.

## Where You See It

**Ready-to-use types:**

```go
var mu sync.Mutex  // Ready to Lock() - no New() needed
var buf bytes.Buffer  // Ready to Write() - no initialization
var wg sync.WaitGroup  // Ready to Add() - works immediately
```

**Numeric zero:**

```go
var count int  // 0
var price float64  // 0.0
var enabled bool  // false
```

**Empty collections:**

```go
var slice []int  // nil but can append
var m map[string]int  // nil, need make() to assign
```

## Real Example

```go
type Config struct {
    Timeout time.Duration  // 0 = no timeout
    Retries int           // 0 = no retries
    Debug   bool          // false = production mode
}

// Zero value is sensible default
var cfg Config
// cfg.Timeout = 0, cfg.Retries = 0, cfg.Debug = false

type Counter struct {
    mu    sync.Mutex  // Zero value ready to use!
    count int
}

func (c *Counter) Inc() {
    c.mu.Lock()  // No initialization needed
    defer c.mu.Unlock()
    c.count++
}

// Works immediately
var counter Counter
counter.Inc()  // Just works
```

## Gotchas

**Nil slices vs empty slices:**

```go
var s1 []int        // nil slice
s2 := []int{}       // empty slice
s3 := make([]int, 0) // empty slice

// All have len 0, but different
fmt.Println(s1 == nil)  // true
fmt.Println(s2 == nil)  // false
fmt.Println(s3 == nil)  // false

// But append works on all
s1 = append(s1, 1)  // Works!
```

**Nil maps can't be assigned:**

```go
var m map[string]int  // nil map
m["key"] = 42  // panic: assignment to entry in nil map

// Must make() first
m = make(map[string]int)
m["key"] = 42  // Works
```

**Nil pointers:**

```go
var ptr *int  // nil
*ptr = 42    // panic: nil pointer dereference

// Check before use
if ptr != nil {
    *ptr = 42
}
```

**String zero value:**

```go
var s string  // ""
if s == "" {
    // Common check
}
```

**Interface zero value:**

```go
var i interface{}  // nil
var err error      // nil

// Safe to check
if err != nil { ... }
```

## Design Pattern

**Zero value constructors optional:**

```go
// Old style: requires New()
type Buffer struct {
    data []byte
}

func NewBuffer() *Buffer {
    return &Buffer{data: make([]byte, 0)}
}

// Go style: zero value works
type Builder struct {
    buf bytes.Buffer  // Zero value ready
}

// No New() needed
var b Builder
b.buf.WriteString("works")
```

**Useful zero values:**

- `sync.Mutex` - ready to lock
- `bytes.Buffer` - ready to write
- `strings.Builder` - ready to build
- Slices - can append even if nil
- Pointers - nil check is standard

**Not useful zero values:**

- Maps - can't assign to nil map
- Channels - nil channel blocks forever
- Functions - nil function panics when called
