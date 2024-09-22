package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type jsonState int

const (
	expectValue jsonState = iota
	expectKey
	expectColon
	expectComma
)

type jsonContext struct {
	state    jsonState
	inString bool
	stack    []rune
}

type jsonReader struct {
	input string
	pos   int
}

func newJSONReader(input string) *jsonReader {
	return &jsonReader{input: input, pos: 0}
}

func (r *jsonReader) peek() (rune, bool) {
	if r.pos >= len(r.input) {
		return 0, false
	}
	return rune(r.input[r.pos]), true
}

func (r *jsonReader) next() (rune, bool) {
	if r.pos >= len(r.input) {
		return 0, false
	}
	ch := rune(r.input[r.pos])
	r.pos++
	return ch, true
}

func ValidateJSONStep1(data []byte) (bool, error) {
	input := strings.TrimSpace(string(data))

	if len(input) == 0 {
		return false, errors.New("empty input")
	}

	if !isValidJSONStart(input) {
		return false, errors.New("JSON must start with '{' or '['")
	}

	return true, nil
}

func ValidateJSONStep2(data []byte) (bool, error) {
	input := strings.TrimSpace(string(data))

	if valid, err := ValidateJSONStep1(data); !valid {
		return false, err
	}

	ctx := &jsonContext{
		state: expectValue,
		stack: []rune{},
	}

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])
		if err := processCharStep2(ctx, ch); err != nil {
			return false, err
		}
	}

	if len(ctx.stack) != 0 {
		return false, errors.New("unbalanced braces or brackets")
	}

	if ctx.state != expectComma && ctx.state != expectValue {
		return false, errors.New("incomplete JSON structure")
	}

	return true, nil
}

func ValidateJSONStep3(data []byte) (bool, error) {
	input := strings.TrimSpace(string(data))

	ctx := &jsonContext{
		state: expectValue,
		stack: []rune{},
	}
	reader := newJSONReader(input)

	for {
		ch, ok := reader.next()
		if !ok {
			break
		}
		if err := processCharStep3(ctx, ch, reader); err != nil {
			return false, err
		}
	}

	if len(ctx.stack) != 0 {
		return false, errors.New("unbalanced braces or brackets")
	}

	if ctx.state != expectComma && ctx.state != expectValue {
		return false, errors.New("incomplete JSON structure")
	}

	return true, nil
}

func ValidateJSONStep4(data []byte) (bool, error) {
	input := strings.TrimSpace(string(data))

	if valid, err := ValidateJSONStep3(data); !valid {
		return false, err
	}

	ctx := &jsonContext{
		state: expectValue,
		stack: []rune{},
	}
	reader := newJSONReader(input)

	for {
		ch, ok := reader.next()
		if !ok {
			break
		}
		if err := processCharStep4(ctx, ch, reader); err != nil {
			return false, err
		}
	}

	return true, nil
}

func isValidJSONStart(input string) bool {
	return input[0] == '{' || input[0] == '['
}

func processCharStep2(ctx *jsonContext, ch rune) error {
	switch ch {
	case '{', '[':
		ctx.stack = append(ctx.stack, ch)
		ctx.state = expectKey
	case '}', ']':
		if len(ctx.stack) == 0 || (ch == '}' && ctx.stack[len(ctx.stack)-1] != '{') || (ch == ']' && ctx.stack[len(ctx.stack)-1] != '[') {
			return errors.New("unbalanced braces or brackets")
		}
		ctx.stack = ctx.stack[:len(ctx.stack)-1]
		if len(ctx.stack) > 0 {
			ctx.state = expectComma
		}
	case '"':
		ctx.inString = !ctx.inString
		if !ctx.inString {
			if ctx.state == expectKey {
				ctx.state = expectColon
			} else if ctx.state == expectValue {
				ctx.state = expectComma
			}
		}
	case ':':
		if ctx.inString {
			return nil
		}
		if ctx.state != expectColon {
			return errors.New("unexpected colon")
		}
		ctx.state = expectValue
	case ',':
		if ctx.inString {
			return nil
		}
		if ctx.state != expectComma {
			return errors.New("unexpected comma")
		}
		if len(ctx.stack) == 0 {
			return errors.New("unexpected comma")
		}
		if ctx.stack[len(ctx.stack)-1] == '{' {
			ctx.state = expectKey
		} else {
			ctx.state = expectValue
		}
	default:
		if !ctx.inString && !unicode.IsSpace(ch) {
			if ctx.state != expectValue {
				return fmt.Errorf("unexpected character: %c", ch)
			}
			ctx.state = expectComma
		}
	}
	return nil
}

func processCharStep3(ctx *jsonContext, ch rune, reader *jsonReader) error {
	if ctx.inString {
		return processCharStep2(ctx, ch)
	}

	switch ch {
	case 't', 'f', 'n':
		return processConstant(ctx, ch, reader)
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return processNumber(ctx, reader)
	default:
		return processCharStep2(ctx, ch)
	}
}

func processConstant(ctx *jsonContext, ch rune, reader *jsonReader) error {
	if ctx.state != expectValue {
		return fmt.Errorf("unexpected constant")
	}

	var constant string
	switch ch {
	case 't':
		constant = "true"
	case 'f':
		constant = "false"
	case 'n':
		constant = "null"
	}

	for i := 1; i < len(constant); i++ {
		nextCh, ok := reader.next()
		if !ok || rune(constant[i]) != nextCh {
			return fmt.Errorf("invalid constant %s", constant)
		}
	}

	ctx.state = expectComma
	return nil
}

func processNumber(ctx *jsonContext, reader *jsonReader) error {
	if ctx.state != expectValue {
		return fmt.Errorf("unexpected number")
	}

	// Simplified number validation
	for {
		ch, ok := reader.peek()
		if !ok || !(unicode.IsDigit(ch) || ch == '.' || ch == 'e' || ch == 'E' || ch == '+' || ch == '-') {
			break
		}
		reader.next()
	}

	ctx.state = expectComma
	return nil
}

func processCharStep4(ctx *jsonContext, ch rune, reader *jsonReader) error {
	return processCharStep3(ctx, ch, reader)
}
