package enc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseNumber(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("456\r\n")
	num, err := readNumber(buf)
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
