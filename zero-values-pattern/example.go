package main

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

// Example 1: Basic zero values
func demonstrateBasicZeroValues() {
	fmt.Println("1. Basic Zero Values:")

	var i int
	var f float64
	var b bool
	var s string
	var ptr *int

	fmt.Printf("   int:     %d\n", i)
	fmt.Printf("   float64: %g\n", f)
	fmt.Printf("   bool:    %v\n", b)
	fmt.Printf("   string:  %q\n", s)
	fmt.Printf("   pointer: %v\n", ptr)
	fmt.Println()
}

// Example 2: sync.Mutex - zero value ready
type Counter struct {
	mu    sync.Mutex // No initialization needed!
	count int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *Counter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func demonstrateMutexZeroValue() {
	fmt.Println("2. sync.Mutex Zero Value (ready to use):")

	var counter Counter // No New() or initialization!

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}
	wg.Wait()

	fmt.Printf("   Counter (after 100 goroutines): %d\n", counter.Get())
	fmt.Println("   ✓ Mutex worked without initialization!")
	fmt.Println()
}

// Example 3: bytes.Buffer - zero value ready
func demonstrateBufferZeroValue() {
	fmt.Println("3. bytes.Buffer Zero Value:")

	var buf bytes.Buffer // No New()!

	buf.WriteString("Hello ")
	buf.WriteString("World")

	fmt.Printf("   Buffer contents: %q\n", buf.String())
	fmt.Println("   ✓ Buffer worked without initialization!")
	fmt.Println()
}

// Example 4: strings.Builder - zero value ready
func demonstrateBuilderZeroValue() {
	fmt.Println("4. strings.Builder Zero Value:")

	var builder strings.Builder // No New()!

	builder.WriteString("Go ")
	builder.WriteString("is ")
	builder.WriteString("awesome")

	fmt.Printf("   Built string: %q\n", builder.String())
	fmt.Println("   ✓ Builder worked without initialization!")
	fmt.Println()
}

// Example 5: Nil slices can append
func demonstrateNilSliceAppend() {
	fmt.Println("5. Nil Slice Can Append:")

	var slice []int // nil slice
	fmt.Printf("   Initial: len=%d, cap=%d, nil=%v\n", len(slice), cap(slice), slice == nil)

	slice = append(slice, 1, 2, 3)
	fmt.Printf("   After append: len=%d, cap=%d, nil=%v\n", len(slice), cap(slice), slice == nil)
	fmt.Printf("   Contents: %v\n", slice)
	fmt.Println("   ✓ Append works on nil slice!")
	fmt.Println()
}

// Example 6: Nil slice vs empty slice
func demonstrateNilVsEmptySlice() {
	fmt.Println("6. Nil Slice vs Empty Slice:")

	var nilSlice []int
	emptySlice := []int{}
	madeSlice := make([]int, 0)

	fmt.Printf("   nil slice:   len=%d, nil=%v\n", len(nilSlice), nilSlice == nil)
	fmt.Printf("   []int{}:     len=%d, nil=%v\n", len(emptySlice), emptySlice == nil)
	fmt.Printf("   make([]int): len=%d, nil=%v\n", len(madeSlice), madeSlice == nil)

	// All work with append
	nilSlice = append(nilSlice, 1)
	emptySlice = append(emptySlice, 1)
	madeSlice = append(madeSlice, 1)

	fmt.Println("   ✓ All work with append!")
	fmt.Println()
}

// Example 7: Nil map gotcha
func demonstrateNilMapGotcha() {
	fmt.Println("7. Nil Map Gotcha:")

	var m map[string]int // nil map
	fmt.Printf("   Map: len=%d, nil=%v\n", len(m), m == nil)

	// Reading works
	val := m["key"]
	fmt.Printf("   Reading from nil map: %d (zero value)\n", val)

	// Writing panics!
	fmt.Println("   ✗ Cannot write to nil map (would panic)")
	// m["key"] = 42  // panic!

	// Must make() first
	m = make(map[string]int)
	m["key"] = 42
	fmt.Printf("   After make(): %v\n", m)
	fmt.Println()
}

// Example 8: Useful zero value struct
type Config struct {
	Port    int    // 0 = auto-assign
	Timeout int    // 0 = no timeout
	Debug   bool   // false = production
	Name    string // "" = default name
}

func (c *Config) GetPort() int {
	if c.Port == 0 {
		return 8080 // Default
	}
	return c.Port
}

func (c *Config) GetName() string {
	if c.Name == "" {
		return "MyApp" // Default
	}
	return c.Name
}

func demonstrateUsefulZeroStruct() {
	fmt.Println("8. Struct with Useful Zero Values:")

	var cfg Config // Zero values
	fmt.Printf("   Port: %d → effective: %d\n", cfg.Port, cfg.GetPort())
	fmt.Printf("   Timeout: %d (no timeout)\n", cfg.Timeout)
	fmt.Printf("   Debug: %v (production)\n", cfg.Debug)
	fmt.Printf("   Name: %q → effective: %q\n", cfg.Name, cfg.GetName())
	fmt.Println("   ✓ Zero values provide sensible defaults!")
	fmt.Println()
}

// Example 9: sync.WaitGroup - zero value ready
func demonstrateWaitGroupZeroValue() {
	fmt.Println("9. sync.WaitGroup Zero Value:")

	var wg sync.WaitGroup // No New()!

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("   Worker %d done\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("   ✓ WaitGroup worked without initialization!")
	fmt.Println()
}

// Example 10: Type with intentional zero value
type Handler struct {
	buf bytes.Buffer // Zero value works
	mu  sync.Mutex   // Zero value works
}

func (h *Handler) Handle(data string) string {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.buf.WriteString(data)
	return h.buf.String()
}

func demonstrateCompositeZeroValue() {
	fmt.Println("10. Composite Type with Zero Values:")

	var handler Handler // Everything ready!

	result1 := handler.Handle("First ")
	result2 := handler.Handle("Second")

	fmt.Printf("   Result 1: %q\n", result1)
	fmt.Printf("   Result 2: %q\n", result2)
	fmt.Println("   ✓ No initialization needed for entire struct!")
	fmt.Println()
}

func main() {
	fmt.Println("=== Zero Values Are Useful ===\n")

	demonstrateBasicZeroValues()
	demonstrateMutexZeroValue()
	demonstrateBufferZeroValue()
	demonstrateBuilderZeroValue()
	demonstrateNilSliceAppend()
	demonstrateNilVsEmptySlice()
	demonstrateNilMapGotcha()
	demonstrateUsefulZeroStruct()
	demonstrateWaitGroupZeroValue()
	demonstrateCompositeZeroValue()

	fmt.Println("Key Takeaway:")
	fmt.Println("Go initializes everything to zero value")
	fmt.Println("Zero values are designed to be useful:")
	fmt.Println("  ✓ sync.Mutex, sync.WaitGroup - ready to use")
	fmt.Println("  ✓ bytes.Buffer, strings.Builder - ready to use")
	fmt.Println("  ✓ nil slices - can append")
	fmt.Println("  ✗ nil maps - must make() first")

	fmt.Println("\n=== Done ===")
}
