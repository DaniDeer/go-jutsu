# go-jutsu Agent Instructions

Repository for Go patterns, idioms, techniques discovered while learning. Keep compact, clear, practical.

## Repository Structure

```
go-jutsu/
├── pattern-name-1/
│   ├── pattern.md     # Explanation, use cases, gotchas
│   └── example.go     # Runnable demonstration
├── pattern-name-2/
│   ├── pattern.md
│   └── example.go
└── README.md          # Main index with all patterns
```

Each pattern is self-contained in its own directory.

## Adding New Patterns

When user discovers/requests new pattern:

1. **Create pattern directory**: `<pattern-name>/`
   - Use kebab-case (e.g., `worker-pool-pattern`)
   - Pattern name should be descriptive and searchable

2. **Create pattern.md**:
   - Title with clear description
   - "What it is (and isn't)" section
   - "Where you see it" with use cases
   - "Real Example" with working code
   - "Gotchas" or "Common mistakes" if relevant
   - Keep concise - respect reader's time

3. **Create example.go**:
   - Standalone, runs with `go run example.go`
   - Demonstrate core concept
   - Include comments explaining key parts
   - Show multiple variations if helpful
   - Keep under 100 lines when possible

4. **Update main README.md**:
   - Add link to pattern directory in appropriate category
   - Keep patterns organized by category
   - Maintain consistent formatting

## Maintaining Consistency

**pattern.md template:**

````markdown
# Pattern Name

Brief one-liner about the pattern.

## What It Is (and Isn't)

Clear explanation. Compare to similar concepts if needed.

## Where You See It

- Use case 1
- Use case 2
- Use case 3

## Real Example

```go
// Working code here
```

## Gotchas

Things to watch out for (optional).
````

**example.go template:**

```go
package main

// Brief description of what this demonstrates

import (...)

// Core pattern implementation with comments
func patternExample() { ... }

func main() {
    // Demonstrate usage
}
```

// Minimal example
\```

## Pattern Name 2

\```go
// Minimal example
\```

````

## Sync Rules

When updating content:

1. **Pattern update** → Check if example needs update
2. **New pattern** → Create directory with pattern.md + example.go, update main README
3. **Example update** → Sync code snippets in pattern.md
4. **Pattern deletion** → Remove directory and update main README

## Philosophy

- **Practical over theoretical** - Real use cases only
- **Clear over clever** - No hand-waving
- **Complete over partial** - Working examples required
- **Compact over comprehensive** - Respect time
- **Idiomatic over exotic** - Standard Go practices

## Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Comment the "why", not the "what"
- Keep examples focused on one concept
- Avoid external dependencies when possible

## When Adding Content

Ask yourself:

- Is this actually useful in real Go code?
- Can I explain it clearly in <500 words?
- Does the example run without setup?
- Would past-me find this helpful?

If all yes → add it.
If any no → rethink or simplify.

## Maintenance Commands

Use `maintain.sh` for repo maintenance tasks:

**After adding/updating patterns:**

```bash
./maintain.sh validate   # Check all patterns have pattern.md + example.go
./maintain.sh check      # Verify all examples compile
```

**Before committing:**

```bash
./maintain.sh format     # Format all Go code with gofmt
./maintain.sh validate   # Final structure check
```

**Testing examples:**

```bash
./maintain.sh run-all    # Run all examples (2s timeout each)
# Or run specific example:
cd <pattern-name>/
go run example.go
```

**Available commands:**

- `check` - Compile all examples individually
- `format` - Run gofmt on all code
- `validate` - Check patterns have both pattern.md and example.go
- `run-all` - Execute all examples with timeout
- `help` - Show command reference

## Workflow Summary

1. Create pattern directory: `mkdir <pattern-name>/`
2. Add `pattern.md` and `example.go` in the directory
3. Update main `README.md` with link to pattern
4. Run `./maintain.sh validate` to verify structure
5. Run `./maintain.sh check` to test compilation
6. Run `./maintain.sh format` before commit
````
