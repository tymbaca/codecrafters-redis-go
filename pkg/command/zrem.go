package command

import (
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type ZRem struct {
	Key    string
	Member string
	Context
}

func ParseZRem(ctx Context, args []string) (ZRem, error) {
	cmd := ZRem{
		Context: ctx,
	}

	iter := iter.New(args)
	var ok bool

	cmd.Key, ok = iter.Next()
	if !ok {
		return ZRem{}, enc.ErrWrongNumberOfArgument("zrem")
	}

	cmd.Member, ok = iter.Next()
	if !ok {
		return ZRem{}, enc.ErrWrongNumberOfArgument("zrem")
	}

	return cmd, nil
}
