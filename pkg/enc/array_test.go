package enc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArray(t *testing.T) {
	a := Array{
		BulkString{Val: "hello"},
		SimpleString("world"),
		Array{
			BulkString{Val: "foo"},
		},
	}

	buf := bytes.NewBuffer(nil)

	err := a.Encode(buf)
	require.NoError(t, err)

	require.Equal(t, []byte("*3\r\n$5\r\nhello\r\n+world\r\n*1\r\n$3\r\nfoo\r\n"), buf.Bytes())

	parsedVal, err := ReadValue(buf)
	require.NoError(t, err)
	require.Equal(t, Array{
		BulkString{Val: "hello"},
		SimpleString("world"),
		Array{
			BulkString{Val: "foo"},
		},
	}, parsedVal)
}
