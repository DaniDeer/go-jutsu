# go-jutsu Agent Instructions

Repository for Go patterns, idioms, techniques discovered while learning. Keep compact, clear, practical.

## Repository Structure

```
go-jutsu/
├── patterns/          # Deep explanations of patterns
├── examples/          # Runnable Go code demonstrating patterns
├── cheatsheets/       # Quick reference cards
└── README.md          # Main index
```

## Adding New Patterns

When user discovers/requests new pattern:

1. **Create pattern file**: `patterns/<pattern-name>.md`
   - Title with clear description
   - "What it is (and isn't)" section
   - "Where you see it" with use cases
   - "Real Example" with working code
   - "Gotchas" or "Common mistakes" if relevant
   - Keep concise - respect reader's time

2. **Create runnable example**: `examples/<pattern-name>.go`
   - Standalone, runs with `go run`
   - Demonstrate core concept
   - Include comments explaining key parts
   - Show multiple variations if helpful
   - Keep under 100 lines when possible

3. **Update indexes**:
   - Add to `patterns/README.md` pattern list
   - Add to `examples/README.md` example list
   - Add to main `README.md` patterns section
   - Keep alphabetical or logical grouping

4. **Optional cheatsheet**: `cheatsheets/<topic>.md`
   - Only if pattern is complex or has variants
   - Single page, printable format
   - Quick reference tables/code snippets
   - No long explanations

## Maintaining Consistency

**Pattern file template:**

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

\```go
// Working code here
\```

## Gotchas

Things to watch out for (optional).
````

**Example file template:**

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

**Cheatsheet template:**

````markdown
# Topic Quick Reference

## Pattern Name 1

\```go
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
2. **New pattern** → Always create example + update all READMEs
3. **Example update** → Sync code snippets in pattern docs
4. **Pattern deletion** → Remove from all READMEs + delete example

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
./maintain.sh validate   # Check pattern↔example sync
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
go run examples/<pattern-name>.go
```

**Available commands:**

- `check` - Compile all examples individually
- `format` - Run gofmt on all code
- `validate` - Check patterns have matching examples
- `run-all` - Execute all examples with timeout
- `help` - Show command reference

## Workflow Summary

1. Add pattern + example files
2. Update all READMEs
3. Run `./maintain.sh validate` to verify sync
4. Run `./maintain.sh check` to test compilation
5. Run `./maintain.sh format` before commit
