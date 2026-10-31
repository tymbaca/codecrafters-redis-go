package service

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/value/geo"
	"github.com/stretchr/testify/require"
)

func TestGeoPos(t *testing.T) {
	ctx := t.Context()
	svc := newService(t)
	buf := bytes.NewBuffer(nil)
	cmdCtx := command.NewContext(buf)

	coord := geo.Coord{Lat: 10, Lon: 10}

	val, err := svc.geoadd(ctx, command.GeoAdd{
		Key:     "test1",
		Coord:   coord,
		Member:  "member1",
		Context: cmdCtx,
	})
	require.NoError(t, err)
	require.Equal(t, enc.Integer(1), val)

	score := geo.ToScore(coord)
	coord = geo.FromScore(score)

	val, err = svc.geopos(ctx, command.GeoPos{
		Key:     "test1",
		Members: []string{"member1"},
		Context: cmdCtx,
	})
	require.NoError(t, err)
	require.Equal(t, enc.Arr(
		enc.Bulk(fmt.Sprint(coord.Lon)),
		enc.Bulk(fmt.Sprint(coord.Lat)),
	), val)
}

func newService(t testing.TB) *Service {
	svc, err := New(t.Context(), Options{})
	require.NoError(t, err)

	return svc
}
