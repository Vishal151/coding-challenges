package main

import (
	"fmt"
	"os"

	"github.com/vishal151/jsonparser/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jsonparser <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	valid, err := parser.ValidateJSONStep4(data)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	if valid {
		fmt.Println("Valid JSON")
		os.Exit(0)
	} else {
		fmt.Println("Invalid JSON")
		os.Exit(1)
	}
}
