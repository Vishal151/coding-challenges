package cutter

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// CutByFields cuts the input by fields
// Step 2: Implement cutting by fields, including ranges and lists
func CutByFields(input string, fieldSpec string, delimiter string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var result strings.Builder

	fields, err := parseFieldSpec(fieldSpec)
	if err != nil {
		return "", err
	}

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

// parseFieldSpec parses the field specification string into a sorted list of unique field numbers
func parseFieldSpec(fieldSpec string) ([]int, error) {
	var fields []int
	specs := strings.Split(fieldSpec, ",")

	for _, spec := range specs {
		if strings.Contains(spec, "-") {
			// Handle range
			rangeParts := strings.Split(spec, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range specification: %s", spec)
			}
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid range start: %s", rangeParts[0])
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range end: %s", rangeParts[1])
			}
			for i := start; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			// Handle single field
			field, err := strconv.Atoi(spec)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", spec)
			}
			fields = append(fields, field)
		}
	}

	// Sort and remove duplicates
	sort.Ints(fields)
	uniqueFields := []int{}
	for i, field := range fields {
		if i == 0 || field != fields[i-1] {
			uniqueFields = append(uniqueFields, field)
		}
	}

	return uniqueFields, nil
}

// ReadFile reads the content of a file
func ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return string(content), nil
}
