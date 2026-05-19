package main

import "fmt"

// Analytics interface
type Analytics interface {
	Track(event string, data map[string]interface{})
}

// Real implementation
type MixpanelAnalytics struct {
	apiKey string
}

func (m *MixpanelAnalytics) Track(event string, data map[string]interface{}) {
	fmt.Printf("[Mixpanel] Event: %s, Data: %v\n", event, data)
}

// Null object - no-op implementation
type NoOpAnalytics struct{}

func (NoOpAnalytics) Track(event string, data map[string]interface{}) {
	// Do nothing
}

// Service using analytics
type OrderService struct {
	analytics Analytics
}

func (s *OrderService) CreateOrder(userID int, amount float64) {
	fmt.Printf("Creating order for user %d: $%.2f\n", userID, amount)

	// Core business logic - ALWAYS executes
	// ... save to database ...
	// ... charge payment ...

	// Analytics is optional side effect - no-op OK
	s.analytics.Track("order_created", map[string]interface{}{
		"user_id": userID,
		"amount":  amount,
	})

	fmt.Println("Order created successfully") // Order IS created, analytics optional
}

// Logger interface
type Logger interface {
	Info(msg string)
	Error(msg string)
}

// Real logger
type ConsoleLogger struct{}

func (ConsoleLogger) Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func (ConsoleLogger) Error(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
}

// Null logger
type NullLogger struct{}

func (NullLogger) Info(msg string)  {}
func (NullLogger) Error(msg string) {}

// Notifier interface - for operations you WANT to skip
type Notifier interface {
	SendEmail(to, subject, body string) error
}

// Real notifier
type EmailNotifier struct{}

func (EmailNotifier) SendEmail(to, subject, body string) error {
	fmt.Printf("[EMAIL] To: %s, Subject: %s\n", to, subject)
	return nil
}

// Null notifier - intentionally skips sending
type NoOpNotifier struct{}

func (NoOpNotifier) SendEmail(to, subject, body string) error {
	// Intentionally do nothing - no email sent
	return nil
}

// Service that sends notifications
type UserService struct {
	notifier Notifier
}

func (s *UserService) RegisterUser(email, name string) error {
	fmt.Printf("Registering user: %s\n", name)

	// Core logic - ALWAYS happens
	// ... save user to database ...

	// Notification is optional - can be no-op
	err := s.notifier.SendEmail(email, "Welcome!", "Welcome to our service")
	if err != nil {
		return err
	}

	fmt.Println("User registered successfully") // User IS registered, email optional
	return nil
}

// Config for initialization examples
type Config struct {
	EnableAnalytics bool
	AnalyticsType   string
	APIKey          string
}

// OPTION 1: Initialize in main
// Direct initialization - main decides based on config
func createServiceInMain(cfg Config) *OrderService {
	// Main creates concrete type based on config
	var analytics Analytics = NoOpAnalytics{} // Default null object
	if cfg.EnableAnalytics {
		analytics = &MixpanelAnalytics{apiKey: cfg.APIKey}
	}

	return &OrderService{analytics: analytics}
}

// OPTION 2: Initialize in constructor
// Constructor handles default - caller doesn't see null object
func NewOrderService(cfg Config) *OrderService {
	// Internal default to null object
	analytics := Analytics(NoOpAnalytics{})
	if cfg.EnableAnalytics {
		analytics = &MixpanelAnalytics{apiKey: cfg.APIKey}
	}

	return &OrderService{analytics: analytics}
}

// OPTION 3: Initialize with factory
// Factory function creates appropriate implementation
func NewAnalytics(cfg Config) Analytics {
	if !cfg.EnableAnalytics {
		return NoOpAnalytics{} // Safe default
	}

	switch cfg.AnalyticsType {
	case "mixpanel":
		return &MixpanelAnalytics{apiKey: cfg.APIKey}
	// case "segment": return &SegmentAnalytics{...}
	default:
		return NoOpAnalytics{}
	}
}

// Use with factory
func NewOrderServiceWithFactory(cfg Config) *OrderService {
	analytics := NewAnalytics(cfg) // Factory handles logic
	return &OrderService{analytics: analytics}
}

// Comparison: with vs without pattern
type OldOrderService struct {
	analytics *MixpanelAnalytics // Pointer can be nil
}

