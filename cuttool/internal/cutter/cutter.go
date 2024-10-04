package cutter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CutByFields cuts the input by fields
// Step 1: Implement cutting by fields
func CutByFields(input string, fields []int, delimiter string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var result strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, delimiter)
		var selectedParts []string

		for _, field := range fields {
			if field > 0 && field <= len(parts) {
				selectedParts = append(selectedParts, parts[field-1])
			}
		}

		result.WriteString(strings.Join(selectedParts, delimiter) + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}

	return result.String(), nil
}

// ReadFile reads the content of a file
func ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return string(content), nil
}
