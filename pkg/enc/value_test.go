package enc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		a := Array{
			BulkString{Val: "hello"},
			SimpleString("world"),
			SimpleError("bad"),
			Array{
				Integer(77),
				Integer(-77),
				BulkString{Val: "foo"},
			},
		}

		buf := bytes.NewBuffer(nil)

		err := a.Encode(buf)
		require.NoError(t, err)

		require.Equal(t, []byte("*4\r\n$5\r\nhello\r\n+world\r\n-bad\r\n*3\r\n:77\r\n:-77\r\n$3\r\nfoo\r\n"), buf.Bytes())

		parsedVal, err := Decode(buf)
		require.NoError(t, err)
		require.Equal(t, Array{
			BulkString{Val: "hello"},
			SimpleString("world"),
			SimpleError("bad"),
			Array{
				Integer(77),
				Integer(-77),
				BulkString{Val: "foo"},
			},
		}, parsedVal)
	})
	t.Run("empty array", func(t *testing.T) {
		a := Array{}

		buf := bytes.NewBuffer(nil)

		err := a.Encode(buf)
		require.NoError(t, err)

		require.Equal(t, []byte("*0\r\n"), buf.Bytes())

		parsedVal, err := Decode(buf)
		require.NoError(t, err)
		require.Equal(t, Array{}, parsedVal)
	})
}

func Test_parseNumber(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("456\r\n")
	num, err := readNumber(buf)
	require.NoError(t, err)
	require.Equal(t, 456, num)

	buf.WriteString("+456\r\n")
	num, err = readNumber(buf)
	require.NoError(t, err)
	require.Equal(t, 456, num)

	buf.WriteString("-456\r\n")
	num, err = readNumber(buf)
	require.NoError(t, err)
	require.Equal(t, -456, num)

	buf.WriteString("-4\r\n")
	num, err = readNumber(buf)
	require.NoError(t, err)
	require.Equal(t, -4, num)

	buf.WriteString("4\r\n")
	num, err = readNumber(buf)
	require.NoError(t, err)
	require.Equal(t, 4, num)

	buf.WriteString("0\r\n")
	num, err = readNumber(buf)
	require.NoError(t, err)
	require.Equal(t, 0, num)

	buf = bytes.NewBufferString("--\r\n")
	_, err = readNumber(buf)
	require.Error(t, err)

	buf = bytes.NewBufferString("af2f\r\n")
	_, err = readNumber(buf)
	require.Error(t, err)

	buf = bytes.NewBufferString("a\r\n")
	_, err = readNumber(buf)
	require.Error(t, err)

	buf = bytes.NewBufferString("3241a\r\n")
	_, err = readNumber(buf)
	require.Error(t, err)

	buf = bytes.NewBufferString("-3241a\r\n")
	_, err = readNumber(buf)
	require.Error(t, err)

	buf = bytes.NewBufferString("-a\r\n")
	_, err = readNumber(buf)
	require.Error(t, err)

	buf = bytes.NewBufferString("555\n\r")
	_, err = readNumber(buf)
	require.Error(t, err)

	buf = bytes.NewBufferString("555\r\r")
	_, err = readNumber(buf)
	require.Error(t, err)
}
