# Table-Driven Design Pattern

Data structures drive logic. Tests, configs, validators as tables. Idiomatic Go.

## What It Is (and Isn't)

Define behavior via data tables. Loop over table rows. Extremely common in Go testing.

Not switch statements. Not if-else chains. Declarative data-driven approach.

## Where You See It

**Table-driven tests:**

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 5, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("got %d, want %d", got, tt.want)
            }
        })
    }
}
```

**Configuration:**

```go
var routes = []struct {
    method  string
    path    string
    handler http.HandlerFunc
}{
    {"GET", "/users", listUsers},
    {"POST", "/users", createUser},
    {"GET", "/users/:id", getUser},
}

for _, route := range routes {
    mux.HandleFunc(route.method, route.path, route.handler)
}
```

**Validation rules:**

```go
var validators = []struct {
    field string
    rule  func(string) bool
    msg   string
}{
    {"email", isEmail, "invalid email"},
    {"phone", isPhone, "invalid phone"},
}
```

## Real Example

```go
// HTTP status code mapping
var statusMessages = map[int]struct {
    message string
    retry   bool
}{
    200: {"OK", false},
    404: {"Not Found", false},
    500: {"Server Error", true},
    503: {"Unavailable", true},
}

func handleStatus(code int) {
    info, ok := statusMessages[code]
    if !ok {
        log.Printf("Unknown status code: %d", code)
        return
    }

    log.Printf("%d: %s", code, info.message)
    if info.retry {
        scheduleRetry()
    }
}

// State machine as table
type State int

const (
    StatePending State = iota
    StateRunning
    StateDone
    StateFailed
)

var transitions = map[State][]State{
    StatePending: {StateRunning, StateFailed},
    StateRunning: {StateDone, StateFailed},
    StateDone:    {},
    StateFailed:  {StatePending}, // Retry
}

func canTransition(from, to State) bool {
    allowed := transitions[from]
    for _, s := range allowed {
        if s == to {
            return true
        }
    }
    return false
}
```

## Gotchas

**Shared state in tests:**

```go
// BAD: mutations affect other tests
tests := []struct {
    input []int
}{
    {[]int{1, 2}},
}

for _, tt := range tests {
    modify(tt.input)  // Modifies original!
}

// GOOD: copy or reinitialize
{[]int{1, 2}},  // Each test gets fresh slice
```

**Table order matters:**

```go
// Execution order matters for setup/teardown
tests := []struct {
    name  string
    setup func()
    test  func()
}{
    {"first", setupDB, testQuery},
    {"second", nil, testQuery},  // Assumes first setup
}
```

## Benefits

- DRY (Don't Repeat Yourself)
- Easy to add test cases
- Self-documenting
- Maintainable
- Standard Go idiom
