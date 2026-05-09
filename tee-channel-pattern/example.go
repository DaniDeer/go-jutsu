package main

import (
	"fmt"
	"time"
)

// tee duplicates channel to two outputs
func tee[T any](done <-chan struct{}, in <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in) {
			var o1, o2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case o1 <- val:
					o1 = nil // Disable this case
				case o2 <- val:
					o2 = nil
				}
			}
		}
	}()

	return out1, out2
}

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
	fmt.Println("=== Tee Channel Pattern ===")

	done := make(chan struct{})
	defer close(done)

	// Generate numbers
	source := generate(done, 1, 2, 3, 4, 5)

	// Split into two channels
	ch1, ch2 := tee(done, source)

	// Consumer 1: Print values
	go func() {
		for val := range ch1 {
			fmt.Printf("Consumer 1: %d\n", val)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// Consumer 2: Print squares
	go func() {
		for val := range ch2 {
			fmt.Printf("Consumer 2: %d² = %d\n", val, val*val)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// Wait for completion
	time.Sleep(500 * time.Millisecond)
	fmt.Println("\nBoth consumers received all values")
}
