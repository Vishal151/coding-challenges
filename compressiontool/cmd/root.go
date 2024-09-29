package cmd

import (
	"bytes"
	"fmt"
	"log" // Add this import
	"os"
	"strings"

	"github.com/vishal151/compression/internal/huffman"

	"encoding/binary"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "huffman",
	Short: "A Huffman coding implementation",
	Long:  `This is a Huffman coding implementation for the coding challenge.`,
}

var inputFile string
var outputFile string
var useTree bool

func init() {
	rootCmd.AddCommand(frequencyCmd)
	rootCmd.AddCommand(treeCmd)
	rootCmd.AddCommand(codesCmd)
	rootCmd.AddCommand(encodeCmd)
	rootCmd.AddCommand(decodeCmd)
	frequencyCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path")
	frequencyCmd.MarkFlagRequired("input")
	treeCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path")
	treeCmd.MarkFlagRequired("input")
	codesCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path")
	codesCmd.MarkFlagRequired("input")
	encodeCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path")
	encodeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	encodeCmd.Flags().BoolVarP(&useTree, "use-tree", "t", false, "Use tree structure instead of frequency table")
	encodeCmd.MarkFlagRequired("input")
	encodeCmd.MarkFlagRequired("output")
	decodeCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path")
	decodeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	decodeCmd.Flags().BoolVarP(&useTree, "use-tree", "t", false, "Use tree structure instead of frequency table")
	decodeCmd.MarkFlagRequired("input")
	decodeCmd.MarkFlagRequired("output")
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

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Build and display the Huffman tree for the input file",
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

		root := huffman.BuildHuffmanTree(frequencies)
		fmt.Println("Huffman Tree built successfully.")
		fmt.Printf("Root frequency: %d\n", root.Freq)
		fmt.Println("Tree structure (simplified):")
		printTree(root, 0)
	},
}

var codesCmd = &cobra.Command{
	Use:   "codes",
	Short: "Generate and display Huffman codes for the input file",
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

		root := huffman.BuildHuffmanTree(frequencies)
		codes := huffman.GenerateHuffmanCodes(root)

		fmt.Println("Huffman Codes:")
		for char, code := range codes {
			fmt.Printf("%q: %s\n", char, code)
		}
	},
}

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode the input file using Huffman coding",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatalf("Error reading input file: %v", err)
		}

		frequencies, err := huffman.CountFrequencies(bytes.NewReader(input))
		if err != nil {
			log.Fatalf("Error counting frequencies: %v", err)
		}

		root := huffman.BuildHuffmanTree(frequencies)
		codes := huffman.GenerateHuffmanCodes(root)
		encoded := huffman.EncodeText(input, codes)

		outputFile, err := os.Create(outputFile)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer outputFile.Close()

		if err := binary.Write(outputFile, binary.LittleEndian, useTree); err != nil {
			log.Fatalf("Error writing use-tree flag: %v", err)
		}

		if useTree {
			if err := huffman.WriteTree(root, outputFile); err != nil {
				log.Fatalf("Error writing tree: %v", err)
			}
		} else {
			if err := huffman.WriteFrequencyTable(frequencies, outputFile); err != nil {
				log.Fatalf("Error writing frequency table: %v", err)
			}
		}

		if err := binary.Write(outputFile, binary.LittleEndian, uint32(root.Freq)); err != nil {
			log.Fatalf("Error writing total frequency: %v", err)
		}

		if _, err := outputFile.Write(encoded); err != nil {
			log.Fatalf("Error writing encoded data: %v", err)
		}

		fmt.Printf("File encoded successfully. Output written to %s\n", outputFile.Name())
		fmt.Printf("Original size: %d bytes\n", len(input))
		fmt.Printf("Compressed size: %d bytes\n", len(encoded))
		fmt.Printf("Compression ratio: %.2f%%\n", float64(len(encoded))/float64(len(input))*100)
	},
}

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a Huffman-encoded file",
	Run: func(cmd *cobra.Command, args []string) {
		inputData, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatalf("Error reading input file: %v", err)
		}

		reader := bytes.NewReader(inputData)

		var usedTree bool
		if err := binary.Read(reader, binary.LittleEndian, &usedTree); err != nil {
			log.Fatalf("Error reading use-tree flag: %v", err)
		}

		var root *huffman.Node
		if usedTree {
			root, err = huffman.ReadTree(reader)
			if err != nil {
				log.Fatalf("Error reading Huffman tree: %v", err)
			}
		} else {
			frequencies, err := huffman.ReadFrequencyTable(reader)
			if err != nil {
				log.Fatalf("Error reading frequency table: %v", err)
			}
			root = huffman.BuildHuffmanTree(frequencies)
		}

		var totalFreq uint32
		if err := binary.Read(reader, binary.LittleEndian, &totalFreq); err != nil {
			log.Fatalf("Error reading total frequency: %v", err)
		}
		root.Freq = int(totalFreq)

		encodedData := inputData[len(inputData)-reader.Len():]
		decoded, err := huffman.DecodeText(encodedData, root)
		if err != nil {
			log.Fatalf("Error decoding text: %v", err)
		}

		// Preserve line endings
		decoded = bytes.ReplaceAll(decoded, []byte{'\r', '\n'}, []byte{'\n'})

		err = os.WriteFile(outputFile, decoded, 0644)
		if err != nil {
			log.Fatalf("Error writing output file: %v", err)
		}

		fmt.Printf("File decoded successfully. Output written to %s\n", outputFile)
		fmt.Printf("Compressed size: %d bytes\n", len(inputData))
		fmt.Printf("Decompressed size: %d bytes\n", len(decoded))
	},
}

func printTree(node *huffman.Node, level int) {
	if node == nil {
		return
	}

	indent := strings.Repeat("  ", level)
	if node.Char != 0 {
		fmt.Printf("%s'%c' (%d)\n", indent, node.Char, node.Freq)
	} else {
		fmt.Printf("%sInternal Node (%d)\n", indent, node.Freq)
	}

	printTree(node.Left, level+1)
	printTree(node.Right, level+1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
