package enc

import (
	"bytes"
	"fmt"
	"io"
)

type BulkString string

func (bs BulkString) Type() ValueType {
	return ValueTypeBulkString
}

func (bs BulkString) String() string {
	return string(bs)
}

func (bs BulkString) Encode(w io.Writer) error {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "$%d\r\n", len(bs))
	buf.WriteString(string(bs))

	_, err := w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("\r\n"))
	return err
}

func readBulkString(r io.Reader) (Value, error) {
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

	return BulkString(buf), nil
}
