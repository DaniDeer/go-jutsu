package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state        State
	failures     int
	threshold    int
	timeout      time.Duration
	lastFailTime time.Time
	mu           sync.Mutex
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:     Closed,
		threshold: threshold,
		timeout:   timeout,
	}
}

var ErrCircuitOpen = errors.New("circuit breaker open")

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Check if we should transition from Open to HalfOpen
	if cb.state == Open {
		if time.Since(cb.lastFailTime) > cb.timeout {
			cb.state = HalfOpen
			fmt.Println("  Circuit: Open -> HalfOpen")
		} else {
			return ErrCircuitOpen
		}
	}

	// Execute function
	err := fn()
	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()
		if cb.failures >= cb.threshold {
			cb.state = Open
			fmt.Printf("  Circuit: -> Open (failures: %d)\n", cb.failures)
		}
		return err
	}

	// Success: reset
	if cb.state == HalfOpen {
		fmt.Println("  Circuit: HalfOpen -> Closed")
	}
	cb.state = Closed
	cb.failures = 0
	return nil
}

// Simulated failing service
var callCount int

func flakeyService() error {
	callCount++
	if callCount <= 5 {
		return fmt.Errorf("service unavailable")
	}
	return nil
}

func main() {
	fmt.Println("=== Circuit Breaker Pattern ===")
	cb := NewCircuitBreaker(3, 2*time.Second)

	// Simulate repeated failures
	for i := 1; i <= 8; i++ {
		err := cb.Call(flakeyService)
		if err != nil {
			if errors.Is(err, ErrCircuitOpen) {
				fmt.Printf("Call %d: Blocked by circuit breaker\n", i)
			} else {
				fmt.Printf("Call %d: Failed - %v\n", i, err)
			}
		} else {
			fmt.Printf("Call %d: Success\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Wait for timeout, then one more call
	fmt.Println("\nWaiting for circuit timeout...")
	time.Sleep(2 * time.Second)

	err := cb.Call(func() error { return nil })
	if err == nil {
		fmt.Println("Call after timeout: Success (circuit recovered)")
	}
}
