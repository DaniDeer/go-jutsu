package main

import (
	"fmt"
	"time"
)

// Example 1: Map access
func demonstrateMapAccess() {
	fmt.Println("1. Map Access with Comma-Ok:")

	scores := map[string]int{
		"Alice": 95,
		"Bob":   0, // Actual score of 0
	}

	// Without comma-ok - ambiguous
	fmt.Println("   Without comma-ok:")
	fmt.Printf("   Bob's score: %d\n", scores["Bob"])         // 0 (exists)
	fmt.Printf("   Charlie's score: %d\n", scores["Charlie"]) // 0 (doesn't exist)

	// With comma-ok - clear
	fmt.Println("\n   With comma-ok:")
	if score, ok := scores["Bob"]; ok {
		fmt.Printf("   ✓ Bob exists with score: %d\n", score)
	}

	if score, ok := scores["Charlie"]; ok {
		fmt.Printf("   Charlie exists with score: %d\n", score)
	} else {
		fmt.Println("   ✗ Charlie not found")
	}
	fmt.Println()
}

// Example 2: Type assertion
func processValue(val interface{}) {
	// Unsafe way (panics)
	// str := val.(string)  // Would panic if val isn't string

	// Safe way with comma-ok
	if str, ok := val.(string); ok {
		fmt.Printf("   String: %q\n", str)
		return
	}

	if num, ok := val.(int); ok {
		fmt.Printf("   Number: %d\n", num)
		return
	}

	if flag, ok := val.(bool); ok {
		fmt.Printf("   Boolean: %v\n", flag)
		return
	}

	fmt.Printf("   Unknown type: %T\n", val)
}

func demonstrateTypeAssertion() {
	fmt.Println("2. Type Assertion with Comma-Ok:")

	values := []interface{}{
		"hello",
		42,
		true,
		3.14,
	}

	for _, val := range values {
		processValue(val)
	}
	fmt.Println()
}

// Example 3: Channel closed detection
func demonstrateChannelClosed() {
	fmt.Println("3. Channel Closed Detection:")

	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	// Read all values and detect closure
	for {
		val, ok := <-ch
		if !ok {
			fmt.Println("   ✓ Channel closed")
			break
		}
		fmt.Printf("   Received: %d\n", val)
	}
	fmt.Println()
}

// Example 4: Channel vs range
func demonstrateChannelRange() {
	fmt.Println("4. Channel with Range (comma-ok built-in):")

	ch := make(chan string, 3)
	ch <- "first"
	ch <- "second"
	ch <- "third"
	close(ch)

	// Range automatically uses comma-ok
	for msg := range ch {
		fmt.Printf("   Got: %s\n", msg)
	}
	fmt.Println("   Range loop ended (channel closed)")
	fmt.Println()
}

// Example 5: Real-world user lookup
type User struct {
	Name  string
	Email string
}

func getUserByID(id int) (*User, error) {
	users := map[int]*User{
		1: {Name: "Alice", Email: "alice@example.com"},
		2: {Name: "Bob", Email: "bob@example.com"},
	}

	user, ok := users[id]
	if !ok {
		return nil, fmt.Errorf("user %d not found", id)
	}

	return user, nil
}

func demonstrateRealWorld() {
	fmt.Println("5. Real-World User Lookup:")

	// Successful lookup
	if user, err := getUserByID(1); err == nil {
		fmt.Printf("   ✓ Found: %s (%s)\n", user.Name, user.Email)
	}

	// Failed lookup
	if user, err := getUserByID(999); err != nil {
		fmt.Printf("   ✗ Error: %v\n", err)
	} else {
		fmt.Printf("   Found: %s\n", user.Name)
	}
	fmt.Println()
}

// Example 6: Type switch (advanced comma-ok)
func identifyType(val interface{}) {
	switch v := val.(type) {
	case string:
		fmt.Printf("   String of length %d: %q\n", len(v), v)
	case int:
		fmt.Printf("   Integer: %d\n", v)
	case bool:
		fmt.Printf("   Boolean: %v\n", v)
	case []int:
		fmt.Printf("   Int slice of length %d\n", len(v))
	default:
		fmt.Printf("   Unknown type: %T\n", v)
	}
}

func demonstrateTypeSwitch() {
	fmt.Println("6. Type Switch (built on comma-ok):")

	values := []interface{}{
		"Go",
		2024,
		true,
		[]int{1, 2, 3},
		3.14,
	}

	for _, val := range values {
		identifyType(val)
	}
	fmt.Println()
}

// Example 7: Channel timeout pattern
func demonstrateChannelTimeout() {
	fmt.Println("7. Channel Timeout Pattern:")

	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "data arrived"
	}()

	select {
	case msg, ok := <-ch:
		if ok {
			fmt.Printf("   ✓ Received: %s\n", msg)
		}
	case <-time.After(1 * time.Second):
		fmt.Println("   ✗ Timeout after 1 second")
	}
	fmt.Println()
}

// Example 8: Multiple type attempts
func parseValue(val interface{}) string {
	// Try multiple types in order
	if s, ok := val.(string); ok {
		return fmt.Sprintf("string: %s", s)
	}
	if n, ok := val.(int); ok {
		return fmt.Sprintf("int: %d", n)
	}
	if f, ok := val.(float64); ok {
		return fmt.Sprintf("float: %.2f", f)
	}
	return fmt.Sprintf("unknown: %T", val)
}

func demonstrateMultipleAttempts() {
	fmt.Println("8. Multiple Type Attempts:")

	values := []interface{}{"hello", 42, 3.14, true}

	for _, val := range values {
		result := parseValue(val)
		fmt.Printf("   %v → %s\n", val, result)
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== Comma-Ok Idiom Examples ===\n")

	demonstrateMapAccess()
	demonstrateTypeAssertion()
	demonstrateChannelClosed()
	demonstrateChannelRange()
	demonstrateRealWorld()
	demonstrateTypeSwitch()
	demonstrateChannelTimeout()
	demonstrateMultipleAttempts()

	fmt.Println("Key Takeaway:")
	fmt.Println("value, ok := operation")
	fmt.Println("- Maps: check existence")
	fmt.Println("- Channels: detect closure")
	fmt.Println("- Type assertions: prevent panics")

	fmt.Println("\n=== Done ===")
}
