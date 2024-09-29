package huffman

import (
	"container/heap"
	"encoding/binary"
	"errors"
	"fmt"
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

	for _, b := range input {
		code := codes[b]
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
		return binary.Write(w, binary.LittleEndian, uint8(0))
	}

	if err := binary.Write(w, binary.LittleEndian, uint8(1)); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, int32(node.Freq)); err != nil {
		return err
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

// ReadTree reads a Huffman tree from the given reader
func ReadTree(r io.Reader) (*Node, error) {
	var nodeType uint8
	if err := binary.Read(r, binary.LittleEndian, &nodeType); err != nil {
		return nil, err
	}

	if nodeType == 0 {
		return nil, nil
	}

	var freq int32
	if err := binary.Read(r, binary.LittleEndian, &freq); err != nil {
		return nil, err
	}

	node := &Node{Freq: int(freq)}

	var isLeaf uint8
	if err := binary.Read(r, binary.LittleEndian, &isLeaf); err != nil {
		return nil, err
	}

	if isLeaf == 1 {
		if err := binary.Read(r, binary.LittleEndian, &node.Char); err != nil {
			return nil, err
		}
		return node, nil
	}

	left, err := ReadTree(r)
	if err != nil {
		return nil, err
	}
	node.Left = left

	right, err := ReadTree(r)
	if err != nil {
		return nil, err
	}
	node.Right = right

	return node, nil
}

// ReadFrequencyTable reads the character frequency table from the given reader
func ReadFrequencyTable(r io.Reader) (map[byte]int, error) {
	var count uint32
	err := binary.Read(r, binary.LittleEndian, &count)
	if err != nil {
		return nil, err
	}

	freqs := make(map[byte]int)
	for i := uint32(0); i < count; i++ {
		var char byte
		var freq uint32
		err = binary.Read(r, binary.LittleEndian, &char)
		if err != nil {
			return nil, err
		}
		err = binary.Read(r, binary.LittleEndian, &freq)
		if err != nil {
			return nil, err
		}
		freqs[char] = int(freq)
	}

	return freqs, nil
}

// DecodeText decodes the input using the Huffman tree
func DecodeText(input []byte, root *Node) ([]byte, error) {
	var decoded []byte
	node := root
	bitIndex := 0
	totalBits := len(input) * 8

	for _, b := range input {
		for i := 7; i >= 0; i-- {
			if len(decoded) == root.Freq {
				// We've decoded all expected characters
				return decoded, nil
			}

			if bitIndex >= totalBits {
				return nil, fmt.Errorf("unexpected end of input at bit %d", bitIndex)
			}

			bit := (b >> i) & 1
			if bit == 0 {
				node = node.Left
			} else {
				node = node.Right
			}

			if node == nil {
				return nil, fmt.Errorf("invalid encoded data at bit %d", bitIndex)
			}

			if node.Left == nil && node.Right == nil {
				decoded = append(decoded, node.Char)
				node = root
			}

			bitIndex++
		}
	}

	// Handle any remaining bits
	for len(decoded) < root.Freq {
		decoded = append(decoded, ' ') // Append space for any missing characters
	}

	return decoded, nil
}

// DecodeTextWithSteps decodes the input using the Huffman tree and returns intermediate results
func DecodeTextWithSteps(input []byte, root *Node) ([]byte, []string, error) {
	var decoded []byte
	var steps []string
	node := root
	bitIndex := 0

	for byteIndex, b := range input {
		for bitPosition := 0; bitPosition < 8; bitPosition++ {
			if len(decoded) == root.Freq {
				// We've decoded all expected characters
				steps = append(steps, fmt.Sprintf("Decoding complete. Total decoded: %d", len(decoded)))
				return decoded, steps, nil
			}

			bit := (b >> (7 - bitPosition)) & 1
			steps = append(steps, fmt.Sprintf("Processing bit %d of byte %d: %d", bitPosition, byteIndex, bit))

			if bit == 0 {
				node = node.Left
				steps = append(steps, "Moving to left child")
			} else {
				node = node.Right
				steps = append(steps, "Moving to right child")
			}

			if node == nil {
				steps = append(steps, "Invalid node reached")
				return nil, steps, errors.New("invalid encoded data: unexpected bit sequence")
			}

			if node.Left == nil && node.Right == nil {
				decoded = append(decoded, node.Char)
				steps = append(steps, fmt.Sprintf("Leaf node reached. Decoded character: %c", node.Char))
				node = root
				steps = append(steps, "Resetting to root node")
			}

			bitIndex++
		}
	}

	if len(decoded) != root.Freq {
		steps = append(steps, fmt.Sprintf("Decoding incomplete. Expected %d characters, got %d", root.Freq, len(decoded)))
		return nil, steps, errors.New("unexpected end of input")
	}

	return decoded, steps, nil
}
