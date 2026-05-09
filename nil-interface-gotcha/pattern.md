# Nil Interface Gotcha

Non-nil interface holding nil concrete value is not nil. Classic Go trap.

## What It Is (and Isn't)

Interface holds two things: type and value. Interface is nil only if **both** are nil.

If type is set, interface isn't nil—even if value is nil pointer. Not intuitive.

## Where You See It

**Error handling gone wrong:**

```go
func doThing() error {
    var err *MyError = nil
    return err  // Returns non-nil error!
}

if err := doThing(); err != nil {
    // This runs even though we returned nil!
}
```

**Checking for nil:**

```go
var val *int = nil
var iface interface{} = val
fmt.Println(iface == nil)  // false!
```

**Method returns:**

```go
func getWriter() io.Writer {
    var w *os.File = nil
    return w  // Non-nil io.Writer
}
```

## Real Example

```go
type MyError struct {
    Msg string
}

func (e *MyError) Error() string {
    return e.Msg
}

// BAD: Returns non-nil error
func badFunc() error {
    var err *MyError = nil
    // Some logic...
    return err  // error interface contains (*MyError, nil)
}

// GOOD: Return typed nil
func goodFunc() error {
    var err error = nil  // Or just: return nil
    // Some logic...
    return err
}

// GOOD: Check before return
func betterFunc() error {
    var err *MyError = nil
    // Some logic...
    if err != nil {
        return err
    }
    return nil  // Return untyped nil
}
```

## Gotchas

**Why this happens:**

```go
// Interface structure (conceptual)
type interface {
    type  *Type   // Concrete type
    value unsafe.Pointer  // Concrete value
}

// nil interface
var err error = nil
// err = {type: nil, value: nil}

// non-nil interface with nil value
var myErr *MyError = nil
var err error = myErr
// err = {type: *MyError, value: nil}
```

**Function returns:**

```go
func process() error {
    var err *MyError  // nil pointer
    if badThing {
        err = &MyError{"failed"}
    }
    return err  // Bug if err is nil!
}

// Fix:
func process() error {
    if badThing {
        return &MyError{"failed"}
    }
    return nil  // Explicit nil
}
```

**Interface comparison:**

```go
var a interface{} = (*int)(nil)
var b interface{} = (*int)(nil)
fmt.Println(a == b)        // true (same type and value)
fmt.Println(a == nil)      // false (type is set)

var c interface{} = nil
fmt.Println(c == nil)      // true (both nil)
```

## How to Avoid

1. Return explicit `nil` for error/interface types
2. Check `!= nil` before returning pointer as interface
3. Use typed `nil` correctly: `var err error = nil`
4. Remember: interface holds (type, value) pair
