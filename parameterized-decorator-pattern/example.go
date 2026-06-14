package main

// Demonstrates the parameterized decorator pattern — a decorator factory
// that captures config at setup time and returns a decorator.
//
// Direct Go translation of the Python pattern:
//
//	@to_uppercase
//	@get_truncate(9)
//	def print_input(s): print(s)
//
// Go equivalent: toUppercase(truncate(9)(printLine))

import (
	"fmt"
	"strings"
)

// Processor is a function that consumes a string (like Python's TextFunc).
type Processor func(string)

// --- Plain decorator (no config) ---

// toUppercase wraps p so the input is uppercased before being passed on.
func toUppercase(next Processor) Processor {
	return func(s string) {
		next(strings.ToUpper(s))
	}
}

// --- Decorator factory (parameterized decorator) ---

// truncate returns a decorator that limits input to the first n characters.
// The config (n) is captured at factory-call time, not at execution time.
func truncate(n int) func(Processor) Processor {
	return func(next Processor) Processor {
		return func(s string) {
			if n < len(s) {
				s = s[:n]
			}
			next(s)
		}
	}
}

// --- chain helper for readable stacking ---

// chain applies decorators right-to-left (innermost first), matching Python's
// bottom-up @decorator evaluation order.
func chain(p Processor, decorators ...func(Processor) Processor) Processor {
	for i := len(decorators) - 1; i >= 0; i-- {
		p = decorators[i](p)
	}
	return p
}

func main() {
	printLine := func(s string) { fmt.Println(s) }

	// Explicit stacking — mirrors Python's @to_uppercase @get_truncate(9)
	// truncate(9) is the factory call; it produces the decorator.
	// Innermost (truncate) runs first: "Keep Calm" → toUppercase → "KEEP CALM"
	printInput := toUppercase(truncate(9)(printLine))
	printInput("Keep Calm and Carry On") // → KEEP CALM

	fmt.Println()

	// chain helper — cleaner when stacking more than two decorators
	p := chain(printLine,
		toUppercase, // outermost: runs last
		truncate(5), // innermost: runs first
		toUppercase, // also runs before the outer one
	)
	p("hello world") // truncate(5)→"hello" → toUppercase→"HELLO" → toUppercase→"HELLO"

	fmt.Println()

	// Demonstrate that the factory captures config independently per call.
	short := truncate(4)(printLine)
	long := truncate(20)(printLine)

	short("Go is great") // → "Go i"
	long("Go is great")  // → "Go is great"
}
