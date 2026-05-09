# go-jutsu 🥷

Advanced Go patterns and techniques discovered while mastering the way of the Gopher.

## What is this?

Collection of Go idioms, patterns, and techniques that aren't always obvious when learning Go. Each pattern includes clear explanations, real-world examples, and gotchas to avoid.

## Patterns

### Foundational Patterns

- [Closure Pattern (Manual "Currying")](./patterns/go-closure-currying-pattern.md) - HTTP middleware factory using closures
- [Defer Cleanup Pattern](./patterns/defer-cleanup-pattern.md) - Stack-based cleanup, LIFO execution (Go-specific)
- [Empty Struct Signal Pattern](./patterns/empty-struct-signal-pattern.md) - Zero-byte signaling with `struct{}` (Go-specific)
- [Zero Values Pattern](./patterns/zero-values-pattern.md) - Useful zero values by design (Go philosophy)

### Gotchas & Idioms

- [Nil Interface Gotcha](./patterns/nil-interface-gotcha.md) - Interface with nil value isn't nil (classic trap)
- [Comma-Ok Idiom](./patterns/comma-ok-idiom.md) - Safe map/channel/type access with two-value returns
- [Blank Identifier Pattern](./patterns/blank-identifier-pattern.md) - Using `_` for ignoring values and side effects

### Concurrency

- [Select Statement Pattern](./patterns/select-statement-pattern.md) - Non-deterministic channel multiplexing
- [Context Propagation Pattern](./patterns/context-propagation-pattern.md) - Cancellation and deadline propagation

### Advanced Patterns

- [Functional Options Pattern](./patterns/functional-options-pattern.md) - Clean configuration with option functions
- [Interface Satisfaction Pattern](./patterns/interface-satisfaction-pattern.md) - Implicit interface implementation
- [Embedding Composition Pattern](./patterns/embedding-composition-pattern.md) - Composition over inheritance
- [Method Values vs Expressions Pattern](./patterns/method-values-expressions-pattern.md) - Methods as closures vs functions
- [Type Switch Pattern](./patterns/type-switch-pattern.md) - Type-specific handling with switch

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
