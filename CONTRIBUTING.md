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

1. **Create pattern file** in `patterns/`
   - Use template from `.github/copilot-instructions.md`
   - Keep under 500 words
   - Include working code

2. **Create example** in `examples/`
   - Must run standalone
   - Add to `examples/README.md`

3. **Update main README**
   - Add to patterns list
   - Keep formatting consistent

4. **Open PR**
   - Title: "Add [pattern name] pattern"
   - Describe why pattern is useful

## Code Guidelines

- Run `gofmt` before committing
- Keep examples under 100 lines
- Use standard library when possible
- Comment the "why", not "what"

## Questions?

Open an issue. Keep it simple.
