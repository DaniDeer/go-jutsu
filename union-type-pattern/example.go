package main

// Union type pattern: two Go approximations of union types.
//
//  1. Runtime union — any + type switch (open set, default required)
//  2. Generic constraint union — ~A | ~B (compile-time, type parameters only)
//
// Compare with sum-type-pattern: sealed interface = closed set, no default needed.

import "fmt"

// ── 1. Runtime union ────────────────────────────────────────────────────────

// ConfigValue represents any scalar value that may appear in a config file.
// Using `any` models an open-set union: callers can pass any type.
type ConfigValue = any

// describe handles the common config types but must include a default —
// the set is open, so unknown types are always possible.
func describe(v ConfigValue) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("string(%q)", val)
	case int:
		return fmt.Sprintf("int(%d)", val)
	case bool:
		return fmt.Sprintf("bool(%v)", val)
	case []string:
		return fmt.Sprintf("[]string%v", val)
	default:
		// default is essential — any type could arrive at runtime
		return fmt.Sprintf("unknown(%T: %v)", val, val)
	}
}

// ── 2. Generic constraint union ─────────────────────────────────────────────

// Number is a type-set union used as a generic constraint.
// It is NOT a value type — you cannot declare `var x Number`.
// It restricts which types may be passed as a type argument.
type Number interface {
	~int | ~int64 | ~float64
}

// Sum works on any slice whose element type satisfies Number.
// The compiler rejects Sum[string](...) — string is not in the type set.
func Sum[N Number](vals []N) N {
	var total N
	for _, v := range vals {
		total += v
	}
	return total
}

// Max returns the larger of two values for any ordered numeric type.
func Max[N Number](a, b N) N {
	if a > b {
		return a
	}
	return b
}

// ── main ────────────────────────────────────────────────────────────────────

func main() {
	// --- Runtime union: any config values ---
	fmt.Println("=== Runtime union (any) ===")
	configs := []ConfigValue{
		"localhost",
		8080,
		true,
		[]string{"alice", "bob"},
		3.14, // float64: hits the default case
	}
	for _, v := range configs {
		fmt.Println(" ", describe(v))
	}

	fmt.Println()

	// --- Generic constraint union ---
	fmt.Println("=== Generic constraint union (Number) ===")

	ints := []int{1, 2, 3, 4, 5}
	fmt.Printf("  Sum[int]     = %d\n", Sum(ints))

	floats := []float64{1.1, 2.2, 3.3}
	fmt.Printf("  Sum[float64] = %.1f\n", Sum(floats))

	fmt.Printf("  Max[int64]   = %d\n", Max[int64](100, 200))

	// Uncommenting the line below would be a compile error:
	// Sum([]string{"a", "b"})  // string does not satisfy Number
}
