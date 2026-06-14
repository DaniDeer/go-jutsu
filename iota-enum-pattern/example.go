package main

// Iota enum pattern: named integer type + const block + String() + Valid().
//
// Key difference from sealed interface (sum type):
//   - All variants share the same shape (just an int, no payload)
//   - NOT truly sealed: OrderStatus(99) compiles — guard with Valid()
//   - Use when variants are data-free labels (flags, status codes, lifecycle states)

import "fmt"

// ── Enum definition ──────────────────────────────────────────────────────────

type OrderStatus int

const (
	Pending   OrderStatus = iota // 0
	Approved                     // 1
	Shipped                      // 2
	Delivered                    // 3
	Cancelled                    // 4

	orderStatusSentinel // unexported sentinel: marks upper bound for Valid()
)

// ── String() ─────────────────────────────────────────────────────────────────

// String makes OrderStatus print as a name instead of a raw integer.
// Without this, fmt.Println(Approved) would print "1".
func (s OrderStatus) String() string {
	names := [...]string{"Pending", "Approved", "Shipped", "Delivered", "Cancelled"}
	if !s.Valid() {
		return fmt.Sprintf("OrderStatus(%d)", int(s))
	}
	return names[s]
}

// ── Valid() guard ─────────────────────────────────────────────────────────────

// Valid returns false for any integer cast that falls outside the defined range.
// This is necessary because Go enums are not truly sealed — any int converts freely.
func (s OrderStatus) Valid() bool {
	return s >= Pending && s < orderStatusSentinel
}

// ── Dispatch (type switch equivalent) ────────────────────────────────────────

// describeTransition explains what happens when entering a state.
// The default panic documents that the switch should be exhaustive —
// if a new variant is added without updating this function, it surfaces fast.
func describeTransition(s OrderStatus) string {
	if !s.Valid() {
		panic(fmt.Sprintf("invalid OrderStatus: %d", int(s)))
	}
	switch s {
	case Pending:
		return "order received, awaiting approval"
	case Approved:
		return "payment confirmed, preparing shipment"
	case Shipped:
		return "handed to carrier"
	case Delivered:
		return "delivered to customer"
	case Cancelled:
		return "order cancelled, refund initiated"
	default:
		panic("unreachable: Valid() guarantees exhaustive switch")
	}
}

// ── Bitmask variant (bonus) ───────────────────────────────────────────────────

// Permission demonstrates 1 << iota for flag/bitmask enums.
// Sequential iota (0,1,2,...) would produce 0,1,2 — not useful for bitmasking.
type Permission int

const (
	Read    Permission = 1 << iota // 1 (0b001)
	Write                          // 2 (0b010)
	Execute                        // 4 (0b100)
)

func (p Permission) String() string {
	var s string
	if p&Read != 0 {
		s += "r"
	} else {
		s += "-"
	}
	if p&Write != 0 {
		s += "w"
	} else {
		s += "-"
	}
	if p&Execute != 0 {
		s += "x"
	} else {
		s += "-"
	}
	return s
}

// ── main ─────────────────────────────────────────────────────────────────────

func main() {
	// --- Lifecycle states ---
	fmt.Println("=== Order lifecycle ===")
	states := []OrderStatus{Pending, Approved, Shipped, Delivered, Cancelled}
	for _, s := range states {
		fmt.Printf("  %-10s → %s\n", s, describeTransition(s))
	}

	fmt.Println()

	// --- The open-integer gotcha ---
	fmt.Println("=== Open-integer gotcha ===")
	ghost := OrderStatus(99)
	fmt.Printf("  OrderStatus(99).Valid()  = %v\n", ghost.Valid())
	fmt.Printf("  OrderStatus(99).String() = %q\n", ghost.String()) // safe: String checks Valid()

	// Contrast: a sealed interface (sum-type-pattern) would make this impossible at compile time.

	fmt.Println()

	// --- Bitmask enum ---
	fmt.Println("=== Permission bitmask ===")
	owner := Read | Write | Execute
	group := Read | Execute
	other := Read

	fmt.Printf("  owner: %s (%d)\n", owner, int(owner))
	fmt.Printf("  group: %s (%d)\n", group, int(group))
	fmt.Printf("  other: %s (%d)\n", other, int(other))
}
