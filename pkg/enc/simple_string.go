package enc

import (
	"bytes"
	"io"
)

type SimpleString string

func (ss SimpleString) Type() ValueType {
	return ValueTypeSimpleString
}

func (ss SimpleString) String() string {
	return string(ss)
}

func (ss SimpleString) Encode(w io.Writer) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("+")
	buf.WriteString(string(ss))
	buf.WriteString("\r\n")

	_, err := w.Write(buf.Bytes())
	return err
}

func readSimpleString(r io.Reader) (Value, error) {
	panic("unimplemented")
}
