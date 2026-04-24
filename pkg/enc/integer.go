package enc

import (
	"bytes"
	"io"
	"strconv"
)

type Integer int

func (in Integer) Type() Type {
	return TypeInteger
}

func (in Integer) String() string {
	return strconv.Itoa(int(in))
}

func (in Integer) Encode(w io.Writer) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(int(in)))
	buf.WriteString("\r\n")

	_, err := w.Write(buf.Bytes())
	return err
}

func decodeInteger(r io.Reader) (Value, error) {
	num, err := readNumber(r)
	if err != nil {
		return nil, err
	}

	return Integer(num), nil
}
