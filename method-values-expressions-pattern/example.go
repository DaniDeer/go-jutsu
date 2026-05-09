package main

import (
	"fmt"
	"strings"
)

// Example type
type Counter struct {
	count int
	name  string
}

func (c *Counter) Inc() {
	c.count++
}

func (c *Counter) Add(n int) {
	c.count += n
}

func (c *Counter) Get() int {
	return c.count
}

func main() {
	fmt.Println("=== Method Values vs Expressions ===\n")

	// Example 1: Method value
	fmt.Println("1. Method Value (bound to instance):")
	c1 := &Counter{name: "counter1"}
	inc := c1.Inc // Closure over c1

	inc()
	inc()
	fmt.Printf("   c1.count: %d\n", c1.Get())
	fmt.Println("   ✓ Method value captured c1")
	fmt.Println()

	// Example 2: Method expression
	fmt.Println("2. Method Expression (unbound):")
	c2 := &Counter{name: "counter2"}
	incFunc := (*Counter).Inc // func(*Counter)

	incFunc(c2)
	incFunc(c2)
	fmt.Printf("   c2.count: %d\n", c2.Get())
	fmt.Println("   ✓ Method expression requires receiver")
	fmt.Println()

	// Example 3: Different instances
	fmt.Println("3. Method Value Captures Specific Instance:")
	c3 := &Counter{name: "c3"}
	c4 := &Counter{name: "c4"}

	inc3 := c3.Inc
	inc4 := c4.Inc

	inc3()
	inc4()
	inc4()

	fmt.Printf("   c3.count: %d\n", c3.Get())
	fmt.Printf("   c4.count: %d\n", c4.Get())
	fmt.Println()

	// Example 4: Using in higher-order functions
	fmt.Println("4. Method Values as Callbacks:")

	counters := []*Counter{
		{name: "a"},
		{name: "b"},
		{name: "c"},
	}

	// Apply Inc to each
	for _, c := range counters {
		c.Inc()
	}

	for _, c := range counters {
		fmt.Printf("   %s: %d\n", c.name, c.Get())
	}
	fmt.Println()

	// Example 5: Method expression with different receivers
	fmt.Println("5. Method Expression with Multiple Types:")

	type MyString string

	toUpper := strings.ToUpper
	s1 := toUpper("hello")
	s2 := toUpper("world")

	fmt.Printf("   %s, %s\n", s1, s2)
	fmt.Println("   ✓ Function value (similar concept)")
	fmt.Println()

	fmt.Println("Key Takeaway:")
	fmt.Println("Method value: obj.Method - bound to obj")
	fmt.Println("Method expression: (*Type).Method - needs receiver")
	fmt.Println("Use for callbacks, delayed execution")

	fmt.Println("\n=== Done ===")
}
