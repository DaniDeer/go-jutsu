# Examples

Runnable Go code demonstrating patterns in action.

Each example is standalone and can be run with `go run <filename>.go`

## Current Examples

All examples runnable with `go run <filename>.go`

### Foundational

- `go-closure-currying-pattern.go` - HTTP middleware using closures
- `defer-cleanup-pattern.go` - LIFO cleanup, gotchas
- `empty-struct-signal-pattern.go` - Zero-byte signaling
- `zero-values-pattern.go` - Useful defaults by design

### Gotchas & Idioms

- `nil-interface-gotcha.go` - Interface (type, value) trap
- `comma-ok-idiom.go` - Safe map/channel/type access
- `blank-identifier-pattern.go` - Using \_ effectively

### Concurrency

- `select-statement-pattern.go` - Channel multiplexing
- `context-propagation-pattern.go` - Cancellation signals

### Advanced

- `functional-options-pattern.go` - Configuration pattern
- `interface-satisfaction-pattern.go` - Implicit interfaces
- `embedding-composition-pattern.go` - Composition not inheritance
- `method-values-expressions-pattern.go` - Methods as values
- `type-switch-pattern.go` - Type-specific handling

### Functional & Declarative

- `map-filter-reduce-generics.go` - Higher-order functions with generics
- `pipeline-pattern.go` - Channel-based stream processing
- `function-composition-pattern.go` - Compose functions
- `table-driven-design-pattern.go` - Data-driven design
- `generator-pattern.go` - Lazy evaluation, infinite sequences
