package cmd

import (
	"fmt"
	"os"

	"github.com/vishal151/compression/internal/huffman"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "huffman",
	Short: "A Huffman coding implementation",
	Long:  `This is a Huffman coding implementation for the coding challenge.`,
}

var inputFile string

func init() {
	rootCmd.AddCommand(frequencyCmd)
	frequencyCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path")
	frequencyCmd.MarkFlagRequired("input")
}

var frequencyCmd = &cobra.Command{
	Use:   "frequency",
	Short: "Count character frequencies in the input file",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		frequencies, err := huffman.CountFrequencies(file)
		if err != nil {
			fmt.Printf("Error counting frequencies: %v\n", err)
			return
		}

		for char, count := range frequencies {
			fmt.Printf("%q: %d\n", char, count)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
