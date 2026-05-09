package main

import (
	"fmt"
	"sync"
)

// Example 1: Basic embedding
type Logger struct {
	prefix string
}

func (l *Logger) Log(msg string) {
	fmt.Printf("%s %s\n", l.prefix, msg)
}

type Service struct {
	Logger // Embedded
	name   string
}

func demonstrateBasic() {
	fmt.Println("1. Basic Embedding:")

	s := &Service{
		Logger: Logger{prefix: "[SVC]"},
		name:   "UserService",
	}

	// Method promoted from Logger
	s.Log("service starting")
	fmt.Println()
}

// Example 2: sync.Mutex embedding
type Counter struct {
	sync.Mutex // Embedded
	count      int
}

func (c *Counter) Inc() {
	c.Lock() // Method promoted!
	defer c.Unlock()
	c.count++
}

func demonstrateMutex() {
	fmt.Println("2. Embed sync.Mutex:")

	c := &Counter{}
	c.Inc()
	c.Inc()

	fmt.Printf("   Count: %d\n", c.count)
	fmt.Println("   ✓ Lock/Unlock methods promoted")
	fmt.Println()
}

// Example 3: Interface embedding
type Reader interface {
	Read() string
}

type Writer interface {
	Write(string)
}

type ReadWriter interface {
	Reader // Embedded interface
	Writer
}

type File struct{}

func (f *File) Read() string   { return "data" }
func (f *File) Write(s string) { fmt.Println("  Write:", s) }

func demonstrateInterfaceEmbedding() {
	fmt.Println("3. Interface Embedding:")

	var rw ReadWriter = &File{}
	data := rw.Read()
	rw.Write(data)
	fmt.Println("   ✓ File satisfies ReadWriter")
	fmt.Println()
}

// Example 4: Name collision
type A struct{}

func (a A) Method() { fmt.Println("   A.Method()") }

type B struct{}

func (b B) Method() { fmt.Println("   B.Method()") }

type C struct {
	A
	B
}

func demonstrateCollision() {
	fmt.Println("4. Name Collision:")

	c := C{}
	// c.Method()  // Ambiguous! Won't compile

	// Must be explicit
	c.A.Method()
	c.B.Method()
	fmt.Println("   ⚠ Collision requires explicit access")
	fmt.Println()
}

// Example 5: Not inheritance
type Animal struct {
	name string
}

func (a *Animal) Speak() {
	fmt.Printf("   %s makes a sound\n", a.name)
}

type Dog struct {
	Animal // Composition, not inheritance
}

func acceptAnimal(a Animal) {
	a.Speak()
}

func demonstrateNotInheritance() {
	fmt.Println("5. Embedding is NOT Inheritance:")

	dog := Dog{Animal: Animal{name: "Buddy"}}
	dog.Speak() // Method promoted

	// acceptAnimal(dog)  // ✗ Won't compile - Dog is not Animal
	acceptAnimal(dog.Animal) // ✓ Must pass embedded type

	fmt.Println("   ⚠ Dog is NOT an Animal (composition)")
	fmt.Println()
}

func main() {
	fmt.Println("=== Embedding Composition Pattern ===\n")

	demonstrateBasic()
	demonstrateMutex()
	demonstrateInterfaceEmbedding()
	demonstrateCollision()
	demonstrateNotInheritance()

	fmt.Println("Key Takeaway:")
	fmt.Println("Embedding != Inheritance")
	fmt.Println("- Methods/fields promoted to outer type")
	fmt.Println("- Composition over inheritance")
	fmt.Println("- Name collisions require explicit access")
	fmt.Println("- Common: embed sync.Mutex, io.Reader, etc.")

	fmt.Println("\n=== Done ===")
}
