package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vishal151/ccwc/internal/counter"
)

var (
	countBytes bool
	countLines bool
	countWords bool
	countChars bool
)

var rootCmd = &cobra.Command{
	Use:   "ccwc [flags] [file]",
	Short: "Count bytes, lines, words, and characters in a file",
	Run:   runCount,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&countBytes, "bytes", "c", false, "Count bytes")
	rootCmd.Flags().BoolVarP(&countLines, "lines", "l", false, "Count lines")
	rootCmd.Flags().BoolVarP(&countWords, "words", "w", false, "Count words")
	rootCmd.Flags().BoolVarP(&countChars, "chars", "m", false, "Count characters")
}

func runCount(cmd *cobra.Command, args []string) {
	var input *os.File
	var err error
	filename := "stdin"

	if len(args) > 0 {
		input, err = os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer input.Close()
		filename = args[0]
	} else {
		input = os.Stdin
	}

	counts, err := counter.Count(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error counting: %v\n", err)
		os.Exit(1)
	}

	if !countBytes && !countLines && !countWords && !countChars {
		fmt.Printf("%7d %7d %7d %s\n", counts.Lines, counts.Words, counts.Bytes, filename)
	} else {
		if countBytes {
			fmt.Printf("%d ", counts.Bytes)
		}
		if countLines {
			fmt.Printf("%d ", counts.Lines)
		}
		if countWords {
			fmt.Printf("%d ", counts.Words)
		}
		if countChars {
			fmt.Printf("%d ", counts.Chars)
		}
		fmt.Printf("%s\n", filename)
	}
}
