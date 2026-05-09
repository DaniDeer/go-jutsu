package main

import (
	"fmt"
	"strings"
)

// Map transforms each element using fn
func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Filter keeps elements where pred returns true
func Filter[T any](slice []T, pred func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce aggregates slice to single value
func Reduce[T, U any](slice []T, init U, fn func(U, T) U) U {
	acc := init
	for _, v := range slice {
		acc = fn(acc, v)
	}
	return acc
}

// Example type
type User struct {
	Name   string
	Age    int
	Active bool
}

func main() {
	fmt.Println("=== Map/Filter/Reduce with Generics ===\n")

	// Example 1: Map
	fmt.Println("1. Map - Transform Elements:")
	numbers := []int{1, 2, 3, 4, 5}
	doubled := Map(numbers, func(x int) int { return x * 2 })
	fmt.Printf("   Original: %v\n", numbers)
	fmt.Printf("   Doubled:  %v\n", doubled)
	fmt.Println()

	// Example 2: Filter
	fmt.Println("2. Filter - Keep Matching:")
	evens := Filter(numbers, func(x int) bool { return x%2 == 0 })
	fmt.Printf("   Original: %v\n", numbers)
	fmt.Printf("   Evens:    %v\n", evens)
	fmt.Println()

	// Example 3: Reduce
	fmt.Println("3. Reduce - Aggregate:")
	sum := Reduce(numbers, 0, func(acc, x int) int { return acc + x })
	product := Reduce(numbers, 1, func(acc, x int) int { return acc * x })
	fmt.Printf("   Numbers: %v\n", numbers)
	fmt.Printf("   Sum:     %d\n", sum)
	fmt.Printf("   Product: %d\n", product)
	fmt.Println()

	// Example 4: Chaining operations
	fmt.Println("4. Chaining - Filter then Map:")
	result := Map(
		Filter(numbers, func(x int) bool { return x > 2 }),
		func(x int) int { return x * 10 },
	)
	fmt.Printf("   Filter > 2, then * 10: %v\n", result)
	fmt.Println()

	// Example 5: Working with structs
	fmt.Println("5. Structs - Users Example:")
	users := []User{
		{Name: "Alice", Age: 30, Active: true},
		{Name: "Bob", Age: 25, Active: false},
		{Name: "Charlie", Age: 35, Active: true},
		{Name: "Diana", Age: 28, Active: false},
	}

	// Get active user names
	activeNames := Map(
		Filter(users, func(u User) bool { return u.Active }),
		func(u User) string { return u.Name },
	)
	fmt.Printf("   Active users: %v\n", activeNames)

	// Total age of active users
	totalAge := Reduce(
		Filter(users, func(u User) bool { return u.Active }),
		0,
		func(sum int, u User) int { return sum + u.Age },
	)
	fmt.Printf("   Total age (active): %d\n", totalAge)
	fmt.Println()

	// Example 6: String operations
	fmt.Println("6. String Operations:")
	words := []string{"hello", "world", "go", "generics"}

	// Uppercase all
	upper := Map(words, strings.ToUpper)
	fmt.Printf("   Uppercase: %v\n", upper)

	// Keep long words
	long := Filter(words, func(s string) bool { return len(s) > 4 })
	fmt.Printf("   Long (>4): %v\n", long)

	// Concatenate
	joined := Reduce(words, "", func(acc, s string) string {
		if acc == "" {
			return s
		}
		return acc + "-" + s
	})
	fmt.Printf("   Joined: %s\n", joined)
	fmt.Println()

	// Example 7: Complex chaining
	fmt.Println("7. Complex Pipeline:")
	pipeline := Reduce(
		Map(
			Filter(numbers, func(x int) bool { return x%2 == 1 }),
			func(x int) int { return x * x },
		),
		0,
		func(acc, x int) int { return acc + x },
	)
	fmt.Printf("   Odd numbers squared and summed: %d\n", pipeline)
	fmt.Printf("   (1² + 3² + 5² = %d)\n", pipeline)
	fmt.Println()

	fmt.Println("Key Takeaway:")
	fmt.Println("Map/Filter/Reduce with generics:")
	fmt.Println("- Type-safe transformations")
	fmt.Println("- Declarative, functional style")
	fmt.Println("- Composable operations")
	fmt.Println("- Not lazy (creates intermediate slices)")

	fmt.Println("\n=== Done ===")
}
