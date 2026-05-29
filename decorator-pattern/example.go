package main

// Demonstrates two forms of the Decorator Pattern in Go:
//  1. Function decorator — wrap a func to add cross-cutting behavior
//  2. Interface decorator — wrap an io.Reader to track bytes read

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ── 1. Function Decorator ────────────────────────────────────────────────────

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// withLogging decorates any HandlerFunc with request timing.
// Go equivalent of Python's @log_duration decorator.
func withLogging(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("[%s] %s %s  (%v)", "INFO", r.Method, r.URL.Path, time.Since(start))
	}
}

// withAuth decorates any HandlerFunc with a token check.
func withAuth(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// serveAPI is the original handler — unaware of logging or auth.
func serveAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from API")
}

// ── 2. Interface Decorator ───────────────────────────────────────────────────

// countingReader wraps any io.Reader and tracks total bytes consumed.
// Implements io.Reader — drop-in replacement, no callers need to change.
type countingReader struct {
	r     io.Reader
	total int64
}

func NewCountingReader(r io.Reader) *countingReader {
	return &countingReader{r: r}
}

func (c *countingReader) Read(p []byte) (n int, err error) {
	n, err = c.r.Read(p)
	c.total += int64(n)
	return
}

// ── main ─────────────────────────────────────────────────────────────────────

func main() {
	// --- Function decorator demo ---
	fmt.Println("=== Function Decorator ===")

	// Stack decorators: withLogging(withAuth(serveAPI))
	// Outer decorator runs first → logs every request; auth may short-circuit.
	handler := withLogging(withAuth(serveAPI))

	// Simulate two requests using net/http/httptest style inline
	req1, _ := http.NewRequest("GET", "/api/data", nil)
	req1.Header.Set("Authorization", "Bearer token123")

	req2, _ := http.NewRequest("GET", "/api/data", nil) // no token

	fmt.Println("Request with valid token:")
	handler(&noopResponseWriter{}, req1)

	fmt.Println("\nRequest without token:")
	handler(&noopResponseWriter{}, req2)

	// --- Interface decorator demo ---
	fmt.Println("\n=== Interface Decorator ===")

	body := strings.NewReader("Go decorators wrap behavior without touching the original.")
	cr := NewCountingReader(body)

	data, _ := io.ReadAll(cr) // cr satisfies io.Reader transparently
	fmt.Printf("Content : %s\n", data)
	fmt.Printf("Bytes read: %d\n", cr.total)
}

// noopResponseWriter satisfies http.ResponseWriter for the demo.
type noopResponseWriter struct{ code int }

func (n *noopResponseWriter) Header() http.Header         { return http.Header{} }
func (n *noopResponseWriter) Write(b []byte) (int, error) { fmt.Print(string(b)); return len(b), nil }
func (n *noopResponseWriter) WriteHeader(code int)        { n.code = code; fmt.Printf("[HTTP %d]\n", code) }
