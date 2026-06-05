package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/internal/service/value/geo"
	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) geoadd(_ context.Context, cmd command.GeoAdd) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var g *geo.Geo

	entr, ok := s.data[cmd.Key]
	if !ok {
		g = geo.New()
		s.data[cmd.Key] = entry{val: g}
	} else {
		g, ok = entr.val.(*geo.Geo)
		if !ok {
			return nil, enc.ErrWrongType
		}
	}

	inserted := g.Add(cmd.Lon, cmd.Lat, cmd.Member)
	if !inserted {
		return enc.Integer(0), nil
	}

	return enc.Integer(1), nil
}
