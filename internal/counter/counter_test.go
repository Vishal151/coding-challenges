package counter

import (
	"io"
	"strings"
	"testing"
)

// Helper function to create a test reader
func createTestReader(content string) io.Reader {
	return strings.NewReader(content)
}

func TestCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Counts
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: Counts{Bytes: 0, Lines: 0, Words: 0, Chars: 0},
		},
		{
			name:     "Single word",
			input:    "hello",
			expected: Counts{Bytes: 5, Lines: 1, Words: 1, Chars: 5},
		},
		{
			name:     "Multiple words",
			input:    "hello world",
			expected: Counts{Bytes: 11, Lines: 1, Words: 2, Chars: 11},
		},
		{
			name:     "Multiple lines",
			input:    "hello\nworld\n",
			expected: Counts{Bytes: 12, Lines: 2, Words: 2, Chars: 12},
		},
		{
			name:     "Mixed content",
			input:    "Hello, World!\nThis is a test.",
			expected: Counts{Bytes: 30, Lines: 2, Words: 6, Chars: 29},
		},
		{
			name:     "Unicode characters",
			input:    "こんにちは\n世界\n",
			expected: Counts{Bytes: 23, Lines: 2, Words: 2, Chars: 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := createTestReader(tt.input)
			counts, err := Count(reader)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if counts != tt.expected {
				t.Errorf("Expected %+v, got %+v", tt.expected, counts)
			}
		})
	}
}

func TestCountLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single line", "hello", 1},
		{"Multiple lines", "hello\nworld\n", 2},
		{"Multiple lines without trailing newline", "hello\nworld", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := createTestReader(tt.input)
			count, err := CountLines(reader)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if count != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, count)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single word", "hello", 1},
		{"Multiple words", "hello world", 2},
		{"Multiple lines", "hello\nworld\n", 2},
		{"Mixed content", "Hello, World!\nThis is a test.", 6},
		{"Punctuation", "Hello, world! This, is a test.", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := createTestReader(tt.input)
			count, err := CountWords(reader)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if count != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, count)
			}
		})
	}
}

func TestCountChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single character", "a", 1},
		{"Multiple characters", "hello", 5},
		{"With newline", "hello\nworld", 11},
		{"Unicode characters", "こんにちは", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := createTestReader(tt.input)
			count, err := CountChars(reader)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if count != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, count)
			}
		})
	}
}
