# Comma-Ok Idiom

Two-value return for safe access. Prevents panics. Everywhere in Go.

## What It Is (and Isn't)

Pattern: `value, ok := operation`. Second return is boolean indicating success.

Not exceptions. Not null checks. Explicit success indication built into the language.

## Where You See It

**Map access:**

```go
value, ok := myMap[key]
if !ok {
    // Key doesn't exist
}
```

**Channel receive:**

```go
value, ok := <-ch
if !ok {
    // Channel closed
}
```

**Type assertion:**

```go
str, ok := value.(string)
if !ok {
    // Not a string
}
```

## Real Example

```go
func processUser(users map[int]User, id int) error {
    // Without comma-ok: panic if key missing
    // user := users[id]  // Dangerous!

    // With comma-ok: safe
    user, ok := users[id]
    if !ok {
        return fmt.Errorf("user %d not found", id)
    }

    return user.Process()
}

func readFromChannel(ch <-chan string) {
    for {
        msg, ok := <-ch
        if !ok {
            log.Println("Channel closed")
            return
        }
        process(msg)
    }
}

func handleValue(val interface{}) {
    // Try as string first
    if str, ok := val.(string); ok {
        fmt.Println("String:", str)
        return
    }

    // Try as int
    if num, ok := val.(int); ok {
        fmt.Println("Number:", num)
        return
    }

    fmt.Println("Unknown type")
}
```

## Gotchas

**Map zero value vs missing:**

```go
m := map[string]int{"a": 0}

// Both return 0, but different meanings
fmt.Println(m["a"])  // 0 (exists)
fmt.Println(m["b"])  // 0 (doesn't exist)

// Use comma-ok to distinguish
if val, ok := m["a"]; ok {
    fmt.Println("Found:", val)  // "Found: 0"
}

if val, ok := m["b"]; ok {
    // Doesn't execute
} else {
    fmt.Println("Not found")
}
```

**Type assertion panic without ok:**

```go
var val interface{} = 42

// Panics if wrong type
str := val.(string)  // panic: interface conversion

// Safe version
str, ok := val.(string)
if !ok {
    // Handle gracefully
}
```

**Channel closed vs no data:**

```go
ch := make(chan int, 1)
ch <- 42
close(ch)

// First receive: gets value
val, ok := <-ch
fmt.Println(val, ok)  // 42, true

// Second receive: channel closed
val, ok = <-ch
fmt.Println(val, ok)  // 0, false

// Without ok, you get zero value
val = <-ch  // 0 (looks like valid data!)
```

**Ignoring ok is common mistake:**

```go
// BAD: assumes key exists
user := users[id]  // Zero value if missing

// GOOD: check existence
if user, ok := users[id]; ok {
    // Use user safely
}
```

**Range already uses comma-ok:**

```go
// These are equivalent
for _, val := range ch {
    // ...
}

// Explicit version
for {
    val, ok := <-ch
    if !ok {
        break
    }
    // ...
}
```

## When to Use

- Always for type assertions (prevents panic)
- Maps when zero value is valid data
- Channels when you need to detect closure
- Anywhere panic would ruin your day
