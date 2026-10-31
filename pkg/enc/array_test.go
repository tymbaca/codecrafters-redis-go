package enc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNestedNullArr(t *testing.T) {
	val := Arr(NullArr(), NullArr(), NullArr())

	buf := bytes.NewBuffer(nil)

	err := val.Encode(buf)
	require.NoError(t, err)

	require.Equal(t, []byte("*3\r\n*-1\r\n*-1\r\n*-1\r\n"), buf.Bytes())
}
