package command

import "fmt"

type Publish struct {
	Chan string
	Val  string
	Context
}

func ParsePublish(ctx Context, args []string) (Publish, error) {
	if len(args) < 2 {
		return Publish{}, fmt.Errorf("PUBLISH must have 2 arguments")
	}

	return Publish{
		Chan:    args[0],
		Val:     args[1],
		Context: ctx,
	}, nil
}
