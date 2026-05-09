package main

import (
	"fmt"
	"sync"
	"time"
)

// Example 1: Channel signaling
func channelSignaling() {
	fmt.Println("1. Channel signaling with struct{}:")

	done := make(chan struct{})

	go func() {
		fmt.Println("   Worker: doing work...")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("   Worker: done!")
		done <- struct{}{} // Send signal (0 bytes)
	}()

	<-done // Wait for signal
	fmt.Println("   Main: received completion signal\n")
}

// Example 2: Set implementation
func setImplementation() {
	fmt.Println("2. Set using map[string]struct{}:")

	// Set of unique strings
	seen := make(map[string]struct{})

	words := []string{"apple", "banana", "apple", "cherry", "banana"}

	for _, word := range words {
		seen[word] = struct{}{} // Add to set
	}

	fmt.Println("   Input:", words)
	fmt.Print("   Unique: ")
	for word := range seen {
		fmt.Print(word, " ")
	}
	fmt.Println("\n")
}

// Example 3: Semaphore pattern
func semaphorePattern() {
	fmt.Println("3. Semaphore for limiting concurrency:")

	const maxWorkers = 3
	sem := make(chan struct{}, maxWorkers)

	var wg sync.WaitGroup
	tasks := 10

	for i := 1; i <= tasks; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire slot
			defer func() { <-sem }() // Release slot

			fmt.Printf("   Task %d running\n", id)
			time.Sleep(50 * time.Millisecond)
		}(i)
	}

	wg.Wait()
	fmt.Println("   All tasks complete\n")
}

// Example 4: Zero-size method receiver
type StatelessLogger struct{}

func (StatelessLogger) Info(msg string) {
	fmt.Printf("   [INFO] %s\n", msg)
}

func (StatelessLogger) Error(msg string) {
	fmt.Printf("   [ERROR] %s\n", msg)
}

func zeroSizeReceiver() {
	fmt.Println("4. Zero-size type as method receiver:")

	logger := StatelessLogger{} // Takes 0 bytes!
	logger.Info("Application started")
	logger.Error("Something went wrong")

	// Prove it's zero-size
	fmt.Printf("   Size of StatelessLogger: %d bytes\n\n",
		0) // Would use unsafe.Sizeof in real code
}

// Example 5: Coordinating multiple goroutines
func coordinateGoroutines() {
	fmt.Println("5. Coordinating multiple goroutines:")

	done := make(chan struct{})
	results := make(chan string, 3)

	workers := []string{"Alice", "Bob", "Charlie"}

	for _, name := range workers {
		go func(n string) {
			time.Sleep(time.Duration(50) * time.Millisecond)
			results <- fmt.Sprintf("%s finished", n)
			done <- struct{}{} // Signal completion
		}(name)
	}

	// Wait for all workers
	for range workers {
		<-done
	}
	close(results)

	fmt.Println("   Results:")
	for result := range results {
		fmt.Printf("   - %s\n", result)
	}
	fmt.Println()
}

// Example 6: Context-like cancellation
func cancellationPattern() {
	fmt.Println("6. Cancellation signal pattern:")

	stop := make(chan struct{})

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("   Tick...")
			case <-stop:
				fmt.Println("   Stopped!")
				return
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	close(stop) // Broadcast stop to all listeners
	time.Sleep(50 * time.Millisecond)
	fmt.Println()
}

// Example 7: Memory comparison
func memoryComparison() {
	fmt.Println("7. Memory comparison (conceptual):")

	// Both channels work the same, but struct{} uses 0 bytes per value
	boolChan := make(chan bool, 1)
	structChan := make(chan struct{}, 1)

	boolChan <- true
	structChan <- struct{}{}

	fmt.Println("   chan bool:     sends/receives 1 byte per signal")
	fmt.Println("   chan struct{}: sends/receives 0 bytes per signal")
	fmt.Println("   For millions of signals, this matters!\n")
}

func main() {
	fmt.Println("=== Empty Struct Pattern Examples ===\n")

	channelSignaling()
	setImplementation()
	semaphorePattern()
	zeroSizeReceiver()
	coordinateGoroutines()
	cancellationPattern()
	memoryComparison()

	fmt.Println("=== Done ===")
}
