package huffman

import (
	"container/heap"
	"encoding/binary"
	"io"
)

// Node represents a node in the Huffman tree
type Node struct {
	Char  byte
	Freq  int
	Left  *Node
	Right *Node
}

// PriorityQueue implements heap.Interface and holds Nodes
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Freq < pq[j].Freq }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// BuildHuffmanTree constructs a Huffman tree from the given frequency map
func BuildHuffmanTree(freqs map[byte]int) *Node {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Create leaf nodes for each character and add to the priority queue
	for char, freq := range freqs {
		heap.Push(&pq, &Node{Char: char, Freq: freq})
	}

	// Build the tree by combining nodes
	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*Node)
		right := heap.Pop(&pq).(*Node)
		parent := &Node{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(&pq, parent)
	}

	// Return the root of the Huffman tree
	return heap.Pop(&pq).(*Node)
}

// HuffmanCode represents the mapping of characters to their Huffman codes
type HuffmanCode map[byte]string

// GenerateHuffmanCodes creates a mapping of characters to their Huffman codes
func GenerateHuffmanCodes(root *Node) HuffmanCode {
	codes := make(HuffmanCode)
	generateCodesRecursive(root, "", codes)
	return codes
}

func generateCodesRecursive(node *Node, currentCode string, codes HuffmanCode) {
	if node == nil {
		return
	}

	// If it's a leaf node, assign the code to the character
	if node.Left == nil && node.Right == nil {
		codes[node.Char] = currentCode
		return
	}

	// Traverse left (add '0' to the code)
	generateCodesRecursive(node.Left, currentCode+"0", codes)

	// Traverse right (add '1' to the code)
	generateCodesRecursive(node.Right, currentCode+"1", codes)
}

// EncodeText encodes the input text using the generated Huffman codes
func EncodeText(input []byte, codes HuffmanCode) []byte {
	var encoded []byte
	var currentByte byte
	bitCount := 0

	for _, char := range input {
		code := codes[char]
		for _, bit := range code {
			if bit == '1' {
				currentByte |= 1 << (7 - bitCount)
			}
			bitCount++
			if bitCount == 8 {
				encoded = append(encoded, currentByte)
				currentByte = 0
				bitCount = 0
			}
		}
	}

	if bitCount > 0 {
		encoded = append(encoded, currentByte)
	}

	return encoded
}

// WriteTree writes the Huffman tree structure to the given writer
func WriteTree(node *Node, w io.Writer) error {
	if node == nil {
		return nil
	}

	if node.Left == nil && node.Right == nil {
		if err := binary.Write(w, binary.LittleEndian, uint8(1)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, node.Char)
	}

	if err := binary.Write(w, binary.LittleEndian, uint8(0)); err != nil {
		return err
	}

	if err := WriteTree(node.Left, w); err != nil {
		return err
	}
	return WriteTree(node.Right, w)
}

// WriteFrequencyTable writes the character frequency table to the given writer
func WriteFrequencyTable(freqs map[byte]int, w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, uint32(len(freqs))); err != nil {
		return err
	}

	for char, freq := range freqs {
		if err := binary.Write(w, binary.LittleEndian, char); err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, uint32(freq)); err != nil {
			return err
		}
	}

	return nil
}
