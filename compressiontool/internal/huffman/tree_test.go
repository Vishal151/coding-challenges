package huffman

import (
	"bytes"
	"encoding/binary"
	"reflect"
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

	readRoot, err := ReadTree(&buf)
	if err != nil {
		t.Fatalf("ReadTree returned an error: %v", err)
	}

	if !compareNodes(root, readRoot) {
		t.Errorf("Read tree doesn't match written tree")
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

func TestReadTree(t *testing.T) {
	// Create a sample tree structure
	treeData := []byte{
		1, 11, 0, 0, 0, 0, // Root node (freq: 11, not a leaf)
		1, 5, 0, 0, 0, 1, 'a', // Left child (freq: 5, leaf: 'a')
		1, 6, 0, 0, 0, 0, // Right child (freq: 6, not a leaf)
		1, 3, 0, 0, 0, 1, 'd', // Right-Left child (freq: 3, leaf: 'd')
		1, 3, 0, 0, 0, 0, // Right-Right child (freq: 3, not a leaf)
		1, 1, 0, 0, 0, 1, 'c', // Right-Right-Left child (freq: 1, leaf: 'c')
		1, 2, 0, 0, 0, 1, 'b', // Right-Right-Right child (freq: 2, leaf: 'b')
	}

	reader := bytes.NewReader(treeData)
	root, err := ReadTree(reader)
	if err != nil {
		t.Fatalf("ReadTree returned an error: %v", err)
	}

	// Verify the tree structure
	if root.Freq != 11 {
		t.Errorf("Root frequency incorrect. Got %d, want 11", root.Freq)
	}
	if root.Left.Char != 'a' || root.Left.Freq != 5 {
		t.Errorf("Left child incorrect. Got {%c, %d}, want {'a', 5}", root.Left.Char, root.Left.Freq)
	}
	if root.Right.Left.Char != 'd' || root.Right.Left.Freq != 3 {
		t.Errorf("Right-Left child incorrect. Got {%c, %d}, want {'d', 3}", root.Right.Left.Char, root.Right.Left.Freq)
	}
	if root.Right.Right.Left.Char != 'c' || root.Right.Right.Left.Freq != 1 {
		t.Errorf("Right-Right-Left child incorrect. Got {%c, %d}, want {'c', 1}", root.Right.Right.Left.Char, root.Right.Right.Left.Freq)
	}
	if root.Right.Right.Right.Char != 'b' || root.Right.Right.Right.Freq != 2 {
		t.Errorf("Right-Right-Right child incorrect. Got {%c, %d}, want {'b', 2}", root.Right.Right.Right.Char, root.Right.Right.Right.Freq)
	}
}

func TestReadFrequencyTable(t *testing.T) {
	tableData := &bytes.Buffer{}
	binary.Write(tableData, binary.LittleEndian, uint32(4)) // Number of entries
	binary.Write(tableData, binary.LittleEndian, byte('a'))
	binary.Write(tableData, binary.LittleEndian, uint32(5))
	binary.Write(tableData, binary.LittleEndian, byte('b'))
	binary.Write(tableData, binary.LittleEndian, uint32(2))
	binary.Write(tableData, binary.LittleEndian, byte('c'))
	binary.Write(tableData, binary.LittleEndian, uint32(1))
	binary.Write(tableData, binary.LittleEndian, byte('d'))
	binary.Write(tableData, binary.LittleEndian, uint32(3))

	freqs, err := ReadFrequencyTable(tableData)
	if err != nil {
		t.Fatalf("ReadFrequencyTable returned an error: %v", err)
	}

	expectedFreqs := map[byte]int{
		'a': 5,
		'b': 2,
		'c': 1,
		'd': 3,
	}

	if !reflect.DeepEqual(freqs, expectedFreqs) {
		t.Errorf("ReadFrequencyTable returned unexpected frequencies. Got %v, want %v", freqs, expectedFreqs)
	}
}

func TestDecodeTextWithSteps(t *testing.T) {
	root := &Node{
		Freq: 4,
		Left: &Node{Char: 'a', Freq: 1},
		Right: &Node{
			Freq: 3,
			Left: &Node{Char: 'b', Freq: 1},
			Right: &Node{
				Freq:  2,
				Left:  &Node{Char: 'c', Freq: 1},
				Right: &Node{Char: 'd', Freq: 1},
			},
		},
	}

	testCases := []struct {
		name     string
		input    []byte
		expected []byte
		wantErr  bool
		errMsg   string
	}{
		{"Valid input", []byte{0b01011011, 0b10000000}, []byte("abcd"), false, ""},
		{"Input too long", []byte{0b01011011, 0b10000001}, []byte("abcd"), false, ""}, // This should now pass without error
		{"Truly invalid bit sequence", []byte{0b11111111}, nil, true, "unexpected end of input"},
		{"Partial decode", []byte{0b10101010}, []byte("bbbb"), false, ""},
		{"Input too short", []byte{0b01011011}, nil, true, "unexpected end of input"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoded, steps, err := DecodeTextWithSteps(tc.input, root)

			t.Logf("Test case: %s", tc.name)
			for i, step := range steps {
				t.Logf("Step %d: %s", i+1, step)
			}

			if tc.wantErr {
				if err == nil {
					t.Errorf("DecodeTextWithSteps should have returned an error")
				} else if err.Error() != tc.errMsg {
					t.Errorf("DecodeTextWithSteps returned unexpected error. Got: %v, Want: %v", err, tc.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("DecodeTextWithSteps returned an unexpected error: %v", err)
				}
				if !bytes.Equal(decoded, tc.expected) {
					t.Errorf("DecodeTextWithSteps returned unexpected result. Got %v, want %v", decoded, tc.expected)
				}
			}
		})
	}
}

func TestFrequencyCounting(t *testing.T) {
	input := []byte("abracadabra")
	expected := map[byte]int{'a': 5, 'b': 2, 'r': 2, 'c': 1, 'd': 1}

	frequencies, err := CountFrequencies(bytes.NewReader(input))
	if err != nil {
		t.Fatalf("CountFrequencies returned an error: %v", err)
	}

	if !reflect.DeepEqual(frequencies, expected) {
		t.Errorf("Frequency count mismatch. Got %v, want %v", frequencies, expected)
	}
}

func TestBuildAndVerifyTree(t *testing.T) {
	freqs := map[byte]int{'a': 5, 'b': 2, 'c': 1, 'd': 3}
	root := BuildHuffmanTree(freqs)

	// Verify tree structure (you can expand this)
	if root.Freq != 11 {
		t.Errorf("Root frequency incorrect. Got %d, want 11", root.Freq)
	}
}

func TestGenerateAndVerifyCodes(t *testing.T) {
	root := &Node{
		Freq: 11,
		Left: &Node{Char: 'a', Freq: 5},
		Right: &Node{
			Freq: 6,
			Left: &Node{Char: 'd', Freq: 3},
			Right: &Node{
				Freq:  3,
				Left:  &Node{Char: 'c', Freq: 1},
				Right: &Node{Char: 'b', Freq: 2},
			},
		},
	}

	codes := GenerateHuffmanCodes(root)
	expected := HuffmanCode{'a': "0", 'd': "10", 'c': "110", 'b': "111"}

	if !reflect.DeepEqual(codes, expected) {
		t.Errorf("Generated codes mismatch. Got %v, want %v", codes, expected)
	}
}

func TestEncodeAndVerify(t *testing.T) {
	codes := HuffmanCode{'a': "0", 'b': "10", 'c': "110", 'd': "111"}
	input := []byte("abcd")
	expected := []byte{0b01011011, 0b10000000}

	encoded := EncodeText(input, codes)
	if !bytes.Equal(encoded, expected) {
		t.Errorf("Encoded output mismatch. Got %v, want %v", encoded, expected)
	}
}

func TestWriteAndReadTree(t *testing.T) {
	root := &Node{
		Freq: 11,
		Left: &Node{Char: 'a', Freq: 5},
		Right: &Node{
			Freq: 6,
			Left: &Node{Char: 'd', Freq: 3},
			Right: &Node{
				Freq:  3,
				Left:  &Node{Char: 'c', Freq: 1},
				Right: &Node{Char: 'b', Freq: 2},
			},
		},
	}

	var buf bytes.Buffer
	err := WriteTree(root, &buf)
	if err != nil {
		t.Fatalf("WriteTree returned an error: %v", err)
	}

	readRoot, err := ReadTree(&buf)
	if err != nil {
		t.Fatalf("ReadTree returned an error: %v", err)
	}

	if !compareNodes(root, readRoot) {
		t.Errorf("Read tree doesn't match written tree")
	}
}

func compareNodes(n1, n2 *Node) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}
	return n1.Freq == n2.Freq && n1.Char == n2.Char &&
		compareNodes(n1.Left, n2.Left) && compareNodes(n1.Right, n2.Right)
}

func TestFullEncodeDecode(t *testing.T) {
	input := []byte("abracadabra")
	var buf bytes.Buffer

	// Encode
	frequencies, _ := CountFrequencies(bytes.NewReader(input))
	root := BuildHuffmanTree(frequencies)
	codes := GenerateHuffmanCodes(root)
	encoded := EncodeText(input, codes)

	WriteTree(root, &buf)
	binary.Write(&buf, binary.LittleEndian, uint32(len(input)))
	buf.Write(encoded)

	// Decode
	readRoot, _ := ReadTree(&buf)
	var length uint32
	binary.Read(&buf, binary.LittleEndian, &length)
	decoded, err := DecodeText(buf.Bytes(), readRoot)

	if err != nil {
		t.Fatalf("DecodeText returned an error: %v", err)
	}

	if !bytes.Equal(input, decoded) {
		t.Errorf("Decoded output doesn't match input. Got %s, want %s", decoded, input)
	}
}
