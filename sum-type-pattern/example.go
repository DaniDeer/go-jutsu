package main

// Sum type pattern: model a closed set of variants using a sealed interface.
// The unexported method (paymentStatusTag) prevents external types from
// implementing PaymentStatus — only the variants defined here are valid.

import (
	"fmt"
)

// --- Sealed interface ---

type PaymentStatus interface {
	paymentStatusTag() // unexported: seals the type set to this package
}

// --- Variants (all data lives in the concrete type) ---

type Pending struct{}
type Completed struct{ TxID string }
type Failed struct{ Reason string }
type Refunded struct{ Amount float64 }

func (Pending) paymentStatusTag()   {}
func (Completed) paymentStatusTag() {}
func (Failed) paymentStatusTag()    {}
func (Refunded) paymentStatusTag()  {}

// --- Exhaustive dispatch ---

// describe handles every variant. The panic documents that the switch is
// exhaustive — if a new variant is added without updating this, it surfaces fast.
func describe(s PaymentStatus) string {
	switch v := s.(type) {
	case Pending:
		return "awaiting payment"
	case Completed:
		return fmt.Sprintf("paid (tx=%s)", v.TxID)
	case Failed:
		return fmt.Sprintf("failed: %s", v.Reason)
	case Refunded:
		return fmt.Sprintf("refunded $%.2f", v.Amount)
	}
	panic("unreachable: sealed interface guarantees exhaustive switch")
}

// --- Bonus: generic Result[T] sum type ---

// Result holds either a success value or an error — explicit alternative to (T, error).
type Result[T any] struct {
	val T
	err error
	ok  bool
}

func Ok[T any](v T) Result[T]      { return Result[T]{val: v, ok: true} }
func Err[T any](e error) Result[T] { return Result[T]{err: e} }

func (r Result[T]) IsOk() bool         { return r.ok }
func (r Result[T]) Unwrap() (T, error) { return r.val, r.err }

// fetchUser simulates a lookup that may succeed or fail
func fetchUser(id int) Result[string] {
	if id == 42 {
		return Ok("alice")
	}
	return Err[string](fmt.Errorf("user %d not found", id))
}

func main() {
	// --- PaymentStatus sum type ---
	statuses := []PaymentStatus{
		Pending{},
		Completed{TxID: "txn_abc123"},
		Failed{Reason: "insufficient funds"},
		Refunded{Amount: 49.99},
	}

	fmt.Println("=== Payment statuses ===")
	for _, s := range statuses {
		fmt.Println(" ", describe(s))
	}

	fmt.Println()

	// --- Result[T] ---
	fmt.Println("=== Result[T] ===")
	for _, id := range []int{42, 99} {
		r := fetchUser(id)
		if name, err := r.Unwrap(); r.IsOk() {
			fmt.Printf("  user %d → %s\n", id, name)
		} else {
			fmt.Printf("  user %d → error: %v\n", id, err)
		}
	}
}
