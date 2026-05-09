# Embedding Composition Pattern

Embed types to compose functionality. Not inheritance. Methods promoted automatically.

## What It Is (and Isn't)

Embedding copies outer type's interface. Promotes embedded type's methods.

Not inheritance. No "is-a" relationship. Composition over inheritance realized.

## Where You See It

**Struct embedding:**

```go
type Reader struct {
    io.Reader  // Embedded
}
// Reader.Read() available automatically
```

**Interface embedding:**

```go
type ReadWriter interface {
    Reader
    Writer
}
```

**Anonymous fields:**

```go
type User struct {
    sync.Mutex  // Can call user.Lock() directly
    name string
}
```

## Real Example

```go
type Logger struct {
    prefix string
}

func (l *Logger) Log(msg string) {
    fmt.Println(l.prefix, msg)
}

// Service embeds Logger
type Service struct {
    Logger  // Embedded
    db      *DB
}

// Can call s.Log() directly
s := &Service{
    Logger: Logger{prefix: "[SVC]"},
}
s.Log("starting")  // Uses embedded Logger.Log
```

## Gotchas

**Name collisions:**

```go
type A struct{}
func (a A) Method() { fmt.Println("A") }

type B struct{}
func (b B) Method() { fmt.Println("B") }

type C struct {
    A
    B
}
// c.Method() // Ambiguous! Won't compile
// Must use c.A.Method() or c.B.Method()
```

**Promotion rules:**

```go
type Inner struct {
    X int
}

func (i Inner) Get() int { return i.X }

type Outer struct {
    Inner  // Embedded
    Y int
}

o := Outer{Inner: Inner{X: 1}, Y: 2}
o.Get()  // Promoted from Inner
o.X      // Promoted field
```

**Pointer vs value embedding:**

```go
type A struct{}
func (a *A) Method() {}

type B struct {
    A   // Value embedding - Method NOT promoted
}

type C struct {
    *A  // Pointer embedding - Method promoted
}
```

## Not Inheritance

```go
type Animal struct{}
func (a Animal) Speak() {}

type Dog struct {
    Animal  // Embedded, not inherited
}

// Dog is NOT an Animal
var _ Animal = Dog{}  // Compile error

// But methods are promoted
dog := Dog{}
dog.Speak()  // Works
```
