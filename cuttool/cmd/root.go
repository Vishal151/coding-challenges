package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vishal151/cut-tool/internal/cutter"
)

var fields string
var delimiter string

var rootCmd = &cobra.Command{
	Use:   "cut-tool [file]",
	Short: "A cut tool implementation",
	Long:  `A cut tool implementation for the coding challenge at https://codingchallenges.fyi/challenges/challenge-cut`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		content, err := cutter.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}

		fieldsList, err := parseFields(fields)
		if err != nil {
			fmt.Println("Error parsing fields:", err)
			os.Exit(1)
		}

		result, err := cutter.CutByFields(content, fieldsList, delimiter)
		if err != nil {
			fmt.Println("Error cutting fields:", err)
			os.Exit(1)
		}

		fmt.Print(result)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&fields, "fields", "f", "", "select only these fields")
	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "\t", "use DELIM instead of TAB for field delimiter")
	rootCmd.MarkFlagRequired("fields")
}

func parseFields(fieldsStr string) ([]int, error) {
	parts := strings.Split(fieldsStr, ",")
	fields := make([]int, 0, len(parts))

	for _, part := range parts {
		field, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid field number: %s", part)
		}
		fields = append(fields, field)
	}

	return fields, nil
}
