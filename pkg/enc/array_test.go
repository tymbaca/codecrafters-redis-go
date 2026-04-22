package enc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArray(t *testing.T) {
	a := Array{
		BulkString("hello"),
		SimpleString("world"),
		Array{
			BulkString("foo"),
		},
	}

	buf := bytes.NewBuffer(nil)

	err := a.Encode(buf)
	require.NoError(t, err)

	require.Equal(t, []byte("*3\r\n$5\r\nhello\r\n+world\r\n*1\r\n$3\r\nfoo\r\n"), buf.Bytes())
}
