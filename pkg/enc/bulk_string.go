package enc

import (
	"bytes"
	"fmt"
	"io"
)

type BulkString struct {
	Val  string
	Null bool
}

func (bs BulkString) Type() Type {
	return TypeBulkString
}

func (bs BulkString) String() string {
	if bs.Null {
		return "<null>"
	}
	return bs.Val
}

func (bs BulkString) Encode(w io.Writer) error {
	buf := bytes.NewBuffer(nil)

	length := len(bs.Val)
	if bs.Null {
		length = -1
	}

	fmt.Fprintf(buf, "$%d\r\n", length)
	if length >= 0 {
		buf.WriteString(bs.Val)
		buf.WriteString("\r\n")
	}

	_, err := w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return err
}

func decodeBulkString(r io.Reader) (Value, error) {
	length, err := readNumber(r)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	cr, err := readByte(r)
	if err != nil {
		return nil, err
	}

	lf, err := readByte(r)
	if err != nil {
		return nil, err
	}

	if cr != '\r' || lf != '\n' {
		return nil, fmt.Errorf("bad CRLF: %w", err)
	}

	return BulkString{Val: string(buf)}, nil
}
