package command

import (
	"fmt"
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

	if cmd.Lon < -180 || cmd.Lon > 180 {
		return GeoAdd{}, fmt.Errorf("ERR longitude value (%f) is invalid", cmd.Lon)
	}

	cmd.Lat, err = parseFloat("geoadd", iter)
	if err != nil {
		return GeoAdd{}, err
	}

	if cmd.Lat < -85.05112878 || cmd.Lat > 85.05112878 {
		return GeoAdd{}, fmt.Errorf("ERR latitude value (%f) is invalid", cmd.Lon)
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
