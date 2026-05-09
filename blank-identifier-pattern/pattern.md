# Blank Identifier Pattern

`_` ignores values, imports for side effects, enforces interface satisfaction.

## What It Is (and Isn't)

Underscore `_` is write-only variable. Compiler discards. Multiple uses.

Not variable. Not accessible after assignment. Special identifier in Go.

## Where You See It

**Ignore return values:**

```go
_, err := doThing()  // Ignore first return
val, _ := map[key]   // Ignore ok boolean
```

**Import for side effects:**

```go
import _ "database/sql/driver"  // Runs init() only
```

**Compile-time interface check:**

```go
var _ io.Writer = (*MyType)(nil)  // Verify satisfaction
```

**Loop variable:**

```go
for _, val := range slice {  // Ignore index
for key := range map {  // Ignore value (implicit)
```

## Real Example

```go
// Import database driver for side effects
import (
    "database/sql"
    _ "github.com/lib/pq"  // Registers driver
)

// Ignore unused variables during development
func debug() {
    x := expensive()
    _ = x  // Silence "declared and not used" error
}

// Verify interface implementation
var _ http.Handler = (*MyHandler)(nil)

// Ignore index in range
for _, user := range users {
    process(user)
}

// Ignore error (not recommended!)
file, _ := os.Open("file.txt")
```

## Gotchas

**Multiple return values:**

```go
// Can ignore any position
v1, _, v3 := multiReturn()
_, v2, _ := multiReturn()
```

**Import side effects:**

```go
// Common for database drivers
import (
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
    _ "github.com/mattn/go-sqlite3"
)
// Drivers register themselves in init()
```

**Interface check position:**

```go
// After type definition
type MyWriter struct{}
var _ io.Writer = (*MyWriter)(nil)

// Or at package level
var (
    _ io.Reader = (*MyReader)(nil)
    _ io.Writer = (*MyWriter)(nil)
)
```

**Don't ignore errors in production:**

```go
// BAD
data, _ := ioutil.ReadFile(path)

// GOOD
data, err := ioutil.ReadFile(path)
if err != nil {
    return err
}
```

**Unused variable workaround:**

```go
// During development
x := getValue()
_ = x  // Will use later
// Or
var _ = x
```

## Common Patterns

**Database drivers:**

```go
import _ "github.com/lib/pq"
```

**Test helpers:**

```go
var _ testing.TB = (*MyMock)(nil)
```

**Interface gates:**

```go
var _ json.Marshaler = MyType{}
```

**Range over map (keys only):**

```go
for key := range myMap {  // Value implicitly ignored
```
