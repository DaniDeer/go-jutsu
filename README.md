# go-jutsu 🥷

Advanced Go patterns and techniques discovered while mastering the way of the Gopher.

## What is this?

Collection of Go idioms, patterns, and techniques that aren't always obvious when learning Go. Each pattern includes clear explanations, real-world examples, and gotchas to avoid.

## Patterns (37 Total)

Each pattern is self-contained in its own directory with documentation (`pattern.md`) and runnable example (`example.go`).

### Foundational Patterns

- [Closure Pattern (Manual "Currying")](./go-closure-currying-pattern/) - HTTP middleware factory using closures
- [Defer Cleanup Pattern](./defer-cleanup-pattern/) - Stack-based cleanup, LIFO execution (Go-specific)
- [Empty Struct Signal Pattern](./empty-struct-signal-pattern/) - Zero-byte signaling with `struct{}` (Go-specific)
- [Zero Values Pattern](./zero-values-pattern/) - Useful zero values by design (Go philosophy)

### Gotchas & Idioms

- [Nil Interface Gotcha](./nil-interface-gotcha/) - Interface with nil value isn't nil (classic trap)
- [Comma-Ok Idiom](./comma-ok-idiom/) - Safe map/channel/type access with two-value returns
- [Blank Identifier Pattern](./blank-identifier-pattern/) - Using `_` for ignoring values and side effects

### Error Handling

- [Sentinel Errors Pattern](./sentinel-errors-pattern/) - Predefined error values with `errors.Is()`
- [Error Wrapping Pattern](./error-wrapping-pattern/) - Add context with `fmt.Errorf("%w")` and `errors.Unwrap/As`
- [Panic and Recover Pattern](./panic-recover-pattern/) - Emergency handling with defer + recover

### Concurrency

- [Select Statement Pattern](./select-statement-pattern/) - Non-deterministic channel multiplexing
- [Context Propagation Pattern](./context-propagation-pattern/) - Cancellation and deadline propagation
- [Worker Pool Pattern](./worker-pool-pattern/) - Fixed goroutines processing job queue
- [Circuit Breaker Pattern](./circuit-breaker-pattern/) - Prevent cascading failures with state machine
- [Rate Limiting Pattern](./rate-limiting-pattern/) - Control request rate with token bucket
- [Bounded Parallelism Pattern](./bounded-parallelism-pattern/) - Semaphore limits concurrent goroutines

### Advanced Patterns

- [Functional Options Pattern](./functional-options-pattern/) - Clean configuration with option functions
- [Interface Satisfaction Pattern](./interface-satisfaction-pattern/) - Implicit interface implementation
- [Embedding Composition Pattern](./embedding-composition-pattern/) - Composition over inheritance
- [Method Values vs Expressions Pattern](./method-values-expressions-pattern/) - Methods as closures vs functions
- [Type Switch Pattern](./type-switch-pattern/) - Type-specific handling with switch
- [Singleton with sync.Once](./singleton-sync-once-pattern/) - Thread-safe lazy initialization
- [Builder Pattern](./builder-pattern/) - Fluent API for complex object construction
- [Factory Pattern](./factory-pattern/) - Registry of constructors for extensibility

### Functional & Declarative

- [Map/Filter/Reduce with Generics](./map-filter-reduce-generics/) - Higher-order functions for transformations (Go 1.18+)
- [Pipeline Pattern](./pipeline-pattern/) - Channel-based stream processing with fan-out/fan-in
- [Function Composition Pattern](./function-composition-pattern/) - Combine functions to build complex operations
- [Table-Driven Design Pattern](./table-driven-design-pattern/) - Data structures drive logic (idiomatic Go)
- [Generator Pattern](./generator-pattern/) - Lazy evaluation with infinite sequences

### Testing Patterns

- [Mocking Interfaces Pattern](./mocking-interfaces-pattern/) - Test doubles with dependency injection
- [Table Subtests Pattern](./table-subtests-pattern/) - Data-driven tests with `t.Run()`
- [Golden Files Pattern](./golden-files-pattern/) - Snapshot testing for complex outputs

### Advanced Channels

- [Or-Done Channel Pattern](./or-done-channel-pattern/) - Wrap channels with cancellation signal
- [Tee Channel Pattern](./tee-channel-pattern/) - Duplicate channel values to multiple consumers
- [Bridge Channel Pattern](./bridge-channel-pattern/) - Flatten channel of channels

### Performance

- [sync.Pool Pattern](./sync-pool-pattern/) - Object reuse to reduce GC pressure
- [String Builder Pattern](./string-builder-pattern/) - Efficient string concatenation with `strings.Builder`

## Structure

Each pattern lives in its own directory containing:

```
pattern-name/
├── pattern.md    # Explanation, examples, gotchas
└── example.go    # Runnable demonstration
```

Browse patterns:

```bash
cd go-jutsu/
ls -d */                    # List all patterns
cd worker-pool-pattern/     # Enter a pattern
cat pattern.md              # Read explanation
go run example.go           # Run example
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

```

See [`.github/copilot-instructions.md`](./.github/copilot-instructions.md) for agent guide on adding patterns.

## Contributing

Found a cool pattern? See [CONTRIBUTING.md](./CONTRIBUTING.md). PRs welcome!
```
