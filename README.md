# go-jutsu 🥷

Advanced Go patterns and techniques discovered while mastering the way of the Gopher.

## What is this?

Collection of Go idioms, patterns, and techniques that aren't always obvious when learning Go. Each pattern includes clear explanations, real-world examples, and gotchas to avoid.

## Patterns

- [Closure Pattern (Manual "Currying")](./patterns/go-closure-currying-pattern.md) - HTTP middleware factory pattern using closures
- [Defer Cleanup Pattern](./patterns/defer-cleanup-pattern.md) - Stack-based cleanup that runs on exit, even on panic (Go-specific)
- [Empty Struct Signal Pattern](./patterns/empty-struct-signal-pattern.md) - Zero-byte signaling with `struct{}` (Go-specific)

## Structure

```
go-jutsu/
├── patterns/          # Pattern explanations and guides
├── examples/          # Runnable code examples
└── cheatsheets/       # Quick reference cards
```

## Philosophy

Keep it simple. Each pattern should be:

- **Practical** - real-world use cases
- **Clear** - no hand-waving
- **Complete** - working examples included
- **Compact** - respect your time

## Maintenance

Helper script for keeping repo clean:

```bash
./maintain.sh check      # Verify all examples compile
./maintain.sh format     # Format all code
./maintain.sh validate   # Check structure consistency
```

See [`.github/copilot-instructions.md`](./.github/copilot-instructions.md) for agent guide on adding patterns.

## Contributing

Found a cool pattern? See [CONTRIBUTING.md](./CONTRIBUTING.md). PRs welcome!
