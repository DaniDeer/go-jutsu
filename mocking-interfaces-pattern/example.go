package main

import "fmt"

// Interfaces for dependencies
type EmailSender interface {
	Send(to, subject, body string) error
}

type UserStore interface {
	Get(id int) (*User, error)
	Save(user *User) error
}

// Domain types
type User struct {
	ID    int
	Email string
	Name  string
}

// Service using interfaces
type UserService struct {
	store UserStore
	email EmailSender
}

func (s *UserService) Register(user *User) error {
	if err := s.store.Save(user); err != nil {
		return err
	}
	return s.email.Send(user.Email, "Welcome!", "Thanks for registering")
}

// Mock implementations
type MockUserStore struct {
	users map[int]*User
	saved []*User
}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		users: make(map[int]*User),
	}
}

func (m *MockUserStore) Get(id int) (*User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (m *MockUserStore) Save(user *User) error {
	m.saved = append(m.saved, user)
	m.users[user.ID] = user
	return nil
}

type MockEmailSender struct {
	sent []Email
}

type Email struct {
	To      string
	Subject string
	Body    string
}

func (m *MockEmailSender) Send(to, subject, body string) error {
	m.sent = append(m.sent, Email{to, subject, body})
	fmt.Printf("  [Mock] Email sent to %s: %s\\n", to, subject)
	return nil
}

func main() {
	fmt.Println("=== Mocking Interfaces Pattern ===")

	// Create mocks
	mockStore := NewMockUserStore()
	mockEmail := &MockEmailSender{}

	// Create service with mocks
	svc := &UserService{
		store: mockStore,
		email: mockEmail,
	}

	// Test registration
	user := &User{
		ID:    1,
		Email: "alice@example.com",
		Name:  "Alice",
	}

	fmt.Println("\\nTesting user registration:")
	err := svc.Register(user)
	if err != nil {
		fmt.Printf("Error: %v\\n", err)
		return
	}

	// Verify mock behavior
	fmt.Println("\\nVerifying mocks:")
	if len(mockStore.saved) == 1 {
		fmt.Printf("✓ User saved: %s\\n", mockStore.saved[0].Name)
	}

	if len(mockEmail.sent) == 1 {
		fmt.Printf("✓ Email sent: %v\\n", mockEmail.sent[0])
	}

	// Verify retrieval
	fmt.Println("\\nTesting retrieval:")
	retrieved, _ := mockStore.Get(1)
	if retrieved.Name == "Alice" {
		fmt.Printf("✓ Retrieved user: %s\\n", retrieved.Name)
	}
}
