package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/value/geo"
	"github.com/codecrafters-io/redis-starter-go/pkg/value/sorted"
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

	inserted := set.Add(geo.ToScore(cmd.Coord), cmd.Member)
	if !inserted {
		return enc.Integer(0), nil
	}

	return enc.Integer(1), nil
}
