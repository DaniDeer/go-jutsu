package main

import (
	"fmt"
	"time"
)

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

// bridge flattens channel of channels
func bridge[T any](done <-chan struct{}, chanStream <-chan <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for {
			var stream <-chan T
			select {
			case <-done:
				return
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			}

			// Drain current stream
			for val := range orDone(done, stream) {
				select {
				case out <- val:
				case <-done:
					return
				}
			}
		}
	}()

	return out
}

// Generate batches of channels
func generateBatches(done <-chan struct{}) <-chan <-chan int {
	chanStream := make(chan (<-chan int))

	go func() {
		defer close(chanStream)

		// Batch 1
		batch1 := make(chan int)
		chanStream <- batch1
		go func() {
			defer close(batch1)
			for i := 1; i <= 3; i++ {
				select {
				case batch1 <- i:
				case <-done:
					return
				}
			}
		}()

		time.Sleep(100 * time.Millisecond)

		// Batch 2
		batch2 := make(chan int)
		chanStream <- batch2
		go func() {
			defer close(batch2)
			for i := 10; i <= 12; i++ {
				select {
				case batch2 <- i:
				case <-done:
					return
				}
			}
		}()

		time.Sleep(100 * time.Millisecond)

		// Batch 3
		batch3 := make(chan int)
		chanStream <- batch3
		go func() {
			defer close(batch3)
			for i := 20; i <= 22; i++ {
				select {
				case batch3 <- i:
				case <-done:
					return
				}
			}
		}()
	}()

	return chanStream
}

func main() {
	fmt.Println("=== Bridge Channel Pattern ===")

	done := make(chan struct{})
	defer close(done)

	// Generate stream of channels
	batches := generateBatches(done)

	// Flatten into single channel
	flattened := bridge(done, batches)

	// Consume flattened stream
	fmt.Println("Receiving values from flattened stream:")
	for val := range flattened {
		fmt.Printf("  Received: %d\n", val)
		time.Sleep(50 * time.Millisecond)
	}

	fmt.Println("\nAll batches processed sequentially")
}
