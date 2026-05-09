package main

import (
	"fmt"
	"time"
)

// Server config with functional options
type Server struct {
	addr     string
	timeout  time.Duration
	maxConns int
	debug    bool
}

// Option is a function that configures Server
type Option func(*Server)

// Option functions
func WithTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.timeout = d
	}
}

func WithMaxConns(n int) Option {
	return func(s *Server) {
		s.maxConns = n
	}
}

func WithDebug(debug bool) Option {
	return func(s *Server) {
		s.debug = debug
	}
}

// Constructor
func NewServer(addr string, opts ...Option) *Server {
	// Defaults
	s := &Server{
		addr:     addr,
		timeout:  30 * time.Second,
		maxConns: 100,
		debug:    false,
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Show() {
	fmt.Printf("   Server{addr:%s, timeout:%v, maxConns:%d, debug:%v}\n",
		s.addr, s.timeout, s.maxConns, s.debug)
}

func main() {
	fmt.Println("=== Functional Options Pattern ===\n")

	fmt.Println("1. Default configuration:")
	s1 := NewServer(":8080")
	s1.Show()

	fmt.Println("\n2. With timeout option:")
	s2 := NewServer(":8080",
		WithTimeout(60*time.Second),
	)
	s2.Show()

	fmt.Println("\n3. Multiple options:")
	s3 := NewServer(":8080",
		WithTimeout(60*time.Second),
		WithMaxConns(200),
		WithDebug(true),
	)
	s3.Show()

	fmt.Println("\n4. Options can be variables:")
	debugOpt := WithDebug(true)
	perfOpts := []Option{
		WithTimeout(120 * time.Second),
		WithMaxConns(500),
	}

	s4 := NewServer(":8080", append(perfOpts, debugOpt)...)
	s4.Show()

	fmt.Println("\n5. Last option wins on conflicts:")
	s5 := NewServer(":8080",
		WithTimeout(10*time.Second),
		WithTimeout(20*time.Second), // Overwrites
	)
	s5.Show()

	fmt.Println("\nKey Takeaway:")
	fmt.Println("Functional Options pattern provides:")
	fmt.Println("- Clean API for optional config")
	fmt.Println("- Backward compatibility")
	fmt.Println("- Self-documenting option names")
	fmt.Println("- No config struct pollution")

	fmt.Println("\n=== Done ===")
}
