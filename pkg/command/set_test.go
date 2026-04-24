package command_test

import (
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/stretchr/testify/require"
)

func TestParseSet(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		args    []string
		want    command.Set
		wantErr bool
	}{
		{
			name:    "basic",
			args:    []string{"k1", "v1"},
			want:    command.Set{Key: "k1", Val: "v1"},
			wantErr: false,
		},
		{
			name:    "EX",
			args:    []string{"k1", "v1", "EX", "10"},
			want:    command.Set{Key: "k1", Val: "v1", ExpireSet: true, Expire: 10 * time.Second},
			wantErr: false,
		},
		{
			name:    "PX",
			args:    []string{"k1", "v1", "PX", "1000"},
			want:    command.Set{Key: "k1", Val: "v1", ExpireSet: true, Expire: 1000 * time.Millisecond},
			wantErr: false,
		},
		{
			name:    "NX",
			args:    []string{"k1", "v1", "NX", "PX", "1000"},
			want:    command.Set{Key: "k1", Val: "v1", Exists: command.ExistsKindNX, ExpireSet: true, Expire: 1000 * time.Millisecond},
			wantErr: false,
		},
		{
			name:    "XX",
			args:    []string{"k1", "v1", "XX", "PX", "1000"},
			want:    command.Set{Key: "k1", Val: "v1", Exists: command.ExistsKindXX, ExpireSet: true, Expire: 1000 * time.Millisecond},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := command.ParseSet(tt.args)
			if tt.wantErr {
				require.Error(t, gotErr)
				return
			}

			got.Time = time.Time{}
			require.NoError(t, gotErr)
			require.Equal(t, tt.want, got)
		})
	}
}
