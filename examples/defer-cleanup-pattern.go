package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Example 1: Basic resource cleanup
func readFileWithDefer(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close() // Guaranteed cleanup

	// Work with file...
	info, err := file.Stat()
	if err != nil {
		return err // defer still runs
	}

	fmt.Printf("File size: %d bytes\n", info.Size())
	return nil
}

// Example 2: Defer for timing
func timeTrack(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func slowOperation() {
	defer timeTrack("slowOperation")() // Note the extra ()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Operation complete")
}

// Example 3: Defer with named return - modify result
func mayFail() (result string, err error) {
	defer func() {
		if err != nil {
			result = "failed" // Modify return on error
			log.Printf("Error: %v", err)
		}
	}()

	// Simulate error
	return "", fmt.Errorf("oops")
	// Actually returns: "failed", error
}

// Example 4: Multiple defers - LIFO
func demonstrateLIFO() {
	fmt.Println("\nDefer execution order (LIFO):")
	defer fmt.Println("  Third defer - executes FIRST")
	defer fmt.Println("  Second defer - executes SECOND")
	defer fmt.Println("  First defer - executes THIRD")
	fmt.Println("Function body")
}

// Example 5: Defer gotcha - loop
func deferInLoopBad() {
	fmt.Println("\n❌ Bad: Defer in loop (files stay open)")
	files := []string{"file1", "file2", "file3"}

	for _, name := range files {
		fmt.Printf("  Opening %s (but not closing until function ends)\n", name)
		// defer would pile up here - don't do this!
		_ = name
	}
}

func deferInLoopGood() {
	fmt.Println("\n✓ Good: Defer in closure")
	files := []string{"file1", "file2", "file3"}

	for _, name := range files {
		func() {
			fmt.Printf("  Opening and closing %s immediately\n", name)
			// defer close() here - runs at end of this func
		}()
	}
}

// Example 6: Defer with immediate vs captured values
func deferArgumentCapture() {
	fmt.Println("\nDefer argument evaluation:")

	i := 0
	defer fmt.Printf("  Immediate arg: %d\n", i)        // Evaluates to 0 NOW
	defer func() { fmt.Printf("  Closure: %d\n", i) }() // Captures variable

	i = 42
	fmt.Println("  Set i = 42")
}

func main() {
	fmt.Println("=== Defer Pattern Examples ===\n")

	// Example 1: File cleanup
	fmt.Println("1. Resource cleanup with defer:")
	err := readFileWithDefer("defer-cleanup-pattern.go")
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	// Example 2: Timing
	fmt.Println("\n2. Function timing:")
	slowOperation()

	// Example 3: Named return
	fmt.Println("\n3. Modifying return value:")
	result, err := mayFail()
	fmt.Printf("  Result: %q, Error: %v\n", result, err)

	// Example 4: LIFO
	demonstrateLIFO()

	// Example 5: Loop gotcha
	deferInLoopBad()
	deferInLoopGood()

	// Example 6: Argument capture
	deferArgumentCapture()

	fmt.Println("\n=== Done ===")
}
