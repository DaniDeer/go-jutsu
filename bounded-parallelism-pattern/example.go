package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ProcessFiles demonstrates bounded parallelism with semaphore
func ProcessFiles(files []string, maxConcurrent int) {
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		sem <- struct{}{} // Acquire token

		go func(f string) {
			defer wg.Done()
			defer func() { <-sem }() // Release token

			// Simulate processing
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Processed: %s\n", f)
		}(file)
	}

	wg.Wait()
}

// ProcessWithContext demonstrates context-aware bounded parallelism
func ProcessWithContext(ctx context.Context, items []int, maxConcurrent int) []int {
	sem := make(chan struct{}, maxConcurrent)
	results := make(chan int, len(items))
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Done()
			continue
		}

		go func(val int) {
			defer wg.Done()
			defer func() { <-sem }()

			select {
			case <-ctx.Done():
				return
			case <-time.After(200 * time.Millisecond):
				results <- val * 2
			}
		}(item)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []int
	for r := range results {
		out = append(out, r)
	}
	return out
}

func main() {
	// Example 1: Basic bounded parallelism
	fmt.Println("=== Example 1: Bounded Parallelism (max 3) ===")
	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt"}
	start := time.Now()
	ProcessFiles(files, 3)
	fmt.Printf("Completed in %v\n", time.Since(start))

	// Example 2: With context cancellation
	fmt.Println("\n=== Example 2: With Context ===")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	results := ProcessWithContext(ctx, items, 2)
	fmt.Printf("Processed %d items: %v\n", len(results), results)

	// Example 3: Compare unlimited vs bounded
	fmt.Println("\n=== Example 3: Comparison ===")
	n := 10

	// Unlimited goroutines
	start = time.Now()
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
		}()
	}
	wg.Wait()
	fmt.Printf("Unlimited (%d goroutines): %v\n", n, time.Since(start))

	// Bounded to 3
	start = time.Now()
	ProcessFiles(make([]string, n), 3)
	fmt.Printf("Bounded (3 workers): %v\n", time.Since(start))
}
