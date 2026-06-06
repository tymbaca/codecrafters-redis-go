package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/internal/service/value/sorted"
	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) geoadd(_ context.Context, cmd command.GeoAdd) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var set *sorted.Set

	entr, ok := s.data[cmd.Key]
	if !ok {
		set = sorted.New()
		s.data[cmd.Key] = entry{val: set}
	} else {
		set, ok = entr.val.(*sorted.Set)
		if !ok {
			return nil, enc.ErrWrongType
		}
	}

	inserted := set.Add(geoScore(cmd.Lon, cmd.Lat), cmd.Member)
	if !inserted {
		return enc.Integer(0), nil
	}

	return enc.Integer(1), nil
}

func geoScore(lon, lat float64) float64 {
	return 0
}
