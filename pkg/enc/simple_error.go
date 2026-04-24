package enc

import (
	"bytes"
	"io"
)

type SimpleError string

func (ss SimpleError) Type() Type {
	return TypeSimpleError
}

func (ss SimpleError) String() string {
	return string(ss)
}

func (ss SimpleError) Encode(w io.Writer) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("-")
	buf.WriteString(string(ss))
	buf.WriteString("\r\n")

	_, err := w.Write(buf.Bytes())
	return err
}

func decodeSimpleError(r io.Reader) (Value, error) {
	buf := bytes.NewBuffer(nil)

	for {
		b, err := readByte(r)
		if err != nil {
			return nil, err
		}

		if b == '\r' {
			return SimpleError(buf.String()), finishCRLF(r)
		}

		buf.WriteByte(b)
	}
}
