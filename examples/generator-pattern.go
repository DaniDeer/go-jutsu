package main

import (
	"context"
	"fmt"
	"time"
)

// Basic generator: count from 0 to max-1
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

// Infinite generator: Fibonacci sequence
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

// Generator with context (cancellable)
func countWithContext(ctx context.Context, max int) <-chan int {
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

// Take first n elements from generator
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

// Filter generator
func filter(input <-chan int, pred func(int) bool) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for val := range input {
			if pred(val) {
				output <- val
			}
		}
	}()
	return output
}

// Map generator
func mapGen(input <-chan int, fn func(int) int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for val := range input {
			output <- fn(val)
		}
	}()
	return output
}

// Prime generator
func primes() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		ch <- 2
		for n := 3; n < 1000; n += 2 {
			if isPrime(n) {
				ch <- n
			}
		}
	}()
	return ch
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Range generator (like Python's range)
func rangeGen(start, end, step int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := start; i < end; i += step {
			ch <- i
		}
	}()
	return ch
}

func main() {
	fmt.Println("=== Generator Pattern ===\n")

	// Example 1: Basic generator
	fmt.Println("1. Basic Generator:")
	for n := range count(5) {
		fmt.Printf("   %d\n", n)
	}
	fmt.Println()

	// Example 2: Infinite generator with take
	fmt.Println("2. Infinite Fibonacci (first 10):")
	for n := range take(10, fibonacci()) {
		fmt.Printf("   %d\n", n)
	}
	fmt.Println()

	// Example 3: Generator with context
	fmt.Println("3. Cancellable Generator:")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	numProcessed := 0
	for n := range countWithContext(ctx, 1000) {
		fmt.Printf("   %d\n", n)
		numProcessed++
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("   (Cancelled after %d items)\n", numProcessed)
	fmt.Println()

	// Example 4: Filtering generator
	fmt.Println("4. Filter Even Numbers:")
	evens := filter(count(10), func(x int) bool { return x%2 == 0 })
	for n := range evens {
		fmt.Printf("   %d\n", n)
	}
	fmt.Println()

	// Example 5: Mapping generator
	fmt.Println("5. Map (Square Numbers):")
	squares := mapGen(count(5), func(x int) int { return x * x })
	for n := range squares {
		fmt.Printf("   %d\n", n)
	}
	fmt.Println()

	// Example 6: Chaining generators
	fmt.Println("6. Chaining: Filter odd, then square, take 5:")
	chain := take(5,
		mapGen(
			filter(count(20), func(x int) bool { return x%2 == 1 }),
			func(x int) int { return x * x },
		),
	)
	for n := range chain {
		fmt.Printf("   %d\n", n)
	}
	fmt.Println()

	// Example 7: Prime numbers
	fmt.Println("7. Prime Numbers (first 10):")
	for n := range take(10, primes()) {
		fmt.Printf("   %d\n", n)
	}
	fmt.Println()

	// Example 8: Range generator
	fmt.Println("8. Range Generator (10 to 30, step 5):")
	for n := range rangeGen(10, 30, 5) {
		fmt.Printf("   %d\n", n)
	}
	fmt.Println()

	fmt.Println("Key Takeaway:")
	fmt.Println("Generator pattern:")
	fmt.Println("- Lazy evaluation (compute on demand)")
	fmt.Println("- Memory efficient")
	fmt.Println("- Infinite sequences possible")
	fmt.Println("- Composable with filter/map/take")
	fmt.Println("- Use context for cancellation")

	fmt.Println("\n=== Done ===")
}
