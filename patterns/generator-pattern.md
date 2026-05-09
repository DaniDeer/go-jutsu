# Generator Pattern

Produce values on demand using channels. Lazy evaluation. Infinite sequences possible.

## What It Is (and Isn't)

Function returns channel that yields values. Values computed as needed. Memory efficient.

Not eager evaluation. Not arrays/slices. Streaming data producer.

## Where You See It

**Basic generator:**

```go
func count(max int) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; i < max; i++ {
            ch <- i
        }
    }()
    return ch
}

for n := range count(5) {
    fmt.Println(n)  // 0, 1, 2, 3, 4
}
```

**Infinite sequence:**

```go
func fibonacci() <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        a, b := 0, 1
        for {
            ch <- a
            a, b = b, a+b
        }
    }()
    return ch
}

// Take first 10
count := 0
for n := range fibonacci() {
    fmt.Println(n)
    count++
    if count >= 10 {
        break
    }
}
```

**Iterator protocol:**

```go
func lines(filename string) <-chan string {
    ch := make(chan string)
    go func() {
        defer close(ch)
        file, _ := os.Open(filename)
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            ch <- scanner.Text()
        }
    }()
    return ch
}
```

## Real Example

```go
// Generate prime numbers
func primes() <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        ch <- 2
        for n := 3; ; n += 2 {
            if isPrime(n) {
                ch <- n
            }
        }
    }()
    return ch
}

// Take with limit
func take(n int, input <-chan int) <-chan int {
    output := make(chan int)
    go func() {
        defer close(output)
        for i := 0; i < n; i++ {
            val, ok := <-input
            if !ok {
                return
            }
            output <- val
        }
    }()
    return output
}

// First 100 primes
for prime := range take(100, primes()) {
    fmt.Println(prime)
}
```

## Gotchas

**Goroutine leak if not consumed:**

```go
gen := count(1000)
// If you don't read from gen, goroutine leaks

// Use context for cancellation
func count(ctx context.Context, max int) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; i < max; i++ {
            select {
            case ch <- i:
            case <-ctx.Done():
                return
            }
        }
    }()
    return ch
}
```

**Break or return needs cleanup:**

```go
for n := range gen() {
    if n > 10 {
        break  // Goroutine may leak
    }
}

// Better: limit generator
for n := range take(10, gen()) {
    process(n)
}
```

**Buffered channels for performance:**

```go
ch := make(chan int, 100)  // Buffer reduces blocking
```

## Benefits

- Lazy evaluation
- Memory efficient (no pre-computation)
- Infinite sequences possible
- Composable with other generators
- Clean iteration protocol
