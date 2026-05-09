# Panic and Recover Pattern

Handle exceptional failures gracefully. Use `panic()` for unrecoverable errors, `recover()` in deferred functions to regain control.

## What It Is (and Isn't)

Emergency stop mechanism. Unwind stack. Catch with `defer + recover()`.

Not normal error handling. Not for flow control. Last resort for programmer errors.

## Where You See It

**Library protection:**

```go
func SafeHandler(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("panic: %v", r)
            http.Error(w, "Internal Error", 500)
        }
    }()
    // handler code that might panic
}
```

**Index bounds:**

```go
if i >= len(slice) {
    panic("index out of range")  // Programmer error
}
```

## Real Example

```go
func ProcessBatch(items []Item) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()

    for _, item := range items {
        process(item)  // Might panic
    }
    return nil
}
```

## Gotchas

**Deferred recover only:**

```go
recover()  // Does nothing
defer recover()  // Does nothing
defer func() { recover() }()  // Works!
```

**When to panic:**

- Initialization failures (e.g., regex compile)
- Impossible scenarios (`default:` in exhaustive switch)
- Programmer errors (nil dereference)

**Don't panic for:**

- Expected errors (file not found → return error)
- Validation failures
- Network errors
