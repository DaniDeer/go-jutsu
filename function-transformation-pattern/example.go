package main

// Function transformation: take a function, return an enhanced version.
// The original function's logic is untouched — behavior is added around it.

import (
	"errors"
	"fmt"
	"time"
)

// withTiming wraps fn, printing how long it takes to run.
func withTiming(name string, fn func()) func() {
	return func() {
		start := time.Now()
		fn()
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

// withRetry wraps fn, retrying up to attempts times on error.
func withRetry(attempts int, fn func() error) func() error {
	return func() error {
		var err error
		for i := range attempts {
			if err = fn(); err == nil {
				return nil
			}
			fmt.Printf("  attempt %d/%d failed: %v\n", i+1, attempts, err)
		}
		return fmt.Errorf("all %d attempts failed: %w", attempts, err)
	}
}

// memoize caches results of a pure function. Generic: works for any key/value types.
func memoize[K comparable, V any](fn func(K) V) func(K) V {
	cache := make(map[K]V)
	return func(k K) V {
		if v, ok := cache[k]; ok {
			fmt.Printf("  cache hit for key=%v\n", k)
			return v
		}
		v := fn(k)
		cache[k] = v
		return v
	}
}

func main() {
	// --- withTiming ---
	slowOp := withTiming("slow-op", func() {
		time.Sleep(5 * time.Millisecond)
	})
	slowOp() // slow-op took 5.Xms

	fmt.Println()

	// --- withRetry ---
	call := 0
	flakyFn := withRetry(3, func() error {
		call++
		if call < 3 {
			return errors.New("service unavailable")
		}
		return nil
	})
	if err := flakyFn(); err != nil {
		fmt.Println("failed:", err)
	} else {
		fmt.Printf("succeeded on attempt %d\n", call)
	}

	fmt.Println()

	// --- memoize ---
	callCount := 0
	expensiveFn := memoize(func(n int) string {
		callCount++
		time.Sleep(1 * time.Millisecond) // simulate work
		return fmt.Sprintf("result-%d", n)
	})

	fmt.Println(expensiveFn(10))                   // computed
	fmt.Println(expensiveFn(10))                   // cache hit
	fmt.Println(expensiveFn(20))                   // computed
	fmt.Printf("actual fn calls: %d\n", callCount) // 2
}
