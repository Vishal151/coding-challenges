package huffman

import (
	"io"
)

// CountFrequencies reads data from the given reader and returns a map of character frequencies
func CountFrequencies(r io.Reader) (map[byte]int, error) {
	frequencies := make(map[byte]int)
	buffer := make([]byte, 1024)

	for {
		n, err := r.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		for i := 0; i < n; i++ {
			frequencies[buffer[i]]++
		}

		if err == io.EOF {
			break
		}
	}

	return frequencies, nil
}
