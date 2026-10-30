package enc

import (
	"fmt"
	"io"
	"strings"
)

func NullArr() Array {
	return Array{len: -1}
}

func Arr(items ...Value) Array {
	return Array{len: len(items), items: items}
}

type Array struct {
	len   int
	items []Value
}

func (a Array) Len() int {
	return a.len
}

func (a Array) Items() []Value {
	return a.items
}

func (a Array) Type() Type {
	return TypeArray
}

func (a Array) String() string {
	if a.len < 0 {
		return "nil"
	}

	b := &strings.Builder{}

	b.WriteString("[")
	for i, elem := range a.items {
		b.WriteString(elem.String())

		if i != len(a.items)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")

	return b.String()
}

func (a Array) Encode(w io.Writer) error {
	_, err := fmt.Fprintf(w, "*%d\r\n", a.len)
	if err != nil {
		return err
	}

	for _, el := range a.items {
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
		return Array{}, fmt.Errorf("read array length: %w", err)
	}

	if num < 0 {
		return NullArr(), nil
	}

	array := make([]Value, 0, num)
	for i := range num {
		val, err := Decode(r)
		if err != nil {
			return Array{}, fmt.Errorf("read elem %d in array: %w", i, err)
		}

		array = append(array, val)
	}

	return Array{len: len(array), items: array}, nil
}
