# Method Values vs Expressions Pattern

Methods can be values (closures) or expressions (functions). Two ways to reference methods.

## What It Is (and Isn't)

Method value: `obj.Method` - bound to instance (closure).
Method expression: `(*Type).Method` - unbound, takes receiver as first param.

Not obvious distinction. Subtle but powerful for callbacks.

## Where You See It

**Method value (bound):**

```go
s := &Server{}
callback := s.Handle  // Closure over s
callback(req)  // Calls s.Handle(req)
```

**Method expression (unbound):**

```go
handler := (*Server).Handle  // Function needing receiver
handler(s, req)  // Calls s.Handle(req)
```

**Callbacks:**

```go
button.OnClick(s.handleClick)  // Method value
```

## Real Example

```go
type Counter struct {
    count int
}

func (c *Counter) Inc() {
    c.count++
}

// Method value - bound to instance
c := &Counter{}
inc := c.Inc  // Closure
inc()  // Calls c.Inc()
fmt.Println(c.count)  // 1

// Method expression - unbound
incFunc := (*Counter).Inc  // func(*Counter)
incFunc(c)  // Needs receiver
fmt.Println(c.count)  // 2
```

## Gotchas

**Method value captures receiver:**

```go
for _, item := range items {
    go item.Process()  // Each goroutine gets its own item
}

// vs

processFunc := Item.Process  // Unbound
for _, item := range items {
    go processFunc(item)  // Same, but explicit
}
```

**Type matters:**

```go
type T struct{}
func (t *T) Method() {}  // Pointer receiver

// Method value
var t T
f := t.Method  // Works (Go takes address)

// Method expression
g := (*T).Method  // func(*T)
// h := T.Method  // ✗ Compile error
```

**Callbacks with method values:**

```go
type Handler struct{}

func (h *Handler) Handle(req Request) {
    // ...
}

// Register callback
http.HandleFunc("/", h.Handle)  // Method value
```

## Use Cases

- Callbacks without closures
- Delayed execution
- Passing methods as first-class functions
- Map/filter/reduce patterns
