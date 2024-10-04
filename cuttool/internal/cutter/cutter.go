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
func CutByFields(input string, fieldSpec string, delimiter string, onlyDelimited bool) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var result strings.Builder

	fields, err := parseFieldSpec(fieldSpec)
	if err != nil {
		return "", err
	}

	for scanner.Scan() {
		line := scanner.Text()
		if onlyDelimited && !strings.Contains(line, delimiter) {
			continue
		}
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

// CutByBytes cuts the input by byte ranges
// Step 3: Implement cutting by byte ranges
func CutByBytes(input string, byteSpec string) (string, error) {
	ranges, err := parseByteSpec(byteSpec)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()
		lineBytes := []byte(line)
		selectedBytes := make([]bool, len(lineBytes))

		for _, r := range ranges {
			start, end := r[0], r[1]
			if start < 1 {
				start = 1
			}
			if end == -1 || end > len(lineBytes) {
				end = len(lineBytes)
			}
			for i := start - 1; i < end && i < len(lineBytes); i++ {
				selectedBytes[i] = true
			}
		}

		for i, selected := range selectedBytes {
			if selected {
				result.WriteByte(lineBytes[i])
			}
		}
		result.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}

	return result.String(), nil
}

// parseByteSpec parses the byte specification string into a sorted list of byte ranges
func parseByteSpec(byteSpec string) ([][2]int, error) {
	var ranges [][2]int
	specs := strings.Split(byteSpec, ",")

	for _, spec := range specs {
		if strings.Contains(spec, "-") {
			parts := strings.Split(spec, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid byte range specification: %s", spec)
			}

			start, err := strconv.Atoi(parts[0])
			if err != nil && parts[0] != "" {
				return nil, fmt.Errorf("invalid start byte: %s", parts[0])
			}

			end, err := strconv.Atoi(parts[1])
			if err != nil && parts[1] != "" {
				return nil, fmt.Errorf("invalid end byte: %s", parts[1])
			}

			if parts[0] == "" {
				start = 1
			}
			if parts[1] == "" {
				end = -1 // Indicates "to the end of the line"
			}

			ranges = append(ranges, [2]int{start, end})
		} else {
			pos, err := strconv.Atoi(spec)
			if err != nil {
				return nil, fmt.Errorf("invalid byte position: %s", spec)
			}
			ranges = append(ranges, [2]int{pos, pos})
		}
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] < ranges[j][0]
	})

	return ranges, nil
}

// ReadFile reads the content of a file
func ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	return string(content), nil
}
