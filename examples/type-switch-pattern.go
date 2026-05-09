package main

import (
	"fmt"
)

func processValue(val interface{}) {
	switch v := val.(type) {
	case string:
		fmt.Printf("   String (len=%d): %q\n", len(v), v)
	case int:
		fmt.Printf("   Int: %d\n", v)
	case []int:
		fmt.Printf("   Int slice (len=%d): %v\n", len(v), v)
	case bool:
		fmt.Printf("   Bool: %v\n", v)
	case nil:
		fmt.Println("   Nil value")
	default:
		fmt.Printf("   Unknown type: %T = %v\n", v, v)
	}
}

func main() {
	fmt.Println("=== Type Switch Pattern ===\n")

	fmt.Println("1. Basic Type Switch:")
	values := []interface{}{
		"hello",
		42,
		[]int{1, 2, 3},
		true,
		3.14,
		nil,
	}

	for _, val := range values {
		processValue(val)
	}
	fmt.Println()

	fmt.Println("2. Type Variable Scope:")
	var x interface{} = "test"
	switch v := x.(type) {
	case string:
		fmt.Printf("   Inside case: v=%q (type: string)\n", v)
	}
	// v not accessible here
	fmt.Println()

	fmt.Println("3. Multiple Types in One Case:")
	checkNumeric := func(val interface{}) {
		switch val.(type) {
		case int, int32, int64, float32, float64:
			fmt.Println("   Numeric type")
		default:
			fmt.Println("   Non-numeric type")
		}
	}
	checkNumeric(42)
	checkNumeric(3.14)
	checkNumeric("text")
	fmt.Println()

	fmt.Println("4. Interface Type Check:")

	// Define interface and type
	type Stringer interface {
		String() string
	}

	type MyType struct {
		name string
	}

	getStringMethod := func(m MyType) string {
		return m.name
	}

	checkStringer := func(val interface{}) {
		switch v := val.(type) {
		case Stringer:
			fmt.Printf("   Stringer: %s\n", v.String())
		default:
			fmt.Printf("   Not a Stringer: %T\n", val)
		}
	}

	// Create value
	mt := MyType{name: "example"}
	fmt.Printf("   MyType.name: %s\n", getStringMethod(mt))
	checkStringer(mt)
	checkStringer(42)
	fmt.Println()

	fmt.Println("Key Takeaway:")
	fmt.Println("Type switch for handling multiple types")
	fmt.Println("- switch v := x.(type)")
	fmt.Println("- Type-specific variable in each case")
	fmt.Println("- Idiomatic for interface{} handling")

	fmt.Println("\n=== Done ===")
}
