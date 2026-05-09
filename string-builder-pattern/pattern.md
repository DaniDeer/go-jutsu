# String Builder Pattern

Efficient string concatenation. `strings.Builder` avoids allocations. Use instead of `+` in loops.

## What It Is (and Isn't)

Mutable byte buffer for strings. Grows dynamically. Much faster than repeated concatenation.

Not needed for few strings. Not thread-safe. Pattern for building strings incrementally.

## Where You See It

**Loop concatenation:**
```go
var b strings.Builder
for _, s := range items {
    b.WriteString(s)
}
result := b.String()
```

**Template-like building:**
```go
b.WriteString("<html>")
b.WriteString("<body>")
b.WriteString(content)
b.WriteString("</body></html>")
```

## Real Example

```go
import "strings"

func BuildSQL(table string, cols []string, rows [][]interface{}) string {
    var b strings.Builder
    b.Grow(256)  // Preallocate if you know size
    
    b.WriteString("INSERT INTO ")
    b.WriteString(table)
    b.WriteString(" (")
    for i, col := range cols {
        if i > 0 {
            b.WriteString(", ")
        }
        b.WriteString(col)
    }
    b.WriteString(") VALUES ")
    
    for i, row := range rows {
        if i > 0 {
            b.WriteString(", ")
        }
        b.WriteString("(")
        for j, val := range row {
            if j > 0 {
                b.WriteString(", ")
            }
            fmt.Fprintf(&b, "'%v'", val)
        }
        b.WriteString(")")
    }
    
    return b.String()
}
```

## Gotchas

**Don't use + in loops:**
```go
// BAD: O(n²) allocations
var s string
for _, item := range items {
    s += item
}

// GOOD: O(n)
var b strings.Builder
for _, item := range items {
    b.WriteString(item)
}
```

**Grow() optimization:**
```go
var b strings.Builder
b.Grow(expectedSize)  // Preallocate capacity
```

**fmt.Fprintf works:**
```go
fmt.Fprintf(&b, "value: %d", num)  // Builder implements io.Writer
```

**Can't copy Builder:**
```go
b1 := strings.Builder{}
b2 := b1  // Panics on use! Pass by pointer
```

