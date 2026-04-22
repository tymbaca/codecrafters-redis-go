package enc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArray(t *testing.T) {
	a := Array{
		BulkString("hello"),
		BulkString("world"),
	}

	buf := bytes.NewBuffer(nil)

	err := a.Encode(buf)
	require.NoError(t, err)

	require.Equal(t, []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"), buf.Bytes())
}
