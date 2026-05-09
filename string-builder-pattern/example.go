package main

import (
	"fmt"
	"strings"
	"time"
)

// Example 1: Basic string building
func BuildGreeting(names []string) string {
	var b strings.Builder
	b.WriteString("Hello, ")

	for i, name := range names {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(name)
	}

	b.WriteString("!")
	return b.String()
}

// Example 2: SQL query builder
func BuildInsertSQL(table string, cols []string, rows [][]interface{}) string {
	var b strings.Builder
	b.Grow(256) // Preallocate

	b.WriteString("INSERT INTO ")
	b.WriteString(table)
	b.WriteString(" (")

	for i, col := range cols {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(col)
	}

	b.WriteString(") VALUES ")

	for i, row := range rows {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("(")
		for j, val := range row {
			if j > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "'%v'", val) // Builder implements io.Writer
		}
		b.WriteString(")")
	}

	return b.String()
}

// Example 3: HTML builder
func BuildHTML(title string, items []string) string {
	var b strings.Builder

	b.WriteString("<html><head><title>")
	b.WriteString(title)
	b.WriteString("</title></head><body><ul>")

	for _, item := range items {
		b.WriteString("<li>")
		b.WriteString(item)
		b.WriteString("</li>")
	}

	b.WriteString("</ul></body></html>")
	return b.String()
}

// Benchmark comparison
func concatenateWithPlus(n int) string {
	var s string
	for i := 0; i < n; i++ {
		s += "x"
	}
	return s
}

func concatenateWithBuilder(n int) string {
	var b strings.Builder
	b.Grow(n) // Optimize if size known
	for i := 0; i < n; i++ {
		b.WriteString("x")
	}
	return b.String()
}

func main() {
	// Example 1: Basic greeting
	fmt.Println("=== Example 1: Basic Builder ===")
	names := []string{"Alice", "Bob", "Charlie"}
	greeting := BuildGreeting(names)
	fmt.Println(greeting)

	// Example 2: SQL generation
	fmt.Println("\\n=== Example 2: SQL Builder ===")
	sql := BuildInsertSQL(
		"users",
		[]string{"name", "age", "email"},
		[][]interface{}{
			{"Alice", 30, "alice@example.com"},
			{"Bob", 25, "bob@example.com"},
		},
	)
	fmt.Println(sql)

	// Example 3: HTML generation
	fmt.Println("\\n=== Example 3: HTML Builder ===")
	html := BuildHTML("My Page", []string{"Item 1", "Item 2", "Item 3"})
	fmt.Println(html)

	// Example 4: Performance comparison
	fmt.Println("\\n=== Example 4: Performance Comparison ===")
	n := 10000

	start := time.Now()
	_ = concatenateWithPlus(n)
	plusDuration := time.Since(start)

	start = time.Now()
	_ = concatenateWithBuilder(n)
	builderDuration := time.Since(start)

	fmt.Printf("Concatenate %d strings:\\n", n)
	fmt.Printf("  Using +: %v\\n", plusDuration)
	fmt.Printf("  Using Builder: %v\\n", builderDuration)
	fmt.Printf("  Speedup: %.1fx\\n",
		float64(plusDuration)/float64(builderDuration))

	// Example 5: Grow optimization
	fmt.Println("\\n=== Example 5: Grow() Optimization ===")

	start = time.Now()
	var b1 strings.Builder
	for i := 0; i < 1000; i++ {
		b1.WriteString("data")
	}
	_ = b1.String()
	noGrow := time.Since(start)

	start = time.Now()
	var b2 strings.Builder
	b2.Grow(4000) // Preallocate
	for i := 0; i < 1000; i++ {
		b2.WriteString("data")
	}
	_ = b2.String()
	withGrow := time.Since(start)

	fmt.Printf("Without Grow: %v\\n", noGrow)
	fmt.Printf("With Grow: %v\\n", withGrow)

	// Example 6: fmt.Fprintf with Builder
	fmt.Println("\\n=== Example 6: fmt.Fprintf Support ===")
	var b strings.Builder
	fmt.Fprintf(&b, "User: %s (ID: %d)\\n", "Alice", 123)
	fmt.Fprintf(&b, "Status: %s\\n", "active")
	fmt.Print(b.String())
}
