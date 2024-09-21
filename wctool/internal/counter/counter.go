package counter

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

type Counts struct {
	Bytes int64
	Lines int
	Words int
	Chars int
}

func Count(r io.Reader) (Counts, error) {
	var counts Counts
	var err error

	// Read all content into a buffer
	content, err := io.ReadAll(r)
	if err != nil {
		return counts, err
	}

	// Count bytes
	counts.Bytes = int64(len(content))

	// Count lines, words, and chars
	counts.Lines, counts.Words, counts.Chars = countLinesWordsChars(content)

	return counts, nil
}

func countLinesWordsChars(content []byte) (lines, words, chars int) {
	lines = bytes.Count(content, []byte{'\n'})
	if len(content) > 0 && content[len(content)-1] != '\n' {
		lines++
	}

	chars = utf8.RuneCount(content)

	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Split(scanWords)
	for scanner.Scan() {
		words++
	}

	return
}

func CountLines(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	return lineCount, scanner.Err()
}

func CountWords(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanWords)

	wordCount := 0
	for scanner.Scan() {
		wordCount++
	}

	return wordCount, scanner.Err()
}

func CountChars(r io.Reader) (int, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	return utf8.RuneCount(content), nil
}

func scanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces and punctuation
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			break
		}
	}

	// Scan until space, newline, or punctuation
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			if i > start {
				return i + width, data[start:i], nil
			}
			start = i + width
		}
	}

	// If we're at EOF, we have a final, non-empty token. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// Request more data.
	return start, nil, nil
}
