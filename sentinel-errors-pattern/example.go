package main

import (
	"errors"
	"fmt"
)

// Sentinel errors
var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("already exists")
	ErrUnauthorized  = errors.New("unauthorized")
)

// User type
type User struct {
	ID    int
	Email string
}

// Mock database
var users = map[int]*User{
	1: {ID: 1, Email: "alice@example.com"},
	2: {ID: 2, Email: "bob@example.com"},
}

func GetUser(id int) (*User, error) {
	user, ok := users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func CreateUser(email string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidInput
	}

	// Check duplicate
	for _, u := range users {
		if u.Email == email {
			return nil, ErrAlreadyExists
		}
	}

	user := &User{ID: len(users) + 1, Email: email}
	users[user.ID] = user
	return user, nil
}

func DeleteUser(id int) error {
	if _, ok := users[id]; !ok {
		return ErrNotFound
	}
	delete(users, id)
	return nil
}

// Wrapped error example
func GetUserWithContext(id int) (*User, error) {
	user, err := GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("GetUserWithContext(id=%d): %w", id, err)
	}
	return user, nil
}

func main() {
	fmt.Println("=== Sentinel Errors Pattern ===\n")

	// Example 1: errors.Is for checking
	fmt.Println("1. Check with errors.Is:")
	_, err := GetUser(999)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("   ✓ User not found (expected)")
	}
	fmt.Println()

	// Example 2: Different sentinel errors
	fmt.Println("2. Multiple Sentinel Errors:")

	_, err = CreateUser("")
	if errors.Is(err, ErrInvalidInput) {
		fmt.Println("   ✓ Invalid input")
	}

	_, err = CreateUser("alice@example.com")
	if errors.Is(err, ErrAlreadyExists) {
		fmt.Println("   ✓ Already exists")
	}
	fmt.Println()

	// Example 3: errors.Is with wrapped errors
	fmt.Println("3. Wrapped Errors Still Match:")
	_, err = GetUserWithContext(999)
	fmt.Printf("   Error: %v\n", err)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("   ✓ Still matches ErrNotFound")
	}
	fmt.Println()

	// Example 4: Successful operations
	fmt.Println("4. Successful Operations:")
	user, err := GetUser(1)
	if err == nil {
		fmt.Printf("   ✓ Found user: %s\n", user.Email)
	}

	newUser, err := CreateUser("charlie@example.com")
	if err == nil {
		fmt.Printf("   ✓ Created user: %s (ID: %d)\n", newUser.Email, newUser.ID)
	}
	fmt.Println()

	fmt.Println("Key Takeaway:")
	fmt.Println("Sentinel errors:")
	fmt.Println("- Package-level error variables")
	fmt.Println("- Use errors.Is() for checking")
	fmt.Println("- Works with error wrapping")
	fmt.Println("- Self-documenting error conditions")

	fmt.Println("\n=== Done ===")
}
