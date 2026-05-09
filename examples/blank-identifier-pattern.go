package main

import (
	"fmt"
	"io"
)

// Example types for interface verification
type MyWriter struct{}

func (m *MyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Compile-time verification
var _ io.Writer = (*MyWriter)(nil)

func main() {
	fmt.Println("=== Blank Identifier Pattern ===\n")

	// Example 1: Ignore return values
	fmt.Println("1. Ignore Return Values:")
	fmt.Println("   _, err := doThing()  // Ignore first value")
	fmt.Println("   val, _ := map[key]   // Ignore ok boolean")
	fmt.Println()

	// Example 2: Range with blank identifier
	fmt.Println("2. Range Loop:")
	numbers := []int{10, 20, 30}

	fmt.Println("   Ignore index:")
	for _, num := range numbers {
		fmt.Printf("   %d\n", num)
	}

	fmt.Println("\n   Ignore value:")
	for i := range numbers {
		fmt.Printf("   Index: %d\n", i)
	}
	fmt.Println()

	// Example 3: Compile-time interface check
	fmt.Println("3. Compile-Time Interface Verification:")
	fmt.Println("   var _ io.Writer = (*MyWriter)(nil)")
	fmt.Println("   ✓ Verified at compile time (no runtime cost)")
	fmt.Println()

	// Example 4: Multiple blanks
	fmt.Println("4. Multiple Blank Identifiers:")

	multiReturn := func() (int, string, bool) {
		return 1, "two", true
	}

	_, str, _ := multiReturn()
	fmt.Printf("   Got middle value: %s\n", str)
	fmt.Println()

	// Example 5: Unused variable workaround
	fmt.Println("5. Unused Variable During Development:")

	debug := true
	if debug {
		x := "debug info"
		_ = x // Silence "declared and not used" error
		fmt.Println("   var x = ...; _ = x  // Compiler happy")
	}
	fmt.Println()

	// Example 6: Map iteration (keys only)
	fmt.Println("6. Map Keys Only:")

	ages := map[string]int{
		"Alice": 30,
		"Bob":   25,
	}

	for name := range ages { // Value implicitly ignored
		fmt.Printf("   Name: %s\n", name)
	}
	fmt.Println()

	fmt.Println("Key Takeaway:")
	fmt.Println("Blank identifier (_) for:")
	fmt.Println("- Ignoring return values")
	fmt.Println("- Import side effects (database drivers)")
	fmt.Println("- Compile-time interface checks")
	fmt.Println("- Unused variables during development")
	fmt.Println("⚠ Don't ignore errors in production!")

	fmt.Println("\n=== Done ===")
}
