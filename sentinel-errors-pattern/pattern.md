# Sentinel Errors Pattern

Predefined error values for comparison. Use `errors.Is()` for checking.

## What It Is (and Isn't)

Package-level error variables. Compare with `errors.Is()` (not `==`). Standard library pattern.

Not error strings. Not types. Specific error values.

## Where You See It

**Standard library:**

```go
if errors.Is(err, io.EOF) {
    // End of file
}

if errors.Is(err, sql.ErrNoRows) {
    // No rows found
}
```

**Custom sentinels:**

```go
var (
    ErrNotFound = errors.New("not found")
    ErrInvalid  = errors.New("invalid input")
    ErrExists   = errors.New("already exists")
)

func Get(id int) error {
    if !exists(id) {
        return ErrNotFound
    }
    return nil
}

// Usage
if errors.Is(err, ErrNotFound) {
    // Handle not found
}
```

## Real Example

```go
package user

import "errors"

var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidEmail    = errors.New("invalid email")
    ErrDuplicateEmail  = errors.New("email already exists")
    ErrUnauthorized    = errors.New("unauthorized")
)

func FindByEmail(email string) (*User, error) {
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }

    user := db.Query(email)
    if user == nil {
        return nil, ErrUserNotFound
    }

    return user, nil
}

// Caller
user, err := FindByEmail("test@example.com")
if errors.Is(err, user.ErrUserNotFound) {
    return http.StatusNotFound
}
if errors.Is(err, user.ErrInvalidEmail) {
    return http.StatusBadRequest
}
```

## Gotchas

**Use errors.Is, not ==:**

```go
// BAD
if err == ErrNotFound {  // Breaks with wrapped errors

// GOOD
if errors.Is(err, ErrNotFound) {  // Works with wrapping
```

**Wrapping preserves sentinel:**

```go
err := fmt.Errorf("database: %w", ErrNotFound)
errors.Is(err, ErrNotFound)  // true
```

**Export sentinels:**

```go
// Public API
var ErrNotFound = errors.New("not found")

// Private helper
var errInternal = errors.New("internal")
```

## Benefits

- Type-safe error handling
- Clear API contracts
- Works with error wrapping
- Self-documenting code
