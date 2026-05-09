package main

import (
	"bytes"
	"fmt"
	"sync"
)

// Example 1: Buffer pool (common use case)
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func ProcessData(data string) string {
	// Get buffer from pool
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset() // MUST reset before use
	defer bufferPool.Put(buf)

	// Use buffer
	buf.WriteString("Processed: ")
	buf.WriteString(data)
	return buf.String()
}

// Example 2: Custom struct pool
type Request struct {
	ID   int
	Data []byte
}

func (r *Request) Reset() {
	r.ID = 0
	r.Data = r.Data[:0] // Keep capacity, reset length
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &Request{
			Data: make([]byte, 0, 1024),
		}
	},
}

func HandleRequest(id int, data []byte) {
	req := requestPool.Get().(*Request)
	req.Reset()
	defer requestPool.Put(req)

	req.ID = id
	req.Data = append(req.Data, data...)
	fmt.Printf("Request %d: %d bytes\\n", req.ID, len(req.Data))
}

// Example 3: Benchmark comparison
func concatenateStrings(n int) string {
	var result string
	for i := 0; i < n; i++ {
		result += "x"
	}
	return result
}

func concatenateWithPool(n int) string {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	for i := 0; i < n; i++ {
		buf.WriteString("x")
	}
	return buf.String()
}

func main() {
	// Example 1: Basic buffer pooling
	fmt.Println("=== Example 1: Buffer Pool ===")
	for i := 0; i < 3; i++ {
		result := ProcessData(fmt.Sprintf("item-%d", i))
		fmt.Println(result)
	}

	// Example 2: Custom struct pool
	fmt.Println("\\n=== Example 2: Request Pool ===")
	for i := 0; i < 3; i++ {
		HandleRequest(i, []byte(fmt.Sprintf("data-%d", i)))
	}

	// Example 3: Concurrent usage
	fmt.Println("\\n=== Example 3: Concurrent Pool Usage ===")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			HandleRequest(id, []byte("concurrent-data"))
		}(i)
	}
	wg.Wait()

	// Example 4: Pool statistics (demonstrate pool may be cleared)
	fmt.Println("\\n=== Example 4: Pool Behavior ===")
	buf1 := bufferPool.Get().(*bytes.Buffer)
	buf1.WriteString("marker")
	bufferPool.Put(buf1)

	buf2 := bufferPool.Get().(*bytes.Buffer)
	if buf2.Len() > 0 {
		fmt.Println("Got same buffer from pool (contains: ", buf2.String(), ")")
	} else {
		fmt.Println("Got new buffer (pool was cleared or different buffer)")
	}
	bufferPool.Put(buf2)
}
