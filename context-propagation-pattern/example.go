package main

import (
	"context"
	"fmt"
	"time"
)

// Example 1: Basic cancellation
func demonstrateCancellation() {
	fmt.Println("1. Basic Cancellation:")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("   Worker: stopping")
				return
			case <-time.After(100 * time.Millisecond):
				fmt.Println("   Worker: tick")
			}
		}
	}()

	time.Sleep(350 * time.Millisecond)
	fmt.Println("   Main: canceling")
	cancel()
	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

// Example 2: Timeout
func slowOperation(ctx context.Context) error {
	select {
	case <-time.After(2 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func demonstrateTimeout() {
	fmt.Println("2. Timeout:")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := slowOperation(ctx)
	elapsed := time.Since(start)

	if err == context.DeadlineExceeded {
		fmt.Printf("   ✓ Timed out after %v\n", elapsed)
	} else {
		fmt.Printf("   Completed: %v\n", err)
	}
	fmt.Println()
}

// Example 3: Parent-child cancellation
func demonstrateParentChild() {
	fmt.Println("3. Parent-Child Cancellation:")

	parent, cancelParent := context.WithCancel(context.Background())
	child1, _ := context.WithCancel(parent)
	child2, _ := context.WithTimeout(parent, 10*time.Second)

	go func() {
		<-child1.Done()
		fmt.Println("   Child1 canceled")
	}()

	go func() {
		<-child2.Done()
		fmt.Println("   Child2 canceled")
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("   Canceling parent")
	cancelParent() // Cancels all children

	time.Sleep(100 * time.Millisecond)
	fmt.Println()
}

// Example 4: Context values (use sparingly)
type userKey struct{}

func demonstrateValues() {
	fmt.Println("4. Context Values (use sparingly):")

	ctx := context.WithValue(context.Background(), userKey{}, "alice")

	if user, ok := ctx.Value(userKey{}).(string); ok {
		fmt.Printf("   User from context: %s\n", user)
	}

	fmt.Println("   ⚠ Prefer explicit parameters over context values")
	fmt.Println()
}

// Example 5: Deadline
func demonstrateDeadline() {
	fmt.Println("5. Deadline:")

	deadline := time.Now().Add(1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("   Deadline: %v\n", d.Format("15:04:05"))
		fmt.Printf("   Time until deadline: %v\n", time.Until(d))
	}

	time.Sleep(1100 * time.Millisecond)

	select {
	case <-ctx.Done():
		fmt.Printf("   ✓ Context done: %v\n", ctx.Err())
	default:
		fmt.Println("   Still running")
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== Context Propagation Pattern ===\n")

	demonstrateCancellation()
	demonstrateTimeout()
	demonstrateParentChild()
	demonstrateValues()
	demonstrateDeadline()

	fmt.Println("Key Takeaway:")
	fmt.Println("Context carries cancellation/timeout signals")
	fmt.Println("- First parameter by convention")
	fmt.Println("- Always defer cancel()")
	fmt.Println("- Check ctx.Done() in select")
	fmt.Println("- Avoid context values for business data")

	fmt.Println("\n=== Done ===")
}
