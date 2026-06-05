package command

import (
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type ZCard struct {
	Key string
	Context
}

func ParseZCard(ctx Context, args []string) (ZCard, error) {
	cmd := ZCard{
		Context: ctx,
	}

	iter := iter.New(args)
	var ok bool

	cmd.Key, ok = iter.Next()
	if !ok {
		return ZCard{}, enc.ErrWrongNumberOfArgument("zcard")
	}

	return cmd, nil
}
