package main

import (
	"fmt"
	"io"
	"strings"
)

// Example 1: Custom interfaces
type Logger interface {
	Log(msg string)
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Log(msg string) {
	fmt.Println("[LOG]", msg)
}

type UpperLogger struct{}

func (u UpperLogger) Log(msg string) {
	fmt.Println("[LOG]", strings.ToUpper(msg))
}

func useLogger(l Logger, msg string) {
	l.Log(msg)
}

func demonstrateImplicit() {
	fmt.Println("1. Implicit Interface Satisfaction:")

	// Both satisfy Logger without declaring it
	useLogger(ConsoleLogger{}, "hello from console")
	useLogger(UpperLogger{}, "hello from upper")
	fmt.Println()
}

// Example 2: Compile-time verification
type Writer interface {
	Write(data string) error
}

type MyWriter struct{}

func (m *MyWriter) Write(data string) error {
	fmt.Println("  ", data)
	return nil
}

// Compile-time check
var _ Writer = (*MyWriter)(nil)

func demonstrateCompileCheck() {
	fmt.Println("2. Compile-Time Verification:")
	fmt.Println("   var _ Writer = (*MyWriter)(nil)")
	fmt.Println("   ✓ Verified at compile time")
	fmt.Println()
}

// Example 3: io.Reader satisfaction
type UppercaseReader struct {
	r io.Reader
}

func (u *UppercaseReader) Read(p []byte) (int, error) {
	n, err := u.r.Read(p)
	for i := 0; i < n; i++ {
		if p[i] >= 'a' && p[i] <= 'z' {
			p[i] = p[i] - 32
		}
	}
	return n, err
}

func demonstrateStdLibInterface() {
	fmt.Println("3. Satisfying stdlib Interface:")

	// Our type satisfies io.Reader
	r := &UppercaseReader{
		r: strings.NewReader("hello world"),
	}

	data, _ := io.ReadAll(r)
	fmt.Printf("   Read: %s\n", data)
	fmt.Println("   ✓ UppercaseReader is an io.Reader")
	fmt.Println()
}

// Example 4: Interface composition
type Closer interface {
	Close() error
}

type ReadCloser interface {
	io.Reader
	Closer
}

type MyReadCloser struct {
	data string
	pos  int
}

func (m *MyReadCloser) Read(p []byte) (int, error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	n := copy(p, m.data[m.pos:])
	m.pos += n
	return n, nil
}

func (m *MyReadCloser) Close() error {
	fmt.Println("   Closed")
	return nil
}

func demonstrateComposition() {
	fmt.Println("4. Interface Composition:")

	rc := &MyReadCloser{data: "test"}

	// Satisfies ReadCloser
	var _ ReadCloser = rc

	data, _ := io.ReadAll(rc)
	fmt.Printf("   Read: %s\n", data)
	rc.Close()
	fmt.Println("   ✓ Satisfies composed interface")
	fmt.Println()
}

// Example 5: Empty interface
func acceptAnything(v interface{}) {
	fmt.Printf("   Got: %v (type: %T)\n", v, v)
}

func demonstrateEmptyInterface() {
	fmt.Println("5. Empty Interface (any type):")

	acceptAnything(42)
	acceptAnything("hello")
	acceptAnything([]int{1, 2, 3})
	fmt.Println("   ✓ interface{} satisfied by all types")
	fmt.Println()
}

func main() {
	fmt.Println("=== Interface Satisfaction Pattern ===\n")

	demonstrateImplicit()
	demonstrateCompileCheck()
	demonstrateStdLibInterface()
	demonstrateComposition()
	demonstrateEmptyInterface()

	fmt.Println("Key Takeaway:")
	fmt.Println("Interfaces are satisfied implicitly")
	fmt.Println("- No 'implements' keyword")
	fmt.Println("- Duck typing at compile time")
	fmt.Println("- Small interfaces are better")
	fmt.Println("- Accept interfaces, return structs")

	fmt.Println("\n=== Done ===")
}
