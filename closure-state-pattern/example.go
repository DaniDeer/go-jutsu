package main

// Closure state: the returned function captures and mutates variables
// in the outer scope — each closure gets its own independent state.

import (
	"fmt"
	"net/http"
	"sync"
)

// idGenerator returns a function that produces unique prefixed IDs.
// Each call to idGenerator creates an independent counter.
func idGenerator(prefix string) func() string {
	n := 0 // captured state — private to this closure
	return func() string {
		n++
		return fmt.Sprintf("%s-%04d", prefix, n)
	}
}

// requestCounter returns middleware that tracks the total request count.
// The count lives inside the closure — no global variables needed.
func requestCounter() func(http.Handler) http.Handler {
	var mu sync.Mutex
	count := 0
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			count++
			n := count
			mu.Unlock()
			fmt.Printf("[counter] request #%d: %s %s\n", n, r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

// newCache returns a (set, get) pair backed by a shared map.
// Both functions close over the same store — one writes, one reads.
func newCache() (set func(key, val string), get func(key string) (string, bool)) {
	var mu sync.Mutex
	store := make(map[string]string)

	set = func(key, val string) {
		mu.Lock()
		store[key] = val
		mu.Unlock()
	}
	get = func(key string) (string, bool) {
		mu.Lock()
		defer mu.Unlock()
		v, ok := store[key]
		return v, ok
	}
	return
}

func main() {
	// --- ID generators: independent counters ---
	newUserID := idGenerator("user")
	newOrderID := idGenerator("order")

	fmt.Println(newUserID())  // user-0001
	fmt.Println(newUserID())  // user-0002
	fmt.Println(newOrderID()) // order-0001 — own counter, starts at 1
	fmt.Println(newUserID())  // user-0003 — resumes from 2

	fmt.Println()

	// --- In-memory cache ---
	set, get := newCache()
	set("token", "abc123")
	set("user", "alice")

	if v, ok := get("token"); ok {
		fmt.Println("cached token:", v) // cached token: abc123
	}
	if _, ok := get("missing"); !ok {
		fmt.Println("key not found") // key not found
	}

	fmt.Println()

	// --- Request counter middleware ---
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	// Wrap once — the closure tracks state across all requests
	handler := requestCounter()(mux)

	fmt.Println("handler ready:", handler != nil)
	fmt.Println("(start server to see request counter in action)")
}
