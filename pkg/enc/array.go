package enc

import (
	"fmt"
	"io"
	"strings"
)

type Array []Value

func (a Array) Type() Type {
	return TypeArray
}

func (a Array) String() string {
	b := &strings.Builder{}

	b.WriteString("[")
	for i, elem := range a {
		b.WriteString(elem.String())

		if i != len(a)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")

	return b.String()
}

func (a Array) Encode(w io.Writer) error {
	_, err := fmt.Fprintf(w, "*%d\r\n", len(a))
	if err != nil {
		return err
	}

	for _, el := range a {
		err := el.Encode(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func decodeArray(r io.Reader) (Array, error) {
	num, err := readNumber(r)
	if err != nil {
		return nil, fmt.Errorf("read array length: %w", err)
	}

	array := make(Array, 0, num)
	for i := range num {
		val, err := Decode(r)
		if err != nil {
			return nil, fmt.Errorf("read elem %d in array: %w", i, err)
		}

		array = append(array, val)
	}

	return array, nil
}
