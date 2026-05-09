package main

import (
	"fmt"
	"time"
)

// Example 1: Basic select with timeout
func demonstrateTimeout() {
	fmt.Println("1. Select with Timeout:")

	ch := make(chan string)

	// Slow sender
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "data arrived"
	}()

	select {
	case msg := <-ch:
		fmt.Printf("   ✓ Received: %s\n", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("   ✗ Timeout after 1 second")
	}
	fmt.Println()
}

// Example 2: Non-blocking receive
func demonstrateNonBlocking() {
	fmt.Println("2. Non-Blocking Channel Operation:")

	ch := make(chan int, 1)
	ch <- 42

	// Try to receive
	select {
	case msg := <-ch:
		fmt.Printf("   ✓ Received: %d\n", msg)
	default:
		fmt.Println("   Channel empty")
	}

	// Try again (now empty)
	select {
	case msg := <-ch:
		fmt.Printf("   Received: %d\n", msg)
	default:
		fmt.Println("   ✓ Channel empty (didn't block)")
	}
	fmt.Println()
}

// Example 3: Multiple channels
func demonstrateMultipleChannels() {
	fmt.Println("3. Multiple Channels:")

	ch1 := make(chan string)
	ch2 := make(chan string)
	done := make(chan struct{})

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from ch1"
	}()

	go func() {
		time.Sleep(150 * time.Millisecond)
		ch2 <- "from ch2"
	}()

	go func() {
		time.Sleep(300 * time.Millisecond)
		close(done)
	}()

	for i := 0; i < 3; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("   Received: %s\n", msg)
		case msg := <-ch2:
			fmt.Printf("   Received: %s\n", msg)
		case <-done:
			fmt.Println("   Done signal received")
			return
		}
	}
	fmt.Println()
}

// Example 4: Fan-in pattern
func fanIn(ch1, ch2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < 4; i++ {
			select {
			case v := <-ch1:
				out <- v
			case v := <-ch2:
				out <- v
			}
		}
	}()
	return out
}

func demonstrateFanIn() {
	fmt.Println("4. Fan-In Pattern:")

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 0; i < 2; i++ {
			ch1 <- i * 10
			time.Sleep(50 * time.Millisecond)
		}
	}()

	go func() {
		for i := 0; i < 2; i++ {
			ch2 <- i * 100
			time.Sleep(75 * time.Millisecond)
		}
	}()

	merged := fanIn(ch1, ch2)
	for val := range merged {
		fmt.Printf("   Received: %d\n", val)
	}
	fmt.Println()
}

// Example 5: Worker with cancellation
func worker(id int, jobs <-chan int, results chan<- int, done <-chan struct{}) {
	for {
		select {
		case job := <-jobs:
			results <- job * 2
		case <-done:
			fmt.Printf("   Worker %d stopping\n", id)
			return
		}
	}
}

func demonstrateWorkerCancellation() {
	fmt.Println("5. Worker with Cancellation:")

	jobs := make(chan int, 5)
	results := make(chan int, 5)
	done := make(chan struct{})

	go worker(1, jobs, results, done)

	// Send jobs
	for i := 1; i <= 3; i++ {
		jobs <- i
	}

	// Get results
	for i := 1; i <= 3; i++ {
		fmt.Printf("   Result: %d\n", <-results)
	}

	// Cancel worker
	close(done)
	time.Sleep(50 * time.Millisecond)
	fmt.Println()
}

// Example 6: Random selection
func demonstrateRandom() {
	fmt.Println("6. Random Selection (when multiple ready):")

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1 <- 1
	ch2 <- 2

	// Both ready - random selection
	count1, count2 := 0, 0
	for i := 0; i < 10; i++ {
		ch1 <- i
		ch2 <- i

		select {
		case <-ch1:
			count1++
		case <-ch2:
			count2++
		}
	}

	fmt.Printf("   ch1 selected: %d times\n", count1)
	fmt.Printf("   ch2 selected: %d times\n", count2)
	fmt.Println("   ✓ Non-deterministic selection!\n")
}

// Example 7: Ticker with select
func demonstrateTicker() {
	fmt.Println("7. Ticker with Select:")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	done := time.After(350 * time.Millisecond)

	for {
		select {
		case t := <-ticker.C:
			fmt.Printf("   Tick at %s\n", t.Format("15:04:05.000"))
		case <-done:
			fmt.Println("   ✓ Done")
			fmt.Println()
			return
		}
	}
}

// Example 8: Priority select
func demonstratePriority() {
	fmt.Println("8. Priority Select:")

	high := make(chan string, 2)
	normal := make(chan string, 2)

	high <- "HIGH1"
	normal <- "normal1"
	high <- "HIGH2"
	normal <- "normal2"

	for i := 0; i < 4; i++ {
		// Check high priority first
		select {
		case msg := <-high:
			fmt.Printf("   %s (priority)\n", msg)
		default:
			// Then normal
			select {
			case msg := <-normal:
				fmt.Printf("   %s\n", msg)
			default:
				fmt.Println("   All empty")
			}
		}
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== Select Statement Pattern ===\n")

	demonstrateTimeout()
	demonstrateNonBlocking()
	demonstrateMultipleChannels()
	demonstrateFanIn()
	demonstrateWorkerCancellation()
	demonstrateRandom()
	demonstrateTicker()
	demonstratePriority()

	fmt.Println("Key Takeaway:")
	fmt.Println("select waits on multiple channels")
	fmt.Println("- First ready case wins")
	fmt.Println("- Non-deterministic if multiple ready")
	fmt.Println("- default makes it non-blocking")
	fmt.Println("- Common patterns: timeout, cancellation, fan-in")

	fmt.Println("\n=== Done ===")
}
