package command

import (
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type GeoPos struct {
	Key     string
	Members []string
	Context
}

func ParseGeoPos(ctx Context, args []string) (cmd GeoPos, err error) {
	cmd = GeoPos{
		Context: ctx,
	}

	iter := iter.New(args)
	var ok bool

	cmd.Key, ok = iter.Next()
	if !ok {
		return GeoPos{}, enc.ErrWrongNumberOfArgument("geopos")
	}

	member, ok := iter.Next()
	if !ok {
		return GeoPos{}, enc.ErrWrongNumberOfArgument("geoadd")
	}

	cmd.Members = append(cmd.Members, member)

	for {
		member, ok := iter.Next()
		if !ok {
			break
		}

		cmd.Members = append(cmd.Members, member)
	}

	return cmd, nil
}
