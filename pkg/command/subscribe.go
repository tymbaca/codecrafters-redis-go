package command

import "github.com/codecrafters-io/redis-starter-go/pkg/enc"

type Subscribe struct {
	Chans []string
	Context
}

func ParseSubscribe(ctx Context, args []string) (Subscribe, error) {
	if len(args) < 1 {
		return Subscribe{}, enc.ErrWrongNumberOfArgument("subscribe")
	}

	return Subscribe{
		Chans:   args,
		Context: ctx,
	}, nil
}
