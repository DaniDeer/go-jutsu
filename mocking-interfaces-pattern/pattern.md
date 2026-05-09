# Mocking Interfaces Pattern

Test doubles for dependencies. Interface-based design enables swapping implementations.

## What It Is (and Isn't)

Fake implementation for testing. Dependency injection. Verify behavior without real dependencies.

Not monkey patching. Not reflection magic. Go way: small interfaces.

## Where You See It

**Database mock:**

```go
type UserStore interface {
    Get(id int) (*User, error)
}

type MockUserStore struct {
    users map[int]*User
}

func (m *MockUserStore) Get(id int) (*User, error) {
    return m.users[id], nil
}
```

**HTTP client mock:**

```go
type HTTPClient interface {
    Get(url string) (*Response, error)
}

type MockHTTPClient struct {
    responses map[string]*Response
}
```

## Real Example

```go
// Production code
type EmailSender interface {
    Send(to, subject, body string) error
}

type UserService struct {
    email EmailSender
}

func (s *UserService) Register(user *User) error {
    // ... validation ...
    return s.email.Send(user.Email, "Welcome", "...")
}

// Test code
type MockEmailSender struct {
    sent []Email
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    m.sent = append(m.sent, Email{to, subject, body})
    return nil
}

func TestRegister(t *testing.T) {
    mock := &MockEmailSender{}
    svc := &UserService{email: mock}

    svc.Register(&User{Email: "test@example.com"})

    if len(mock.sent) != 1 {
        t.Error("email not sent")
    }
}
```

## Gotchas

**Small interfaces:**

```go
type Reader interface {
    Read([]byte) (int, error)  // Small, focused
}
```

**Accept interfaces, return structs:**

```go
func NewService(db UserStore) *Service  // Accept interface
```
