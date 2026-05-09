# Error Wrapping Pattern

Add context to errors while preserving original. Use `%w` verb and `errors.Unwrap/Is/As`.

## What It Is (and Isn't)

Wrap errors with context. Build error chains. Preserve type information.

Not losing original error. Not string concatenation. Go 1.13+ feature.

## Where You See It

**Basic wrapping:**

```go
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}
```

**Check wrapped:**

```go
if errors.Is(err, io.EOF) {  // Works through wrapping

if var pathErr *os.PathError; errors.As(err, &pathErr) {
    // Extract wrapped type
}
```

## Real Example

```go
func ProcessFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("ProcessFile: %w", err)
    }

    if err := validate(data); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    return nil
}

// Caller can check original error
if errors.Is(err, os.ErrNotExist) {
    log.Println("File doesn't exist")
}
```

## Gotchas

**%w vs %v:**

```go
fmt.Errorf("context: %v", err)  // String only
fmt.Errorf("context: %w", err)  // Wraps (use this)
```

**Multiple wraps:**

```go
err1 := errors.New("original")
err2 := fmt.Errorf("layer 2: %w", err1)
err3 := fmt.Errorf("layer 3: %w", err2)
errors.Is(err3, err1)  // true
```
