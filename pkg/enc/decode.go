// Package enc provides RESP types and encoding/decoding functions.
package enc

import (
	"fmt"
	"io"
	"strconv"
	"unicode"
)

func Decode(r io.Reader) (Value, error) {
	head, err := readByte(r)
	if err != nil {
		return nil, fmt.Errorf("read value header: %w", err)
	}

	switch head {
	case '*':
		return decodeArray(r)
	case '$':
		return decodeBulkString(r)
	case '+':
		return decodeSimpleString(r)
	case '-':
		return decodeSimpleError(r)
	case ':':
		return decodeInteger(r)
	}

	return nil, fmt.Errorf("unsupported value header: %c", head)
}

func readNumber(r io.Reader) (int, error) {
	first, err := readByte(r)
	if err != nil {
		return 0, err
	}

	var digits []byte

	sign := 1
	if first == '-' {
		sign = -1
	} else if first == '+' {
		sign = +1
	} else if unicode.IsDigit(rune(first)) {
		digits = append(digits, first)
	} else {
		return 0, fmt.Errorf("invalid character '%c', expected number", first)
	}

	for {
		b, err := readByte(r)
		if err != nil {
			return 0, fmt.Errorf("read number: %w", err)
		}

		if unicode.IsDigit(rune(b)) {
			digits = append(digits, b)
		} else if b == '\r' {
			err := finishCRLF(r)
			if err != nil {
				return 0, err
			}

			num, err := strconv.Atoi(string(digits))
			if err != nil {
				return 0, fmt.Errorf("convert digits to number: %w", err)
			}

			return sign * num, nil
		} else {
			return 0, fmt.Errorf("invalid character '%c', expected number", first)
		}
	}
}

func finishCRLF(r io.Reader) error {
	b, err := readByte(r)
	if err != nil {
		return fmt.Errorf("read CRLF: %w", err)
	}

	if b != '\n' {
		return fmt.Errorf("read CRLF: no \\n character")
	}

	return nil
}

func readByte(r io.Reader) (byte, error) {
	var buf [1]byte
	_, err := io.ReadFull(r, buf[:])
	return buf[0], err
}
