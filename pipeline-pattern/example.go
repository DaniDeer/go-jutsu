package main

import (
	"fmt"
	"sync"
	"time"
)

// Stage 1: Generator
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Stage 2: Square
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// Stage 3: Add offset
func add(in <-chan int, offset int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n + offset
		}
	}()
	return out
}

// Fan-out: distribute work to multiple workers
func fanOut(in <-chan int, workers int) []<-chan int {
	channels := make([]<-chan int, workers)
	for i := 0; i < workers; i++ {
		channels[i] = worker(in, i+1)
	}
	return channels
}

func worker(in <-chan int, id int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			time.Sleep(10 * time.Millisecond) // Simulate work
			fmt.Printf("   Worker %d processed: %d\n", id, n)
			out <- n * 2
		}
	}()
	return out
}

// Fan-in: merge multiple channels into one
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				out <- n
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	fmt.Println("=== Pipeline Pattern ===\n")

	// Example 1: Simple pipeline
	fmt.Println("1. Simple Pipeline:")
	nums := generator(1, 2, 3, 4)
	squared := square(nums)
	result := add(squared, 10)

	for n := range result {
		fmt.Printf("   Result: %d\n", n)
	}
	fmt.Println()

	// Example 2: Inline pipeline
	fmt.Println("2. Inline Pipeline:")
	for n := range add(square(generator(5, 6, 7)), 100) {
		fmt.Printf("   Result: %d\n", n)
	}
	fmt.Println()

	// Example 3: Fan-out/Fan-in
	fmt.Println("3. Fan-Out/Fan-In (Parallel Processing):")
	input := generator(1, 2, 3, 4, 5, 6)
	workers := fanOut(input, 3) // 3 parallel workers
	merged := fanIn(workers...)

	var results []int
	for n := range merged {
		results = append(results, n)
	}
	fmt.Printf("   Final results: %v\n", results)
	fmt.Println()

	// Example 4: Pipeline with cancellation
	fmt.Println("4. Pipeline with Done Signal:")
	done := make(chan struct{})
	nums4 := generatorWithDone(done, 1, 2, 3, 4, 5)

	go func() {
		time.Sleep(50 * time.Millisecond)
		close(done) // Cancel after 50ms
	}()

	count := 0
	for n := range nums4 {
		fmt.Printf("   Got: %d\n", n)
		count++
		time.Sleep(20 * time.Millisecond)
	}
	fmt.Printf("   Processed %d items before cancellation\n", count)
	fmt.Println()

	fmt.Println("Key Takeaway:")
	fmt.Println("Pipeline pattern enables:")
	fmt.Println("- Concurrent stream processing")
	fmt.Println("- Composable stages")
	fmt.Println("- Fan-out for parallelism")
	fmt.Println("- Fan-in for aggregation")
	fmt.Println("- Memory efficient (no buffering all data)")

	fmt.Println("\n=== Done ===")
}

// Generator with done channel for cancellation
func generatorWithDone(done <-chan struct{}, nums ...int) <-chan int {
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
