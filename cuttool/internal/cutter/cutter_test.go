package cutter

import (
	"os"
	"reflect"
	"testing"
)

func TestCutByFields(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		fields    []int
		delimiter string
		expected  string
		wantErr   bool
	}{
		{
			name:      "Basic case",
			input:     "a,b,c\nd,e,f\ng,h,i",
			fields:    []int{1, 3},
			delimiter: ",",
			expected:  "a,c\nd,f\ng,i\n",
			wantErr:   false,
		},
		{
			name:      "Single field",
			input:     "a,b,c\nd,e,f\ng,h,i",
			fields:    []int{2},
			delimiter: ",",
			expected:  "b\ne\nh\n",
			wantErr:   false,
		},
		{
			name:      "Out of range field",
			input:     "a,b,c\nd,e,f\ng,h,i",
			fields:    []int{1, 4},
			delimiter: ",",
			expected:  "a\nd\ng\n",
			wantErr:   false,
		},
		{
			name:      "Custom delimiter",
			input:     "a:b:c\nd:e:f\ng:h:i",
			fields:    []int{1, 3},
			delimiter: ":",
			expected:  "a:c\nd:f\ng:i\n",
			wantErr:   false,
		},
		{
			name:      "Empty input",
			input:     "",
			fields:    []int{1},
			delimiter: ",",
			expected:  "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CutByFields(tt.input, tt.fields, tt.delimiter)
			if (err != nil) != tt.wantErr {
				t.Errorf("CutByFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("CutByFields() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		expected string
	}{
		{
			name:     "Normal file",
			content:  "test content",
			wantErr:  false,
			expected: "test content",
		},
		{
			name:     "Empty file",
			content:  "",
			wantErr:  false,
			expected: "",
		},
		{
			name:     "Multi-line file",
			content:  "line1\nline2\nline3",
			wantErr:  false,
			expected: "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for testing
			tmpfile, err := os.CreateTemp("", "testfile")
			if err != nil {
				t.Fatalf("Cannot create temporary file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatalf("Failed to write to temporary file: %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Failed to close temporary file: %v", err)
			}

			// Test ReadFile function
			result, err := ReadFile(tmpfile.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReadFile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReadFile_NonExistentFile(t *testing.T) {
	_, err := ReadFile("non_existent_file.txt")
	if err == nil {
		t.Error("ReadFile() expected an error for non-existent file, but got nil")
	}
}
