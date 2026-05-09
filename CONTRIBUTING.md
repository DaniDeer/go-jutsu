# Contributing to go-jutsu

Found a cool Go pattern? Want to add an example? Contributions welcome!

## What to Contribute

**Good candidates:**

- Patterns you wish you knew earlier
- Non-obvious Go idioms
- Techniques from real projects
- Clarifications of confusing concepts

**Not a good fit:**

- Basic Go syntax tutorials
- Copy-paste from official docs
- Language debate/opinions
- Framework-specific patterns

## How to Add a Pattern

1. **Create pattern directory**

   ```bash
   mkdir pattern-name/
   ```

   - Use kebab-case (e.g., `worker-pool-pattern`)
   - Name should be descriptive and searchable

2. **Create pattern.md**
   - Use template from `.github/copilot-instructions.md`
   - Keep under 500 words
   - Include working code examples

3. **Create example.go**
   - Must run standalone: `go run example.go`
   - Demonstrates core concept clearly
   - Keep under 100 lines when possible

4. **Update main README.md**
   - Add link to your pattern directory
   - Place in appropriate category
   - Keep formatting consistent

5. **Validate your work**

   ```bash
   cd pattern-name/
   go run example.go              # Test it works
   cd ..
   ./maintain.sh validate         # Check structure
   ./maintain.sh check            # Verify compilation
   ./maintain.sh format           # Format code
   ```

6. **Open PR**
   - Title: "Add [pattern name] pattern"
   - Describe why pattern is useful
   - Show example output if relevant

## Code Guidelines

- Run `gofmt` before committing (or use `./maintain.sh format`)
- Keep examples under 100 lines when possible
- Use standard library when possible
- Comment the "why", not "what"
- Each pattern must be self-contained in its directory
- Examples must run with `go run example.go` from pattern directory

## Repository Structure

Each pattern lives in its own directory:

```
pattern-name/
├── pattern.md    # Documentation with explanation and gotchas
└── example.go    # Runnable code demonstrating the pattern
```

Browse existing patterns for reference before adding new ones.

## Questions?

Open an issue. Keep it simple.
