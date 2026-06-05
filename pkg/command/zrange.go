package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type ZRange struct {
	Key      string
	Min, Max int
	Context
}

func ParseZRange(ctx Context, args []string) (cmd ZRange, err error) {
	cmd = ZRange{
		Context: ctx,
	}

	iter := iter.New(args)
	var ok bool

	cmd.Key, ok = iter.Next()
	if !ok {
		return ZRange{}, enc.ErrWrongNumberOfArgument("zrange")
	}

	minStr, ok := iter.Next()
	if !ok {
		return ZRange{}, enc.ErrWrongNumberOfArgument("zrange")
	}

	cmd.Min, err = strconv.Atoi(minStr)
	if err != nil {
		return ZRange{}, enc.ErrNotInteger
	}

	maxStr, ok := iter.Next()
	if !ok {
		return ZRange{}, enc.ErrWrongNumberOfArgument("zrange")
	}

	cmd.Max, err = strconv.Atoi(maxStr)
	if err != nil {
		return ZRange{}, enc.ErrNotInteger
	}

	return cmd, nil
}
