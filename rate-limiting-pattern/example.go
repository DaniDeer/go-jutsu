package main

import (
	"context"
	"fmt"
	"time"
)

// Simple rate limiter using time.Ticker
type TickerLimiter struct {
	ticker *time.Ticker
}

func NewTickerLimiter(rps int) *TickerLimiter {
	interval := time.Second / time.Duration(rps)
	return &TickerLimiter{
		ticker: time.NewTicker(interval),
	}
}

func (l *TickerLimiter) Wait() {
	<-l.ticker.C
}

func (l *TickerLimiter) Stop() {
	l.ticker.Stop()
}

// Token bucket limiter
type TokenBucket struct {
	tokens   chan struct{}
	rate     time.Duration
	capacity int
}

func NewTokenBucket(rps, burst int) *TokenBucket {
	tb := &TokenBucket{
		tokens:   make(chan struct{}, burst),
		rate:     time.Second / time.Duration(rps),
		capacity: burst,
	}

	// Fill bucket
	for i := 0; i < burst; i++ {
		tb.tokens <- struct{}{}
	}

	// Refill tokens
	go func() {
		ticker := time.NewTicker(tb.rate)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case tb.tokens <- struct{}{}:
			default:
				// Bucket full
			}
		}
	}()

	return tb
}

func (tb *TokenBucket) Allow() bool {
	select {
	case <-tb.tokens:
		return true
	default:
		return false
	}
}

func (tb *TokenBucket) Wait(ctx context.Context) error {
	select {
	case <-tb.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	// Example 1: Simple ticker-based limiting
	fmt.Println("=== Example 1: Ticker Rate Limiter (5 req/sec) ===")
	limiter := NewTickerLimiter(5)
	defer limiter.Stop()

	start := time.Now()
	for i := 0; i < 8; i++ {
		limiter.Wait()
		fmt.Printf("Request %d at %v\n", i+1, time.Since(start).Round(time.Millisecond))
	}

	// Example 2: Token bucket with burst
	fmt.Println("\n=== Example 2: Token Bucket (2/sec, burst 5) ===")
	bucket := NewTokenBucket(2, 5)

	// Burst requests (should allow first 5 immediately)
	start = time.Now()
	fmt.Println("Burst requests:")
	for i := 0; i < 7; i++ {
		if bucket.Allow() {
			fmt.Printf("  Request %d allowed at %v\n", i+1, time.Since(start).Round(time.Millisecond))
		} else {
			fmt.Printf("  Request %d BLOCKED at %v\n", i+1, time.Since(start).Round(time.Millisecond))
		}
	}

	// Example 3: Wait for token
	fmt.Println("\n=== Example 3: Wait for Token ===")
	ctx := context.Background()
	start = time.Now()
	for i := 0; i < 3; i++ {
		bucket.Wait(ctx)
		fmt.Printf("Request %d at %v\n", i+1, time.Since(start).Round(10*time.Millisecond))
	}

	// Example 4: With context cancellation
	fmt.Println("\n=== Example 4: Context Cancellation ===")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := bucket.Wait(ctx)
	if err != nil {
		fmt.Printf("✓ Wait cancelled: %v\n", err)
	}
}
