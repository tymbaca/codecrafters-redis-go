package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type GeoAdd struct {
	Key      string
	Lon, Lat float64
	Member   string
	Context
}

func ParseGeoAdd(ctx Context, args []string) (cmd GeoAdd, err error) {
	cmd = GeoAdd{
		Context: ctx,
	}

	iter := iter.New(args)
	var ok bool

	cmd.Key, ok = iter.Next()
	if !ok {
		return GeoAdd{}, enc.ErrWrongNumberOfArgument("geoadd")
	}

	cmd.Lon, err = parseFloat("geoadd", iter)
	if err != nil {
		return GeoAdd{}, err
	}

	cmd.Lat, err = parseFloat("geoadd", iter)
	if err != nil {
		return GeoAdd{}, err
	}

	cmd.Member, ok = iter.Next()
	if !ok {
		return GeoAdd{}, enc.ErrWrongNumberOfArgument("geoadd")
	}

	return cmd, nil
}

func parseFloat(cmd string, it *iter.Iterator[string]) (float64, error) {
	str, ok := it.Next()
	if !ok {
		return 0, enc.ErrWrongNumberOfArgument(cmd)
	}

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, enc.ErrNotFloat
	}

	return f, nil
}
