package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Compose combines two functions: f(g(x))
func Compose[A, B, C any](f func(B) C, g func(A) B) func(A) C {
	return func(x A) C {
		return f(g(x))
	}
}

// Transform is a function that transforms T to T (enables chaining)
type Transform[T any] func(T) T

// Then chains transformations
func (t Transform[T]) Then(next Transform[T]) Transform[T] {
	return func(x T) T {
		return next(t(x))
	}
}

func main() {
	fmt.Println("=== Function Composition Pattern ===\n")

	// Example 1: Basic composition
	fmt.Println("1. Basic Composition:")

	add5 := func(x int) int { return x + 5 }
	double := func(x int) int { return x * 2 }

	// Compose: double(add5(x))
	add5ThenDouble := Compose(double, add5)
	result1 := add5ThenDouble(3)
	fmt.Printf("   add5ThenDouble(3) = %d\n", result1)
	fmt.Printf("   (3 + 5) * 2 = %d\n", result1)
	fmt.Println()

	// Example 2: Transform chaining
	fmt.Println("2. Transform Chaining:")

	square := Transform[int](func(x int) int { return x * x })
	add10 := Transform[int](func(x int) int { return x + 10 })
	half := Transform[int](func(x int) int { return x / 2 })

	pipeline := add10.Then(square).Then(half)
	result2 := pipeline(5)
	fmt.Printf("   Pipeline(5) = %d\n", result2)
	fmt.Printf("   ((5 + 10)² / 2) = %d\n", result2)
	fmt.Println()

	// Example 3: String processing
	fmt.Println("3. String Processing Pipeline:")

	type StringTransform func(string) string

	toLower := StringTransform(strings.ToLower)
	trim := StringTransform(strings.TrimSpace)
	removePunct := StringTransform(func(s string) string {
		return strings.Map(func(r rune) rune {
			if unicode.IsPunct(r) {
				return -1
			}
			return r
		}, s)
	})

	normalize := func(s string) string {
		s = trim(s)
		s = toLower(s)
		s = removePunct(s)
		return s
	}

	input := "  Hello, World!  "
	output := normalize(input)
	fmt.Printf("   Input:  %q\n", input)
	fmt.Printf("   Output: %q\n", output)
	fmt.Println()

	// Example 4: Multiple compositions
	fmt.Println("4. Multiple Compositions:")

	inc := func(x int) int { return x + 1 }
	triple := func(x int) int { return x * 3 }
	negate := func(x int) int { return -x }

	// Build different pipelines
	f1 := Compose(triple, inc)                  // triple(inc(x))
	f2 := Compose(negate, f1)                   // negate(triple(inc(x)))
	f3 := Compose(Compose(negate, triple), inc) // Same as f2

	x := 5
	fmt.Printf("   x = %d\n", x)
	fmt.Printf("   triple(inc(x)) = %d\n", f1(x))
	fmt.Printf("   negate(triple(inc(x))) = %d\n", f2(x))
	fmt.Printf("   Same result: %d\n", f3(x))
	fmt.Println()

	// Example 5: Validation chain
	fmt.Println("5. Validation Chain:")

	type Validator func(string) bool

	notEmpty := Validator(func(s string) bool {
		return len(strings.TrimSpace(s)) > 0
	})

	minLength := func(min int) Validator {
		return func(s string) bool {
			return len(s) >= min
		}
	}

	hasUpper := Validator(func(s string) bool {
		for _, r := range s {
			if unicode.IsUpper(r) {
				return true
			}
		}
		return false
	})

	// Combine validators
	validate := func(s string) bool {
		return notEmpty(s) && minLength(5)(s) && hasUpper(s)
	}

	tests := []string{"Hello", "hi", "", "WORLD", "test"}
	for _, test := range tests {
		fmt.Printf("   %q valid: %v\n", test, validate(test))
	}
	fmt.Println()

	// Example 6: Currying with composition
	fmt.Println("6. Currying with Composition:")

	add := func(a int) func(int) int {
		return func(b int) int {
			return a + b
		}
	}

	multiply := func(a int) func(int) int {
		return func(b int) int {
			return a * b
		}
	}

	add10 = add(10)
	multiplyBy3 := multiply(3)

	composed := Compose(multiplyBy3, add10)
	result6 := composed(5)
	fmt.Printf("   Compose(×3, +10)(5) = %d\n", result6)
	fmt.Printf("   (5 + 10) × 3 = %d\n", result6)
	fmt.Println()

	fmt.Println("Key Takeaway:")
	fmt.Println("Function composition:")
	fmt.Println("- Build complex functions from simple ones")
	fmt.Println("- Compose(f, g)(x) = f(g(x))")
	fmt.Println("- Reusable, testable building blocks")
	fmt.Println("- Functional programming in Go")

	fmt.Println("\n=== Done ===")
}
