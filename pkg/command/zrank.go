package command

import (
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type ZRank struct {
	Key    string
	Member string
	Context
}

func ParseZRank(ctx Context, args []string) (ZRank, error) {
	cmd := ZRank{
		Context: ctx,
	}

	iter := iter.New(args)
	var ok bool

	cmd.Key, ok = iter.Next()
	if !ok {
		return ZRank{}, enc.ErrWrongNumberOfArgument("zrank")
	}

	cmd.Member, ok = iter.Next()
	if !ok {
		return ZRank{}, enc.ErrWrongNumberOfArgument("zrank")
	}

	return cmd, nil
}
