package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FormatJSON formats data as JSON
func FormatJSON(data map[string]interface{}) string {
	b, _ := json.MarshalIndent(data, "", "  ")
	return string(b)
}

// FormatUser formats user as text
func FormatUser(name string, age int) string {
	return fmt.Sprintf("Name: %s\nAge: %d\n", name, age)
}

// Simple test framework
func testCase(name, got, want string) {
	if got == want {
		fmt.Printf("✓ %s\n", name)
	} else {
		fmt.Printf("✗ %s\n  Got:\n%s\n  Want:\n%s\n", name, got, want)
	}
}

func main() {
	fmt.Println("=== Golden Files Pattern ===")

	// Create testdata directory
	os.MkdirAll("testdata", 0755)
	defer os.RemoveAll("testdata")

	// Example 1: JSON formatting
	fmt.Println("\nExample 1: JSON Golden File")
	jsonData := map[string]interface{}{
		"name":  "Alice",
		"age":   30,
		"email": "alice@example.com",
	}
	gotJSON := FormatJSON(jsonData)

	goldenJSON := filepath.Join("testdata", "user.golden.json")

	// Simulate -update flag (write golden file)
	update := true
	if update {
		os.WriteFile(goldenJSON, []byte(gotJSON), 0644)
		fmt.Println("  Golden file created")
	}

	// Read and compare
	wantJSON, _ := os.ReadFile(goldenJSON)
	testCase("JSON format", gotJSON, string(wantJSON))

	// Example 2: Text formatting
	fmt.Println("\nExample 2: Text Golden File")
	gotText := FormatUser("Bob", 25)

	goldenText := filepath.Join("testdata", "user.golden.txt")
	os.WriteFile(goldenText, []byte(gotText), 0644)

	wantText, _ := os.ReadFile(goldenText)
	testCase("Text format", gotText, string(wantText))

	// Example 3: Detect changes
	fmt.Println("\nExample 3: Detect Output Changes")
	modifiedJSON := FormatJSON(map[string]interface{}{
		"name": "Alice",
		"age":  31, // Changed!
	})

	testCase("Modified JSON", modifiedJSON, string(wantJSON))

	// Example 4: Multiple test cases
	fmt.Println("\nExample 4: Multiple Golden Files")
	tests := []struct {
		name string
		data map[string]interface{}
	}{
		{"simple", map[string]interface{}{"key": "value"}},
		{"nested", map[string]interface{}{"a": map[string]int{"b": 1}}},
	}

	for _, tt := range tests {
		got := FormatJSON(tt.data)
		goldenFile := filepath.Join("testdata", tt.name+".golden.json")

		if update {
			os.WriteFile(goldenFile, []byte(got), 0644)
		}

		want, _ := os.ReadFile(goldenFile)
		testCase(tt.name, got, string(want))
	}

	fmt.Println("\n💡 Tip: Run with -update flag to regenerate golden files")
}
