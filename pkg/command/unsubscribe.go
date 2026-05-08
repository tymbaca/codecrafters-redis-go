package command

import (
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

type Unubscribe struct {
	Chans []string
	Context
}

func ParseUnubscribe(ctx Context, args []string) (Unubscribe, error) {
	if len(args) < 1 {
		return Unubscribe{}, enc.ErrWrongNumberOfArgument("unsubscribe")
	}

	return Unubscribe{
		Chans:   args,
		Context: ctx,
	}, nil
}
