package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// Custom errors
var (
	ErrValidation = errors.New("validation error")
	ErrProcessing = errors.New("processing error")
)

// ProcessFile wraps errors with context
func ProcessFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ProcessFile(%q): %w", path, err)
	}

	if err := validateData(data); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := processData(data); err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}

	return nil
}

func validateData(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty file: %w", ErrValidation)
	}
	return nil
}

func processData(data []byte) error {
	if len(data) > 1000 {
		return fmt.Errorf("file too large: %w", ErrProcessing)
	}
	return nil
}

func main() {
	// Example 1: File not found (wrapped)
	fmt.Println("=== Example 1: File Not Found ===")
	err := ProcessFile("nonexistent.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		// Check original error through wrapping
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("✓ Detected: file doesn't exist")
		}

		// Extract PathError type
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			fmt.Printf("✓ PathError Op: %s, Path: %s\n", pathErr.Op, pathErr.Path)
		}
	}

	// Example 2: Create temp file and validate
	fmt.Println("\n=== Example 2: Validation Error ===")
	tmpFile, _ := os.CreateTemp("", "test")
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = ProcessFile(tmpFile.Name())
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		if errors.Is(err, ErrValidation) {
			fmt.Println("✓ Detected: validation error")
		}
	}

	// Example 3: Large file processing error
	fmt.Println("\n=== Example 3: Processing Error ===")
	largeFile, _ := os.CreateTemp("", "large")
	largeFile.Write(make([]byte, 1500))
	largeFile.Close()
	defer os.Remove(largeFile.Name())

	err = ProcessFile(largeFile.Name())
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		if errors.Is(err, ErrProcessing) {
			fmt.Println("✓ Detected: processing error")
		}
	}

	// Example 4: Demonstrate Unwrap
	fmt.Println("\n=== Example 4: Manual Unwrap ===")
	err1 := errors.New("base error")
	err2 := fmt.Errorf("level 2: %w", err1)
	err3 := fmt.Errorf("level 3: %w", err2)

	fmt.Printf("err3: %v\n", err3)
	unwrapped := errors.Unwrap(err3)
	fmt.Printf("Unwrapped once: %v\n", unwrapped)
	fmt.Printf("Is original? %v\n", errors.Is(err3, err1))
}
