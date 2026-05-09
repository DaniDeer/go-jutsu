package main

import (
	"fmt"
	"sync"
)

// Logger interface
type Logger interface {
	Log(msg string)
}

// Concrete implementations
type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(msg string) {
	fmt.Printf("[CONSOLE] %s\n", msg)
}

type FileLogger struct {
	path string
}

func (l *FileLogger) Log(msg string) {
	fmt.Printf("[FILE:%s] %s\n", l.path, msg)
}

type NullLogger struct{}

func (l *NullLogger) Log(msg string) {
	// No-op
}

// Factory
type LoggerFactory func(config map[string]string) Logger

var (
	loggerRegistry = make(map[string]LoggerFactory)
	mu             sync.RWMutex
)

func RegisterLogger(name string, factory LoggerFactory) {
	mu.Lock()
	defer mu.Unlock()
	loggerRegistry[name] = factory
}

func NewLogger(name string, config map[string]string) (Logger, error) {
	mu.RLock()
	factory, ok := loggerRegistry[name]
	mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown logger type: %s", name)
	}

	return factory(config), nil
}

func init() {
	RegisterLogger("console", func(cfg map[string]string) Logger {
		return &ConsoleLogger{}
	})

	RegisterLogger("file", func(cfg map[string]string) Logger {
		path := cfg["path"]
		if path == "" {
			path = "app.log"
		}
		return &FileLogger{path: path}
	})

	RegisterLogger("null", func(cfg map[string]string) Logger {
		return &NullLogger{}
	})
}

func main() {
	fmt.Println("=== Factory Pattern ===")

	// Create different loggers
	console, _ := NewLogger("console", nil)
	console.Log("This goes to console")

	file, _ := NewLogger("file", map[string]string{"path": "/var/log/app.log"})
	file.Log("This goes to file")

	null, _ := NewLogger("null", nil)
	null.Log("This is silenced")

	// Try unknown type
	_, err := NewLogger("unknown", nil)
	if err != nil {
		fmt.Printf("✓ Error caught: %v\n", err)
	}

	// Dynamic registration (plugin system)
	fmt.Println("\nDynamic registration:")
	RegisterLogger("custom", func(cfg map[string]string) Logger {
		return &ConsoleLogger{} // Custom implementation
	})

	custom, _ := NewLogger("custom", nil)
	custom.Log("Custom logger working")
}
