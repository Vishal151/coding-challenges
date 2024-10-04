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

func TestCutByBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		byteSpec string
		expected string
		wantErr  bool
	}{
		{
			name:     "Single byte",
			input:    "abcde\nfghij\nklmno",
			byteSpec: "3",
			expected: "c\nh\nm\n",
			wantErr:  false,
		},
		{
			name:     "Byte range",
			input:    "abcde\nfghij\nklmno",
			byteSpec: "2-4",
			expected: "bcd\nghi\nlmn\n",
			wantErr:  false,
		},
		{
			name:     "Multiple ranges",
			input:    "abcde\nfghij\nklmno",
			byteSpec: "1-2,4-5",
			expected: "abde\nfgij\nklno\n",
			wantErr:  false,
		},
		{
			name:     "Open-ended range",
			input:    "abcde\nfghij\nklmno",
			byteSpec: "3-",
			expected: "cde\nhij\nmno\n",
			wantErr:  false,
		},
		{
			name:     "Range from beginning",
			input:    "abcde\nfghij\nklmno",
			byteSpec: "-3",
			expected: "abc\nfgh\nklm\n",
			wantErr:  false,
		},
		{
			name:     "Out of range",
			input:    "abcde\nfghij\nklmno",
			byteSpec: "1-10",
			expected: "abcde\nfghij\nklmno\n",
			wantErr:  false,
		},
		{
			name:     "Mixed ranges",
			input:    "abcde\nfghij\nklmno",
			byteSpec: "1,3-4,2-",
			expected: "abcde\nfghij\nklmno\n",
			wantErr:  false,
		},
		{
			name:     "Overlapping ranges",
			input:    "abcde\nfghij\nklmno",
			byteSpec: "1-3,2-4",
			expected: "abcd\nfghi\nklmn\n",
			wantErr:  false,
		},
		{
			name:     "Single character lines",
			input:    "a\nb\nc",
			byteSpec: "1-",
			expected: "a\nb\nc\n",
			wantErr:  false,
		},
		{
			name:     "Empty lines",
			input:    "abc\n\ndef",
			byteSpec: "1-",
			expected: "abc\n\ndef\n",
			wantErr:  false,
		},
		{
			name:     "Invalid byte spec",
			input:    "abcde",
			byteSpec: "a-b",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CutByBytes(tt.input, tt.byteSpec)
			if (err != nil) != tt.wantErr {
				t.Errorf("CutByBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("CutByBytes() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestParseByteSpec(t *testing.T) {
	tests := []struct {
		name     string
		byteSpec string
		expected [][2]int
		wantErr  bool
	}{
		{"Single byte", "3", [][2]int{{3, 3}}, false},
		{"Byte range", "2-4", [][2]int{{2, 4}}, false},
		{"Multiple ranges", "1-2,4-5", [][2]int{{1, 2}, {4, 5}}, false},
		{"Open-ended range", "3-", [][2]int{{3, -1}}, false},
		{"Range from beginning", "-3", [][2]int{{1, 3}}, false},
		{"Mixed specifications", "1,3-5,7-", [][2]int{{1, 1}, {3, 5}, {7, -1}}, false},
		{"Invalid spec", "a-b", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseByteSpec(tt.byteSpec)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseByteSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseByteSpec() = %v, want %v", result, tt.expected)
			}
		})
	}
}
