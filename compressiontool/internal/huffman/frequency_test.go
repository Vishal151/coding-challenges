package huffman

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCountFrequencies(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[byte]int
	}{
		{
			name:     "Empty input",
			input:    "",
			expected: map[byte]int{},
		},
		{
			name:     "Single character",
			input:    "a",
			expected: map[byte]int{'a': 1},
		},
		{
			name:     "Multiple characters",
			input:    "abracadabra",
			expected: map[byte]int{'a': 5, 'b': 2, 'r': 2, 'c': 1, 'd': 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.input))
			result, err := CountFrequencies(reader)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
