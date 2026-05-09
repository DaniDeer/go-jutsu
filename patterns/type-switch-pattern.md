# Type Switch Pattern

Switch on interface's concrete type. Type assertion in switch form.

## What It Is (and Isn't)

Pattern to handle different types from interface{}. Built on type assertion and switch.

Not reflection (compile-time checks). Not if-else chains. Idiomatic type handling.

## Where You See It

**Basic type switch:**

```go
switch v := val.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
default:
    fmt.Printf("unknown: %T\n", v)
}
```

**Multiple types:**

```go
switch v := val.(type) {
case string, []byte:
    // v is interface{} here
case int, int64:
    // Handle integers
}
```

**Type with method:**

```go
switch v := val.(type) {
case Stringer:
    fmt.Println(v.String())
}
```

## Real Example

```go
func process(val interface{}) {
    switch v := val.(type) {
    case string:
        fmt.Printf("String of length %d: %q\n", len(v), v)
    case int:
        fmt.Printf("Integer: %d\n", v)
    case []int:
        fmt.Printf("Int slice of length %d\n", len(v))
    case error:
        fmt.Printf("Error: %v\n", v)
    case nil:
        fmt.Println("Nil value")
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}
```

## Gotchas

**Type variable scope:**

```go
switch v := val.(type) {
case string:
    // v is string here
case int:
    // v is int here
default:
    // v is interface{} here
}
// v not accessible here
```

**Multiple types lose type info:**

```go
switch v := val.(type) {
case string, []byte:
    // v is interface{} (not string or []byte)
    s := v.(string)  // Need assertion
}
```

**Nil interface vs nil value:**

```go
var i interface{} = (*int)(nil)

switch i.(type) {
case nil:
    // Doesn't match (i is not nil)
case *int:
    // Matches
}
```

## vs Type Assertion

```go
// Type assertion (single type)
if s, ok := val.(string); ok {
    // Use s
}

// Type switch (multiple types)
switch v := val.(type) {
case string:
    // Use v
case int:
    // Use v
}
```
