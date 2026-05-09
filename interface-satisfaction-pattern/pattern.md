# Interface Satisfaction Pattern

Implicit interface implementation. No "implements" keyword. Duck typing at compile time.

## What It Is (and Isn't)

Types satisfy interfaces automatically if they have matching methods. Compiler checks.

Not inheritance. Not explicit declaration. Structural typing unique to Go.

## Where You See It

**Standard library:**

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Any type with Read() is a Reader
type MyType struct{}
func (m MyType) Read(p []byte) (int, error) { return 0, nil }
// MyType is now a Reader!
```

**Small interfaces:**

```go
type Stringer interface {
    String() string
}
// Implement one method, satisfy interface
```

**Compile-time check:**

```go
// Verify *MyType implements io.Writer
var _ io.Writer = (*MyType)(nil)
```

## Real Example

```go
type Logger interface {
    Log(msg string)
}

// FileLogger implements Logger (implicitly)
type FileLogger struct {
    file *os.File
}

func (f *FileLogger) Log(msg string) {
    f.file.WriteString(msg + "\n")
}

// ConsoleLogger also implements Logger
type ConsoleLogger struct{}

func (c *ConsoleLogger) Log(msg string) {
    fmt.Println(msg)
}

// Accept interface, return struct
func NewService(logger Logger) *Service {
    return &Service{logger: logger}
}

// Works with anyLogger
service1 := NewService(&FileLogger{file})
service2 := NewService(&ConsoleLogger{})
```

## Gotchas

**Pointer vs value receivers:**

```go
type Writer interface {
    Write([]byte) error
}

type MyWriter struct{}

func (m *MyWriter) Write(b []byte) error { return nil }

// Only *MyWriter satisfies Writer, not MyWriter
var _ Writer = (*MyWriter)(nil)  // ✓
var _ Writer = MyWriter{}         // ✗ compile error
```

**Accept interfaces, return structs:**

```go
// GOOD
func NewThing(r io.Reader) *Thing

// BAD (too restrictive)
func NewThing(f *os.File) *Thing
```

**Interface segregation:**

```go
// GOOD: small, focused interfaces
type Reader interface { Read([]byte) (int, error) }
type Writer interface { Write([]byte) (int, error) }
type ReadWriter interface { Reader; Writer }

// BAD: kitchen sink interface
type Everything interface {
    Read()
    Write()
    Close()
    Seek()
    // ...
}
```

## Compile-Time Verification

```go
// Ensure type satisfies interface
var _ io.Reader = (*MyType)(nil)
var _ fmt.Stringer = MyType{}

// Common in packages
var (
    _ Handler = (*MyHandler)(nil)
    _ Writer  = (*MyWriter)(nil)
)
```

## Benefits

- Minimal coupling
- Easy testing (mock implementations)
- Compose interfaces from smaller ones
- Types can satisfy interfaces without importing them