func (s *OldOrderService) CreateOrderOld(userID int, amount float64) {
	fmt.Printf("Creating order (old way) for user %d: $%.2f\n", userID, amount)

	// Scattered nil checks everywhere
	if s.analytics != nil {
		s.analytics.Track("order_created", map[string]interface{}{
			"user_id": userID,
			"amount":  amount,
		})
	}

	fmt.Println("Order created successfully")
}

func main() {
	fmt.Println("=== KEY CONCEPT: Null Object for OPTIONAL Side Effects ===\n")

	// Analytics example - order creation ALWAYS happens
	fmt.Println("Example: Order with analytics (side effect)")
	svc := &OrderService{analytics: &MixpanelAnalytics{apiKey: "key"}}
	svc.CreateOrder(100, 50.00)
	fmt.Println("→ Order created, analytics tracked\n")

	svcNoOp := &OrderService{analytics: NoOpAnalytics{}}
	svcNoOp.CreateOrder(101, 60.00)
	fmt.Println("→ Order STILL created, analytics skipped (no-op)\n")

	// Notification example - user registration ALWAYS happens
	fmt.Println("Example: User registration with email notification")
	userSvc := &UserService{notifier: EmailNotifier{}}
	userSvc.RegisterUser("alice@example.com", "Alice")
	fmt.Println("→ User registered, email sent\n")

	userSvcNoOp := &UserService{notifier: NoOpNotifier{}}
	userSvcNoOp.RegisterUser("bob@example.com", "Bob")
	fmt.Println("→ User STILL registered, email skipped (no-op)\n")

	fmt.Println("=== INITIALIZATION PATTERNS ===\n")

	// OPTION 1: In main
	fmt.Println("Option 1: Initialize in main")
	cfg1 := Config{EnableAnalytics: true, APIKey: "main-key"}
	svc1 := createServiceInMain(cfg1)
	svc1.CreateOrder(200, 75.00)

	cfg1NoOp := Config{EnableAnalytics: false}
	svc1NoOp := createServiceInMain(cfg1NoOp)
	svc1NoOp.CreateOrder(201, 85.00)
	fmt.Println("→ Main controls initialization logic\n")

	// OPTION 2: In constructor
	fmt.Println("Option 2: Initialize in constructor")
	cfg2 := Config{EnableAnalytics: true, APIKey: "constructor-key"}
	svc2 := NewOrderService(cfg2)
	svc2.CreateOrder(300, 99.00)

	cfg2NoOp := Config{EnableAnalytics: false}
	svc2NoOp := NewOrderService(cfg2NoOp)
	svc2NoOp.CreateOrder(301, 109.00)
	fmt.Println("→ Constructor handles default safely\n")

	// OPTION 3: With factory
	fmt.Println("Option 3: Initialize with factory")
	cfg3 := Config{
		EnableAnalytics: true,
		AnalyticsType:   "mixpanel",
		APIKey:          "factory-key",
	}
	svc3 := NewOrderServiceWithFactory(cfg3)
	svc3.CreateOrder(400, 125.00)

	cfg3NoOp := Config{EnableAnalytics: false}
	svc3NoOp := NewOrderServiceWithFactory(cfg3NoOp)
	svc3NoOp.CreateOrder(401, 135.00)
	fmt.Println("→ Factory encapsulates creation logic\n")

	fmt.Println("=== COMPARISON: OLD WAY ===\n")
	oldSvc1 := &OldOrderService{analytics: &MixpanelAnalytics{apiKey: "test"}}
	oldSvc1.CreateOrderOld(500, 50.00)

	oldSvc2 := &OldOrderService{analytics: nil}
	oldSvc2.CreateOrderOld(501, 75.00)

	fmt.Println("\n=== SUMMARY ===")
	fmt.Println("✓ Null object = skip OPTIONAL operations (analytics, logging, notifications)")
	fmt.Println("✓ Core business logic ALWAYS executes")
	fmt.Println("✓ Eliminates scattered nil checks for side effects")
	fmt.Println("✓ Option 1 (main): Full control, explicit")
	fmt.Println("✓ Option 2 (constructor): Encapsulated, simple API")
	fmt.Println("✓ Option 3 (factory): Extensible, multiple types")
}
