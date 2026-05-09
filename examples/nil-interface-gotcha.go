package main

import (
	"fmt"
)

// Custom error type
type MyError struct {
	Msg string
}

func (e *MyError) Error() string {
	return e.Msg
}

// Example 1: The classic gotcha
func returnsNonNilError() error {
	var err *MyError = nil
	// Do some work...
	return err // Returns non-nil error!
}

func demonstrateGotcha() {
	fmt.Println("1. The Classic Gotcha:")
	err := returnsNonNilError()
	fmt.Printf("   err == nil: %v\n", err == nil)
	fmt.Printf("   err value: %v\n", err)
	if err != nil {
		fmt.Println("   ⚠ This branch executes even though we returned nil!")
	}
	fmt.Println()
}

// Example 2: Correct way
func returnsCorrectNil() error {
	var err *MyError = nil
	// Do some work...
	if err != nil {
		return err
	}
	return nil // Explicit nil
}

func demonstrateCorrect() {
	fmt.Println("2. Correct Approach:")
	err := returnsCorrectNil()
	fmt.Printf("   err == nil: %v\n", err == nil)
	fmt.Printf("   err value: %v\n", err)
	if err != nil {
		fmt.Println("   This won't execute")
	} else {
		fmt.Println("   ✓ Correctly recognized as nil")
	}
	fmt.Println()
}

// Example 3: Interface internals demonstration
func showInterfaceInternals() {
	fmt.Println("3. Understanding Interface Internals:")

	// Nil interface
	var err1 error = nil
	fmt.Printf("   var err error = nil\n")
	fmt.Printf("   err == nil: %v\n", err1 == nil)
	fmt.Printf("   Type: %T, Value: %v\n", err1, err1)

	// Non-nil interface with nil pointer
	var myErr *MyError = nil
	var err2 error = myErr
	fmt.Printf("\n   var myErr *MyError = nil\n")
	fmt.Printf("   var err error = myErr\n")
	fmt.Printf("   err == nil: %v\n", err2 == nil)
	fmt.Printf("   Type: %T, Value: %v\n", err2, err2)
	fmt.Println()
}

// Example 4: Real-world scenario
type Database struct {
	connected bool
}

func (db *Database) Query() error {
	if !db.connected {
		return &MyError{"not connected"}
	}
	return nil
}

func connectDB() (*Database, error) {
	// Simulating connection failure
	var db *Database = nil
	return db, nil // Oops, returning nil pointer
}

func demonstrateRealWorld() {
	fmt.Println("4. Real-World Scenario (Database):")

	db, err := connectDB()
	if err != nil {
		fmt.Println("   Connection error:", err)
		return
	}

	fmt.Printf("   db == nil: %v\n", db == nil)

	// This will panic! db is nil
	// err = db.Query()
	fmt.Println("   ⚠ Can't call db.Query() - db is nil pointer")
	fmt.Println()
}

// Example 5: Function that returns interface
func getWriter() interface{} {
	var w *Database = nil
	return w // Returns non-nil interface
}

func demonstrateInterfaceReturn() {
	fmt.Println("5. Function Returning Interface:")

	writer := getWriter()
	fmt.Printf("   writer == nil: %v\n", writer == nil)
	fmt.Printf("   Type: %T, Value: %v\n", writer, writer)
	fmt.Println()
}

// Example 6: Comparison between nil interfaces
func demonstrateComparison() {
	fmt.Println("6. Interface Comparison:")

	var a interface{} = (*int)(nil)
	var b interface{} = (*int)(nil)
	var c interface{} = nil

	fmt.Printf("   a := (*int)(nil) as interface{}\n")
	fmt.Printf("   b := (*int)(nil) as interface{}\n")
	fmt.Printf("   c := nil as interface{}\n\n")

	fmt.Printf("   a == b: %v (same type and value)\n", a == b)
	fmt.Printf("   a == nil: %v (type is set)\n", a == nil)
	fmt.Printf("   c == nil: %v (both nil)\n", c == nil)
	fmt.Println()
}

// Example 7: How to fix - pattern 1
func safeFuncPattern1() error {
	var err *MyError = nil

	// Some logic that might set err
	shouldFail := false
	if shouldFail {
		err = &MyError{"something went wrong"}
	}

	// Check before returning
	if err != nil {
		return err
	}
	return nil
}

// Example 8: How to fix - pattern 2
func safeFuncPattern2() error {
	// Use error type directly
	var err error = nil

	shouldFail := false
	if shouldFail {
		err = &MyError{"something went wrong"}
	}

	return err // Safe because err is already error type
}

// Example 9: How to fix - pattern 3
func safeFuncPattern3() error {
	shouldFail := false
	if shouldFail {
		return &MyError{"something went wrong"}
	}
	return nil // Most explicit and clear
}

func demonstrateFixes() {
	fmt.Println("7. Three Safe Patterns:")

	err1 := safeFuncPattern1()
	fmt.Printf("   Pattern 1 (check before return): err == nil: %v\n", err1 == nil)

	err2 := safeFuncPattern2()
	fmt.Printf("   Pattern 2 (use error type): err == nil: %v\n", err2 == nil)

	err3 := safeFuncPattern3()
	fmt.Printf("   Pattern 3 (explicit return): err == nil: %v\n", err3 == nil)
	fmt.Println()
}

func main() {
	fmt.Println("=== Nil Interface Gotcha Examples ===\n")

	demonstrateGotcha()
	demonstrateCorrect()
	showInterfaceInternals()
	demonstrateRealWorld()
	demonstrateInterfaceReturn()
	demonstrateComparison()
	demonstrateFixes()

	fmt.Println("Key Takeaway:")
	fmt.Println("Interface = (type, value)")
	fmt.Println("nil only if BOTH are nil")
	fmt.Println("Always return explicit nil for interfaces!")

	fmt.Println("\n=== Done ===")
}
