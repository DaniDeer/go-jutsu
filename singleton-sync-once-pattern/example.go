package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Example 1: Basic singleton with sync.Once
type Config struct {
	AppName string
	Version string
}

var (
	config Config
	once   sync.Once
)

func GetConfig() Config {
	once.Do(func() {
		fmt.Println("  [Initializing config...]")
		time.Sleep(100 * time.Millisecond) // Simulate expensive load
		config = Config{
			AppName: "MyApp",
			Version: "1.0.0",
		}
	})
	return config
}

// Example 2: Singleton with error handling
type Logger struct {
	file *os.File
}

var (
	logger    *Logger
	loggerErr error
	logOnce   sync.Once
)

func GetLogger() (*Logger, error) {
	logOnce.Do(func() {
		fmt.Println("  [Opening log file...]")
		file, err := os.CreateTemp("", "app-*.log")
		if err != nil {
			loggerErr = err
			return
		}
		logger = &Logger{file: file}
	})
	return logger, loggerErr
}

func (l *Logger) Log(msg string) {
	if l != nil && l.file != nil {
		fmt.Fprintf(l.file, "[%s] %s\\n", time.Now().Format(time.RFC3339), msg)
	}
}

func (l *Logger) Close() error {
	if l != nil && l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Example 3: Database singleton
type Database struct {
	dsn string
}

var (
	db     *Database
	dbOnce sync.Once
)

func GetDB() *Database {
	dbOnce.Do(func() {
		fmt.Println("  [Connecting to database...]")
		time.Sleep(200 * time.Millisecond) // Simulate connection
		db = &Database{dsn: "postgres://localhost/mydb"}
	})
	return db
}

func main() {
	// Example 1: Basic singleton - called multiple times
	fmt.Println("=== Example 1: Config Singleton ===")
	for i := 0; i < 3; i++ {
		cfg := GetConfig()
		fmt.Printf("Call %d: %s v%s\\n", i+1, cfg.AppName, cfg.Version)
	}

	// Example 2: Concurrent access (initialization only once)
	fmt.Println("\\n=== Example 2: Concurrent Singleton ===")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db := GetDB()
			fmt.Printf("Goroutine %d: DB DSN = %s\\n", id, db.dsn)
		}(i)
	}
	wg.Wait()

	// Example 3: Singleton with error handling
	fmt.Println("\\n=== Example 3: Logger Singleton ===")
	log, err := GetLogger()
	if err != nil {
		fmt.Printf("Error: %v\\n", err)
		return
	}
	defer log.Close()

	log.Log("Application started")
	log.Log("Processing request")

	// Verify file contents
	logFile := log.file.Name()
	log.Close()
	content, _ := os.ReadFile(logFile)
	fmt.Printf("Log contents:\\n%s", content)
	os.Remove(logFile)

	// Example 4: Demonstrate sync.Once never resets
	fmt.Println("\\n=== Example 4: Once Never Resets ===")
	var counter int
	var testOnce sync.Once

	testOnce.Do(func() {
		counter++
		fmt.Println("  First Do() call")
	})

	testOnce.Do(func() {
		counter++
		fmt.Println("  Second Do() call (never printed)")
	})

	fmt.Printf("Counter: %d (should be 1)\\n", counter)

	// Example 5: Benchmark - compare with regular initialization
	fmt.Println("\\n=== Example 5: Performance Note ===")
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		_ = GetConfig() // Cached, very fast
	}
	fmt.Printf("1M singleton calls: %v\\n", time.Since(start))
}
