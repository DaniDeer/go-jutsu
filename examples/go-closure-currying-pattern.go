package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// requestLogger creates a middleware that logs HTTP requests
func requestLogger(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

// rateLimit creates a middleware that limits requests per minute
func rateLimit(reqPerMin int) func(http.Handler) http.Handler {
	ticker := time.NewTicker(time.Minute / time.Duration(reqPerMin))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-ticker.C
			next.ServeHTTP(w, r)
		})
	}
}

// chain combines multiple middleware into one
func chain(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from go-jutsu!\n")
}

func main() {
	logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	// Example 1: Single middleware
	handler1 := requestLogger(logger)(mux)

	// Example 2: Multiple middleware chained
	handler2 := chain(mux,
		requestLogger(logger),
		rateLimit(60), // 60 requests per minute
	)

	_ = handler1
	_ = handler2

	fmt.Println("Server starting on :8080")
	fmt.Println("Try: curl http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", handler2))
}
