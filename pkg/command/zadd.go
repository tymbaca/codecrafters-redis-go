package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type ZAdd struct {
	Key    string
	Score  float64
	Member string
	Context
}

func ParseZAdd(ctx Context, args []string) (ZAdd, error) {
	cmd := ZAdd{
		Context: ctx,
	}

	iter := iter.New(args)
	var ok bool

	cmd.Key, ok = iter.Next()
	if !ok {
		return ZAdd{}, enc.ErrWrongNumberOfArgument("zadd")
	}

	scoreStr, ok := iter.Next()
	if !ok {
		return ZAdd{}, enc.ErrWrongNumberOfArgument("zadd")
	}

	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		return ZAdd{}, enc.ErrNotFloat
	}
	cmd.Score = score

	cmd.Member, ok = iter.Next()
	if !ok {
		return ZAdd{}, enc.ErrWrongNumberOfArgument("zadd")
	}

	return cmd, nil
}
