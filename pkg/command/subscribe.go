package command

type Subscribe struct {
	Chans []string
	Context
}

func ParseSubscribe(ctx Context, args []string) (Subscribe, error) {
	return Subscribe{
		Chans:   args,
		Context: ctx,
	}, nil
}
