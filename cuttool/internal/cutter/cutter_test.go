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
		fieldSpec string
		delimiter string
		expected  string
		wantErr   bool
	}{
		{
			name:      "Basic case",
			input:     "a,b,c\nd,e,f\ng,h,i",
			fieldSpec: "1,3",
			delimiter: ",",
			expected:  "a,c\nd,f\ng,i\n",
			wantErr:   false,
		},
		{
			name:      "Single field",
			input:     "a,b,c\nd,e,f\ng,h,i",
			fieldSpec: "2",
			delimiter: ",",
			expected:  "b\ne\nh\n",
			wantErr:   false,
		},
		{
			name:      "Out of range field",
			input:     "a,b,c\nd,e,f\ng,h,i",
			fieldSpec: "1,4",
			delimiter: ",",
			expected:  "a\nd\ng\n",
			wantErr:   false,
		},
		{
			name:      "Custom delimiter",
			input:     "a:b:c\nd:e:f\ng:h:i",
			fieldSpec: "1,3",
			delimiter: ":",
			expected:  "a:c\nd:f\ng:i\n",
			wantErr:   false,
		},
		{
			name:      "Empty input",
			input:     "",
			fieldSpec: "1",
			delimiter: ",",
			expected:  "",
			wantErr:   false,
		},
		{
			name:      "Field range",
			input:     "a,b,c,d,e\n1,2,3,4,5",
			fieldSpec: "2-4",
			delimiter: ",",
			expected:  "b,c,d\n2,3,4\n",
			wantErr:   false,
		},
		{
			name:      "Field range and single field",
			input:     "a,b,c,d,e\n1,2,3,4,5",
			fieldSpec: "1,3-5",
			delimiter: ",",
			expected:  "a,c,d,e\n1,3,4,5\n",
			wantErr:   false,
		},
		{
			name:      "Overlapping ranges",
			input:     "a,b,c,d,e\n1,2,3,4,5",
			fieldSpec: "1-3,2-4",
			delimiter: ",",
			expected:  "a,b,c,d\n1,2,3,4\n",
			wantErr:   false,
		},
		{
			name:      "Invalid field spec",
			input:     "a,b,c\n1,2,3",
			fieldSpec: "1,a",
			delimiter: ",",
			expected:  "",
			wantErr:   true,
		},
		{
			name:      "Invalid range spec",
			input:     "a,b,c\n1,2,3",
			fieldSpec: "1-3-5",
			delimiter: ",",
			expected:  "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CutByFields(tt.input, tt.fieldSpec, tt.delimiter)
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

func TestParseFieldSpec(t *testing.T) {
	tests := []struct {
		name      string
		fieldSpec string
		expected  []int
		wantErr   bool
	}{
		{"Single field", "3", []int{3}, false},
		{"Multiple fields", "1,3,5", []int{1, 3, 5}, false},
		{"Range", "2-5", []int{2, 3, 4, 5}, false},
		{"Mixed", "1,3-5,7", []int{1, 3, 4, 5, 7}, false},
		{"Overlapping", "1-3,2-4", []int{1, 2, 3, 4}, false},
		{"Unsorted", "5,1,3", []int{1, 3, 5}, false},
		{"Invalid field", "a", nil, true},
		{"Invalid range", "1-3-5", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseFieldSpec(tt.fieldSpec)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFieldSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseFieldSpec() = %v, want %v", result, tt.expected)
			}
		})
	}
}
