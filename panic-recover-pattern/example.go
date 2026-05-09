package main

import (
	"fmt"
	"log"
)

// SafeDivide demonstrates recover in function
func SafeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()

	result = a / b // Panics if b == 0
	return result, nil
}

// ProcessItems demonstrates recovering from panics in loops
func ProcessItems(items []int) (processed []int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("processing failed: %v", r)
		}
	}()

	for _, item := range items {
		// Simulate processing that might panic
		if item < 0 {
			panic("negative values not allowed")
		}
		processed = append(processed, item*2)
	}

	return processed, nil
}

// BadRecover demonstrates WRONG usage
func BadRecover() {
	// This does NOT work - recover must be in deferred function
	recover()
	panic("this will not be caught")
}

// Server simulates HTTP handler with panic protection
type Server struct{}

func (s *Server) HandleRequest(id int) (response string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Handler panic: %v", r)
			response = "Internal Server Error"
		}
	}()

	if id < 0 {
		panic("invalid request ID")
	}

	return fmt.Sprintf("Request %d processed", id)
}

// MustCompile demonstrates acceptable panic usage
func MustCompile(pattern string) {
	// Panic during initialization is acceptable
	if pattern == "" {
		panic("empty pattern not allowed")
	}
	fmt.Printf("Pattern compiled: %s\n", pattern)
}

func main() {
	// Example 1: Safe division
	fmt.Println("=== Example 1: Safe Division ===")
	result, err := SafeDivide(10, 2)
	fmt.Printf("10 / 2 = %d, error: %v\n", result, err)

	result, err = SafeDivide(10, 0)
	fmt.Printf("10 / 0 = %d, error: %v\n", result, err)

	// Example 2: Processing with recovery
	fmt.Println("\n=== Example 2: Process Items ===")
	items1 := []int{1, 2, 3, 4}
	processed, err := ProcessItems(items1)
	fmt.Printf("Processed %v: %v, error: %v\n", items1, processed, err)

	items2 := []int{1, 2, -3, 4}
	processed, err = ProcessItems(items2)
	fmt.Printf("Processed %v: %v, error: %v\n", items2, processed, err)

	// Example 3: Server panic handling
	fmt.Println("\n=== Example 3: Server Handler ===")
	server := &Server{}
	fmt.Println(server.HandleRequest(123))
	fmt.Println(server.HandleRequest(-1)) // Panics internally

	// Example 4: Init-time panic (acceptable)
	fmt.Println("\n=== Example 4: Must Compile ===")
	MustCompile("valid-pattern")

	// Uncomment to see panic:
	// MustCompile("")

	// Example 5: Demonstrate nested panics
	fmt.Println("\n=== Example 5: Nested Recovery ===")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Caught panic: %v\n", r)
			}
		}()

		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Inner recovery: %v\n", r)
					panic("re-panicking") // Can re-panic if needed
				}
			}()
			panic("original panic")
		}()
	}()

	fmt.Println("\n✓ All examples completed without crashing")
}
