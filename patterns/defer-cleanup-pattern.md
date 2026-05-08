# Defer Cleanup Pattern

Stack cleanup operations that execute in reverse order on function exit, guaranteed even on panic.

## What It Is (and Isn't)

`defer` schedules function call to run **after** surrounding function returns. Stack-based LIFO execution.

Not like try-finally (no catch). Not async (runs on same goroutine). Unique to Go.

## Where You See It

**Resource cleanup:**

```go
file, _ := os.Open("file.txt")
defer file.Close()  // Guaranteed cleanup
```

**Unlock mutexes:**

```go
mu.Lock()
defer mu.Unlock()  // Can't forget
```

**Timing functions:**

```go
func trace() func() {
    start := time.Now()
    return func() { log.Printf("took %v", time.Since(start)) }
}

func operation() {
    defer trace()()  // Measures operation time
    // ... work ...
}
```

**Modify return value:**

```go
func increment() (result int) {
    defer func() { result++ }()  // Runs after return, modifies named result
    return 5  // Actually returns 6
}
```

## Real Example

```go
func processFile(path string) (err error) {
    // Multiple defers stack up - execute in reverse
    defer log.Printf("processFile done, err=%v", err)

    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // Runs before log

    data, err := io.ReadAll(f)
    if err != nil {
        return err  // defer still runs
    }

    return process(data)
}
```

## Gotchas

**Defer in loop** - creates stack of defers (memory leak):

```go
// BAD
for _, file := range files {
    f, _ := os.Open(file)
    defer f.Close()  // Defers pile up, don't run until function ends
}

// GOOD
for _, file := range files {
    func() {
        f, _ := os.Open(file)
        defer f.Close()  // Runs at end of this func
    }()
}
```

**Defer evaluates args immediately:**

```go
func bad() {
    i := 0
    defer log.Println(i)  // Captures 0 now
    i++
    // Prints: 0 (not 1)
}

func good() {
    i := 0
    defer func() { log.Println(i) }()  // Closure captures variable
    i++
    // Prints: 1
}
```

**Execution order** - LIFO (Last In First Out):

```go
defer log.Println("1")
defer log.Println("2")
defer log.Println("3")
// Prints: 3, 2, 1
```
