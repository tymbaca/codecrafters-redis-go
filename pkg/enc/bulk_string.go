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
	panic("unimplemented")
}
