package main

import (
	"fmt"
	"time"
)

// Example 1: HTTP Request Builder
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

type HTTPRequestBuilder struct {
	req HTTPRequest
}

func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		req: HTTPRequest{
			Method:  "GET",
			Headers: make(map[string]string),
			Timeout: 30 * time.Second,
		},
	}
}

func (b *HTTPRequestBuilder) Method(method string) *HTTPRequestBuilder {
	b.req.Method = method
	return b
}

func (b *HTTPRequestBuilder) URL(url string) *HTTPRequestBuilder {
	b.req.URL = url
	return b
}

func (b *HTTPRequestBuilder) Header(key, value string) *HTTPRequestBuilder {
	b.req.Headers[key] = value
	return b
}

func (b *HTTPRequestBuilder) Body(body []byte) *HTTPRequestBuilder {
	b.req.Body = body
	return b
}

func (b *HTTPRequestBuilder) Timeout(timeout time.Duration) *HTTPRequestBuilder {
	b.req.Timeout = timeout
	return b
}

func (b *HTTPRequestBuilder) Build() HTTPRequest {
	return b.req // Return copy
}

// Example 2: Server Configuration Builder
type Server struct {
	host         string
	port         int
	readTimeout  time.Duration
	writeTimeout time.Duration
	maxConns     int
}

type ServerBuilder struct {
	server Server
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		server: Server{
			host:         "localhost",
			port:         8080,
			readTimeout:  10 * time.Second,
			writeTimeout: 10 * time.Second,
			maxConns:     100,
		},
	}
}

func (b *ServerBuilder) Host(host string) *ServerBuilder {
	b.server.host = host
	return b
}

func (b *ServerBuilder) Port(port int) *ServerBuilder {
	b.server.port = port
	return b
}

func (b *ServerBuilder) Timeouts(read, write time.Duration) *ServerBuilder {
	b.server.readTimeout = read
	b.server.writeTimeout = write
	return b
}

func (b *ServerBuilder) MaxConnections(max int) *ServerBuilder {
	b.server.maxConns = max
	return b
}

func (b *ServerBuilder) Build() (Server, error) {
	// Validation
	if b.server.port < 1 || b.server.port > 65535 {
		return Server{}, fmt.Errorf("invalid port: %d", b.server.port)
	}
	return b.server, nil
}

func main() {
	// Example 1: HTTP Request Builder
	fmt.Println("=== Example 1: HTTP Request Builder ===")
	req := NewHTTPRequestBuilder().
		Method("POST").
		URL("https://api.example.com/users").
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer token123").
		Body([]byte(`{"name":"Alice"}`)).
		Timeout(5 * time.Second).
		Build()

	fmt.Printf("Request: %s %s\n", req.Method, req.URL)
	fmt.Printf("Headers: %v\n", req.Headers)
	fmt.Printf("Timeout: %v\n", req.Timeout)

	// Example 2: Server Builder with defaults
	fmt.Println("\n=== Example 2: Server with Defaults ===")
	server1, _ := NewServerBuilder().Build()
	fmt.Printf("Default: %s:%d (max: %d conns)\n",
		server1.host, server1.port, server1.maxConns)

	// Example 3: Custom server configuration
	fmt.Println("\n=== Example 3: Custom Server ===")
	server2, err := NewServerBuilder().
		Host("0.0.0.0").
		Port(9000).
		Timeouts(30*time.Second, 30*time.Second).
		MaxConnections(500).
		Build()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Custom: %s:%d\n", server2.host, server2.port)
		fmt.Printf("Timeouts: read=%v, write=%v\n",
			server2.readTimeout, server2.writeTimeout)
	}

	// Example 4: Validation
	fmt.Println("\n=== Example 4: Validation ===")
	_, err = NewServerBuilder().Port(99999).Build()
	if err != nil {
		fmt.Printf("✓ Validation caught error: %v\n", err)
	}

	// Example 5: Reusable builder (antipattern)
	fmt.Println("\n=== Example 5: Builder Reuse (Careful!) ===")
	builder := NewHTTPRequestBuilder().
		Method("POST").
		URL("/api/data")

	// Each Build() call returns copy, safe to modify builder
	req1 := builder.Header("X-Version", "1").Build()
	req2 := builder.Header("X-Version", "2").Build()

	fmt.Printf("Req1 version: %s\n", req1.Headers["X-Version"])
	fmt.Printf("Req2 version: %s\n", req2.Headers["X-Version"])
}
