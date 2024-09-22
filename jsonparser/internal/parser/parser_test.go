package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateJSONStep1(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
		wantErr  bool
	}{
		{"Valid Empty Object", "tests/step1/valid.json", true, false},
		{"Invalid Empty", "tests/step1/invalid.json", false, true},
	}

	runTests(t, tests, ValidateJSONStep1)
}

func TestValidateJSONStep2(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
		wantErr  bool
	}{
		{"Step 2 - Valid", "tests/step2/valid.json", true, false},
		{"Step 2 - Valid 2", "tests/step2/valid2.json", true, false},
		{"Step 2 - Invalid", "tests/step2/invalid.json", false, true},
		{"Step 2 - Invalid 2", "tests/step2/invalid2.json", false, true},
	}

	runTests(t, tests, ValidateJSONStep2)
}

func TestValidateJSONStep3(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
		wantErr  bool
	}{
		{"Step 3 - Valid", "tests/step3/valid.json", true, false},
		{"Step 3 - Invalid", "tests/step3/invalid.json", false, true},
	}

	runTests(t, tests, ValidateJSONStep3)
}

func TestValidateJSONStep4(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
		wantErr  bool
	}{
		{"Step 4 - Valid", "tests/step4/valid.json", true, false},
		{"Step 4 - Valid 2", "tests/step4/valid2.json", true, false},
		{"Step 4 - Invalid", "tests/step4/invalid.json", false, true},
	}

	runTests(t, tests, ValidateJSONStep4)
}

func runTests(t *testing.T, tests []struct {
	name     string
	filename string
	want     bool
	wantErr  bool
}, validateFunc func([]byte) (bool, error)) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "..", tt.filename))
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			got, err := validateFunc(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
