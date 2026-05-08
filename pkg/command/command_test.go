package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetName(t *testing.T) {
	require.Equal(t, "GET", GetName(Command(Get{})))
	require.Equal(t, "GET", GetName(Get{}))
	require.Equal(t, "EXEC", GetName(Exec{}))
	require.Equal(t, "SUBSCRIBE", GetName(Subscribe{}))
}
