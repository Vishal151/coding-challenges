package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/vishal151/cut-tool/internal/cutter"
)

var fields string
var bytes string
var delimiter string
var onlyDelimited bool

var rootCmd = &cobra.Command{
	Use:   "cut-tool [file]",
	Short: "A cut tool implementation",
	Long:  `A cut tool implementation for the coding challenge at https://codingchallenges.fyi/challenges/challenge-cut`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var input io.Reader
		var err error

		if len(args) == 0 {
			input = os.Stdin
		} else {
			input, err = cutter.ReadFile(args[0])
			if err != nil {
				fmt.Println("Error reading file:", err)
				os.Exit(1)
			}
		}

		var result string
		if fields != "" {
			result, err = cutter.CutByFields(input, fields, delimiter, onlyDelimited)
		} else if bytes != "" {
			result, err = cutter.CutByBytes(input, bytes)
		} else {
			fmt.Println("Error: either -f or -b flag must be specified")
			os.Exit(1)
		}

		if err != nil {
			fmt.Println("Error cutting content:", err)
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
	rootCmd.Flags().StringVarP(&fields, "fields", "f", "", "select only these fields; also print any line that contains no delimiter character, unless the -s option is specified")
	rootCmd.Flags().StringVarP(&bytes, "bytes", "b", "", "select only these bytes")
	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "\t", "use DELIM instead of TAB for field delimiter")
	rootCmd.Flags().BoolVarP(&onlyDelimited, "only-delimited", "s", false, "do not print lines not containing delimiters")
}
