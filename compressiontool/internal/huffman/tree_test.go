package huffman

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestBuildHuffmanTree(t *testing.T) {
	freqs := map[byte]int{
		'a': 5,
		'b': 2,
		'c': 1,
		'd': 3,
	}

	root := BuildHuffmanTree(freqs)

	// Check if the root frequency is the sum of all frequencies
	if root.Freq != 11 {
		t.Errorf("Expected root frequency to be 11, got %d", root.Freq)
	}

	// Check the structure and frequencies of the tree
	if root.Left.Freq != 5 || root.Right.Freq != 6 {
		t.Errorf("Unexpected root children frequencies. Left: %d, Right: %d", root.Left.Freq, root.Right.Freq)
	}

	// Check the left subtree (should be a leaf node with 'a')
	if root.Left.Char != 'a' || root.Left.Left != nil || root.Left.Right != nil {
		t.Errorf("Unexpected left child. Char: %c, Left: %v, Right: %v", root.Left.Char, root.Left.Left, root.Left.Right)
	}

	// Check the right subtree
	rightChild := root.Right
	if rightChild.Left.Freq != 3 || rightChild.Right.Freq != 3 {
		t.Errorf("Unexpected right subtree structure. Left freq: %d, Right freq: %d", rightChild.Left.Freq, rightChild.Right.Freq)
	}

	// Check the leaves of the right subtree
	if rightChild.Left.Char != 'd' || rightChild.Left.Left != nil || rightChild.Left.Right != nil {
		t.Errorf("Unexpected 'd' node. Char: %c, Left: %v, Right: %v", rightChild.Left.Char, rightChild.Left.Left, rightChild.Left.Right)
	}

	lastNode := rightChild.Right
	if lastNode.Left.Char != 'c' || lastNode.Left.Freq != 1 || lastNode.Right.Char != 'b' || lastNode.Right.Freq != 2 {
		t.Errorf("Unexpected 'b' and 'c' nodes. Left: {%c, %d}, Right: {%c, %d}",
			lastNode.Left.Char, lastNode.Left.Freq, lastNode.Right.Char, lastNode.Right.Freq)
	}
}

func TestGenerateHuffmanCodes(t *testing.T) {
	freqs := map[byte]int{
		'a': 5,
		'b': 2,
		'c': 1,
		'd': 3,
	}

	root := BuildHuffmanTree(freqs)
	codes := GenerateHuffmanCodes(root)

	expectedCodes := HuffmanCode{
		'a': "0",
		'b': "111",
		'c': "110",
		'd': "10",
	}

	if len(codes) != len(expectedCodes) {
		t.Errorf("Expected %d codes, got %d", len(expectedCodes), len(codes))
	}

	for char, code := range expectedCodes {
		if codes[char] != code {
			t.Errorf("For character '%c', expected code '%s', got '%s'", char, code, codes[char])
		}
	}
}

func TestEncodeText(t *testing.T) {
	codes := HuffmanCode{
		'a': "0",
		'b': "10",
		'c': "110",
		'd': "111",
	}

	input := []byte("abcd")
	expected := []byte{0b01011011, 0b10000000}

	encoded := EncodeText(input, codes)

	if !bytes.Equal(encoded, expected) {
		t.Errorf("Expected encoded bytes %v, got %v", expected, encoded)
	}
}

// Add more tests for WriteTree and WriteFrequencyTable

func TestWriteTree(t *testing.T) {
	root := &Node{
		Freq: 11,
		Left: &Node{
			Char: 'a',
			Freq: 5,
		},
		Right: &Node{
			Freq: 6,
			Left: &Node{
				Char: 'd',
				Freq: 3,
			},
			Right: &Node{
				Freq: 3,
				Left: &Node{
					Char: 'c',
					Freq: 1,
				},
				Right: &Node{
					Char: 'b',
					Freq: 2,
				},
			},
		},
	}

	var buf bytes.Buffer
	err := WriteTree(root, &buf)
	if err != nil {
		t.Fatalf("WriteTree returned an error: %v", err)
	}

	expected := []byte{
		0,      // Internal node
		1, 'a', // Leaf node 'a'
		0,      // Internal node
		1, 'd', // Leaf node 'd'
		0,      // Internal node
		1, 'c', // Leaf node 'c'
		1, 'b', // Leaf node 'b'
	}

	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("WriteTree output doesn't match expected. Got %v, want %v", buf.Bytes(), expected)
	}
}

func TestWriteFrequencyTable(t *testing.T) {
	freqs := map[byte]int{
		'a': 5,
		'b': 2,
		'c': 1,
		'd': 3,
	}

	var buf bytes.Buffer
	err := WriteFrequencyTable(freqs, &buf)
	if err != nil {
		t.Fatalf("WriteFrequencyTable returned an error: %v", err)
	}

	// Read the number of entries
	var count uint32
	err = binary.Read(&buf, binary.LittleEndian, &count)
	if err != nil {
		t.Fatalf("Error reading count: %v", err)
	}
	if count != 4 {
		t.Errorf("Expected 4 entries, got %d", count)
	}

	// Read and check each entry
	expectedFreqs := map[byte]uint32{
		'a': 5,
		'b': 2,
		'c': 1,
		'd': 3,
	}
	for i := 0; i < int(count); i++ {
		var char byte
		var freq uint32
		err = binary.Read(&buf, binary.LittleEndian, &char)
		if err != nil {
			t.Fatalf("Error reading char: %v", err)
		}
		err = binary.Read(&buf, binary.LittleEndian, &freq)
		if err != nil {
			t.Fatalf("Error reading freq: %v", err)
		}
		if expectedFreq, ok := expectedFreqs[char]; !ok || expectedFreq != freq {
			t.Errorf("Unexpected frequency for char %c: got %d, want %d", char, freq, expectedFreq)
		}
		delete(expectedFreqs, char)
	}

	if len(expectedFreqs) != 0 {
		t.Errorf("Not all expected frequencies were found in the output")
	}
}
