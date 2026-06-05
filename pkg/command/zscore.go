package command

import (
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type ZScore struct {
	Key    string
	Member string
	Context
}

func ParseZScore(ctx Context, args []string) (cmd ZScore, err error) {
	cmd = ZScore{
		Context: ctx,
	}

	iter := iter.New(args)
	var ok bool

	cmd.Key, ok = iter.Next()
	if !ok {
		return ZScore{}, enc.ErrWrongNumberOfArgument("zscore")
	}

	cmd.Member, ok = iter.Next()
	if !ok {
		return ZScore{}, enc.ErrWrongNumberOfArgument("zscore")
	}

	return cmd, nil
}
