package command

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type Get struct {
	Key  string
	Time time.Time
	Context
}

func ParseGet(ctx Context, args []string) (Get, error) {
	cmd := Get{}
	cmd.Time = time.Now()
	cmd.Context = ctx

	iter := iter.Iter(args)
	key, ok := iter.Next()
	if !ok {
		return Get{}, fmt.Errorf("expected GET key")
	}
	cmd.Key = key

	return cmd, nil
}
