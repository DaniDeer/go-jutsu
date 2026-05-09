package main

import (
	"fmt"
	"strings"
)

// Function to test
func Square(n int) int {
	return n * n
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func ParseURL(url string) (map[string]string, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("invalid scheme")
	}
	parts := strings.Split(strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://"), "/")
	return map[string]string{"host": parts[0]}, nil
}

func main() {
	// Example 1: Simple table test
	fmt.Println("=== Example 1: Square Function ===")
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"zero", 0, 0},
		{"positive", 5, 25},
		{"negative", -3, 9},
		{"one", 1, 1},
	}

	for _, tt := range tests {
		got := Square(tt.input)
		if got != tt.want {
			fmt.Printf("FAIL %s: Square(%d) = %d, want %d\\n",
				tt.name, tt.input, got, tt.want)
		} else {
			fmt.Printf("PASS %s\\n", tt.name)
		}
	}

	// Example 2: With error handling
	fmt.Println("\\n=== Example 2: Divide Function ===")
	divTests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
	}{
		{"normal", 10, 2, 5, false},
		{"by zero", 10, 0, 0, true},
		{"negative", -10, 2, -5, false},
	}

	for _, tt := range divTests {
		got, err := Divide(tt.a, tt.b)
		hasErr := err != nil
		if hasErr != tt.wantErr {
			fmt.Printf("FAIL %s: error = %v, wantErr %v\\n",
				tt.name, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			fmt.Printf("FAIL %s: got %d, want %d\\n", tt.name, got, tt.want)
			continue
		}
		fmt.Printf("PASS %s\\n", tt.name)
	}

	// Example 3: Multiple assertions
	fmt.Println("\\n=== Example 3: URL Parser ===")
	urlTests := []struct {
		name     string
		input    string
		wantErr  bool
		wantHost string
	}{
		{
			name:     "valid http",
			input:    "http://example.com/path",
			wantErr:  false,
			wantHost: "example.com",
		},
		{
			name:    "invalid scheme",
			input:   "ftp://example.com",
			wantErr: true,
		},
		{
			name:     "https",
			input:    "https://secure.com",
			wantErr:  false,
			wantHost: "secure.com",
		},
	}

	for _, tt := range urlTests {
		got, err := ParseURL(tt.input)
		hasErr := err != nil

		if hasErr != tt.wantErr {
			fmt.Printf("FAIL %s: error = %v, wantErr %v\\n",
				tt.name, err, tt.wantErr)
			continue
		}

		if !tt.wantErr {
			if got["host"] != tt.wantHost {
				fmt.Printf("FAIL %s: host = %s, want %s\\n",
					tt.name, got["host"], tt.wantHost)
				continue
			}
		}

		fmt.Printf("PASS %s\\n", tt.name)
	}
}
