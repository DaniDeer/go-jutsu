# Null Object Pattern

Replace nil checks with no-op implementation. Use interface + empty struct instead of scattered `if x != nil`.

## What It Is (and Isn't)

No-op implementation satisfies interface. No nil checks in business logic. Default safe behavior.

Not about avoiding nil everywhere. For optional dependencies. When nil = "do nothing".

**Critical distinction:**

- Null object for **optional side effects** (logging, analytics, notifications)
- Core business logic **always executes**
- Pattern skips non-critical operations, not main functionality

Example:

```go
func CreateOrder(analytics Analytics) {
    // Core logic - ALWAYS runs
    saveOrderToDatabase()
    chargePayment()

    // Side effect - can be no-op
    analytics.Track("order_created") // Skipped if NoOpAnalytics

    return "Order created" // Order IS created, analytics optional
}
```

## Where You See It

**Optional analytics:**

```go
type Tracker interface {
    Track(event string)
}

type NoOpTracker struct{}
func (NoOpTracker) Track(event string) {} // Does nothing

// No nil check needed
func ProcessOrder(tracker Tracker) {
    tracker.Track("order_processed") // Safe even if NoOpTracker
}
```

**Feature flags:**

```go
type Logger interface {
    Log(msg string)
}

type NullLogger struct{}
func (NullLogger) Log(msg string) {}

// Wire up based on config
func NewService(debugMode bool) *Service {
    logger := Logger(NullLogger{})
    if debugMode {
        logger = RealLogger{}
    }
    return &Service{logger: logger}
}
```

## Real Example

```go
type Notifier interface {
    Notify(user, message string) error
}

// Production implementation
type EmailNotifier struct {
    client *smtp.Client
}

func (e *EmailNotifier) Notify(user, message string) error {
    return e.client.SendEmail(user, message)
}

// Null object
type NoOpNotifier struct{}

func (NoOpNotifier) Notify(user, message string) error {
    return nil // Silent success
}

// Service code clean - no nil checks
type UserService struct {
    notifier Notifier
}

func (s *UserService) RegisterUser(user User) error {
    // ... registration logic ...

    // No if s.notifier != nil check
    return s.notifier.Notify(user.Email, "Welcome!")
}
```

## Gotchas

**When NOT to use:**

```go
// BAD: Core operation that must succeed
db, err := sql.Open(dsn)
if db == nil { // Need real error handling
    return err
}

// BAD: Operation failure matters
payment, err := ProcessPayment(amount)
if err != nil { // Can't use no-op - payment MUST succeed
    return err
}

// GOOD: Optional side effect/notification
analytics := Analytics(NoOpAnalytics{}) // Safe default
if cfg.EnableAnalytics {
    analytics = RealAnalytics{}
}

// GOOD: Non-critical logging
logger := Logger(NullLogger{}) // Silent OK
if cfg.DebugMode {
    logger = ConsoleLogger{}
}
```

**Use null object for:**

- Analytics tracking
- Debug logging
- Email notifications (when not critical)
- Metrics collection
- Audit trails (when optional)

**Don't use for:**

- Database connections
- Payment processing
- Authentication
- Data validation
- Any operation where failure matters

**Interface return values:**

```go
// Careful: nil interface != nil value
func GetTracker() Tracker {
    var t *RealTracker = nil
    return t // Returns non-nil interface!
}

// Prefer explicit null object
func GetTracker() Tracker {
    return NoOpTracker{}
}
```

**Error handling:**

```go
// Decide: silent success or explicit no-op error?
type NullValidator struct{}

func (NullValidator) Validate(data string) error {
    return nil // Silent success
    // Or: return ErrValidationSkipped
}
```

## Package Organization

### Where to Define Interface

**Option 1: With consumer (most common)**

```
myservice/
├── service.go       # Interface + service using it
└── null_notifier.go # Null object
```

**Option 2: Shared package**

```
notifier/
├── interface.go     # Interface
└── noop.go          # Null object

myservice/
└── service.go       # Import and use notifier.Notifier
```

**Go idiom:** Define interface where consumed, not where implemented.

### Where to Create Null Object

**Same package as interface (convenient):**

```go
// service/notifier.go
type Notifier interface {
    Notify(user, msg string) error
}

// service/null_notifier.go (same package)
type NullNotifier struct{}
func (NullNotifier) Notify(user, msg string) error { return nil }
```

**With other implementations:**

```go
// notifier/email.go
type Email struct{}

// notifier/slack.go
type Slack struct{}

// notifier/noop.go (alongside real ones)
type NoOp struct{}
```

### Where to Initialize

**In main (typical):**

```go
func main() {
    var notifier Notifier = NullNotifier{} // Default null object
    if config.EnableEmail {
        notifier = &EmailNotifier{...}
    }

    svc := NewUserService(notifier)
}
```

**In constructor:**

```go
func NewUserService(cfg Config) *UserService {
    notifier := Notifier(NullNotifier{}) // Default safe
    if cfg.EmailEnabled {
        notifier = NewEmailNotifier(cfg.SMTP)
    }
    return &UserService{notifier: notifier}
}
```

**In factory:**

```go
func NewNotifier(cfg Config) Notifier {
    switch cfg.Type {
    case "email": return &EmailNotifier{...}
    case "slack": return &SlackNotifier{...}
    default: return NullNotifier{} // Safe default
    }
}
```

**Rule:** Interface near consumer, null object nearby or with implementations. Initialize in main/factory based on config.

## Benefits

vs Nil Checks:

- Before: `if x != nil { x.Do() }` scattered everywhere
- After: `x.Do()` always safe

vs Optional/Maybe types:

- No generics/wrappers needed
- Pure Go interfaces
- Zero allocation (empty struct)

When optional dependency = "do nothing", null object > nil checks.

## References

- [Cleaner Go Code with the Null Object Pattern](https://medium.com/@mahmoud-magdy/cleaner-go-code-with-the-null-object-pattern-2464094de7bd) - Mahmoud Magdy's article on eliminating nil checks for optional services
