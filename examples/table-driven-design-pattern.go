package main

import (
	"fmt"
	"strings"
)

// Example 1: Calculator with table-driven tests
func Add(a, b int) int { return a + b }
func Sub(a, b int) int { return a - b }
func Mul(a, b int) int { return a * b }
func Div(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

func testCalculator() {
	fmt.Println("1. Table-Driven Calculator Tests:")

	tests := []struct {
		name string
		a, b int
		op   func(int, int) int
		want int
	}{
		{"add positive", 2, 3, Add, 5},
		{"add negative", -1, -2, Add, -3},
		{"sub", 10, 3, Sub, 7},
		{"mul", 4, 5, Mul, 20},
		{"div", 15, 3, Div, 5},
		{"div by zero", 10, 0, Div, 0},
	}

	for _, tt := range tests {
		got := tt.op(tt.a, tt.b)
		status := "✓"
		if got != tt.want {
			status = "✗"
		}
		fmt.Printf("   %s %s: got %d, want %d\n", status, tt.name, got, tt.want)
	}
	fmt.Println()
}

// Example 2: Validation rules
type Validator struct {
	name string
	rule func(string) bool
	msg  string
}

var validators = []Validator{
	{"not empty", func(s string) bool { return len(strings.TrimSpace(s)) > 0 }, "must not be empty"},
	{"min length", func(s string) bool { return len(s) >= 3 }, "must be at least 3 characters"},
	{"no spaces", func(s string) bool { return !strings.Contains(s, " ") }, "must not contain spaces"},
}

func validate(input string) []string {
	var errors []string
	for _, v := range validators {
		if !v.rule(input) {
			errors = append(errors, v.msg)
		}
	}
	return errors
}

func testValidation() {
	fmt.Println("2. Table-Driven Validation:")

	inputs := []string{"hello", "hi", "", "hello world", "ok"}

	for _, input := range inputs {
		errors := validate(input)
		if len(errors) == 0 {
			fmt.Printf("   ✓ %q is valid\n", input)
		} else {
			fmt.Printf("   ✗ %q: %s\n", input, strings.Join(errors, ", "))
		}
	}
	fmt.Println()
}

// Example 3: State machine
type State int

const (
	StatePending State = iota
	StateRunning
	StateDone
	StateFailed
)

func (s State) String() string {
	return [...]string{"Pending", "Running", "Done", "Failed"}[s]
}

var transitions = map[State][]State{
	StatePending: {StateRunning, StateFailed},
	StateRunning: {StateDone, StateFailed},
	StateDone:    {},
	StateFailed:  {StatePending}, // Can retry
}

func canTransition(from, to State) bool {
	allowed := transitions[from]
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func testStateMachine() {
	fmt.Println("3. Table-Driven State Machine:")

	tests := []struct {
		from, to State
		valid    bool
	}{
		{StatePending, StateRunning, true},
		{StatePending, StateDone, false},
		{StateRunning, StateDone, true},
		{StateRunning, StateFailed, true},
		{StateDone, StatePending, false},
		{StateFailed, StatePending, true},
	}

	for _, tt := range tests {
		got := canTransition(tt.from, tt.to)
		status := "✓"
		if got != tt.valid {
			status = "✗"
		}
		fmt.Printf("   %s %s -> %s: %v (expected %v)\n",
			status, tt.from, tt.to, got, tt.valid)
	}
	fmt.Println()
}

// Example 4: HTTP status handling
type StatusInfo struct {
	message string
	retry   bool
}

var statusCodes = map[int]StatusInfo{
	200: {"OK", false},
	201: {"Created", false},
	400: {"Bad Request", false},
	404: {"Not Found", false},
	500: {"Server Error", true},
	503: {"Service Unavailable", true},
}

func handleStatus(code int) string {
	info, ok := statusCodes[code]
	if !ok {
		return fmt.Sprintf("Unknown status: %d", code)
	}

	msg := fmt.Sprintf("%d %s", code, info.message)
	if info.retry {
		msg += " (will retry)"
	}
	return msg
}

func testStatusHandling() {
	fmt.Println("4. Table-Driven Status Codes:")

	codes := []int{200, 404, 500, 503, 999}

	for _, code := range codes {
		result := handleStatus(code)
		fmt.Printf("   %s\n", result)
	}
	fmt.Println()
}

// Example 5: Data transformation
type Transform struct {
	name string
	fn   func(string) string
}

var transforms = []Transform{
	{"lowercase", strings.ToLower},
	{"uppercase", strings.ToUpper},
	{"trim", strings.TrimSpace},
	{"reverse", func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}},
}

func testTransforms() {
	fmt.Println("5. Table-Driven Transforms:")

	input := "  Hello World  "
	fmt.Printf("   Input: %q\n", input)

	for _, t := range transforms {
		result := t.fn(input)
		fmt.Printf("   %s: %q\n", t.name, result)
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== Table-Driven Design Pattern ===\n")

	testCalculator()
	testValidation()
	testStateMachine()
	testStatusHandling()
	testTransforms()

	fmt.Println("Key Takeaway:")
	fmt.Println("Table-driven design:")
	fmt.Println("- Data structures drive behavior")
	fmt.Println("- Idiomatic in Go (especially tests)")
	fmt.Println("- Easy to add cases")
	fmt.Println("- Self-documenting")
	fmt.Println("- DRY and maintainable")

	fmt.Println("\n=== Done ===")
}
