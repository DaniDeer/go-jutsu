package main

import "fmt"

// orDone wraps channel with done signal
func orDone[T any](done <-chan struct{}, c <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
					return
				}
			}
		}
	}()
	return out
}

// Generator for testing
func generate(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func main() {
	fmt.Println("=== Or-Done Channel Pattern ===")

	// Example 1: Normal completion
	fmt.Println("\\nExample 1: Normal completion")
	done := make(chan struct{})
	nums := generate(done, 1, 2, 3, 4, 5)

	for val := range orDone(done, nums) {
		fmt.Printf("Received: %d\\n", val)
	}

	// Example 2: Early cancellation
	fmt.Println("\\nExample 2: Early cancellation")
	done2 := make(chan struct{})
	nums2 := generate(done2, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	count := 0
	for val := range orDone(done2, nums2) {
		fmt.Printf("Received: %d\\n", val)
		count++
		if count == 3 {
			close(done2) // Cancel early
			break
		}
	}
	fmt.Println("Cancelled after 3 values")

	// Example 3: Compare with manual select
	fmt.Println("\\nExample 3: Manual vs OrDone")

	// Manual approach (verbose)
	done3 := make(chan struct{})
	nums3 := generate(done3, 1, 2, 3)
	fmt.Println("Manual select:")
	for {
		select {
		case <-done3:
			fmt.Println("  Done signal received")
			return
		case val, ok := <-nums3:
			if !ok {
				fmt.Println("  Channel closed")
				goto useOrDone
			}
			fmt.Printf("  Got: %d\\n", val)
		}
	}

useOrDone:
	// OrDone approach (clean)
	done4 := make(chan struct{})
	nums4 := generate(done4, 1, 2, 3)
	fmt.Println("Using orDone:")
	for val := range orDone(done4, nums4) {
		fmt.Printf("  Got: %d\\n", val)
	}
}
