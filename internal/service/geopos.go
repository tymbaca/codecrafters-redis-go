package service

import (
	"context"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/value/geo"
	"github.com/codecrafters-io/redis-starter-go/pkg/value/sorted"
)

func (s *Service) geopos(_ context.Context, cmd command.GeoPos) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var set *sorted.Set

	entr, ok := s.data[cmd.Key]
	if !ok {
		return enc.NullArr(), nil
	}

	set, ok = entr.val.(*sorted.Set)
	if !ok {
		return nil, enc.ErrWrongType
	}

	var result []enc.Value

	for _, member := range cmd.Members {
		score, ok := set.Score(member)
		if !ok {
			return enc.NullArr(), nil
		}

		coord := geo.FromScore(score)
		result = append(result, enc.Arr(
			enc.Bulk(fmt.Sprint(coord.Lon)),
			enc.Bulk(fmt.Sprint(coord.Lat)),
		))
	}

	return enc.Arr(result...), nil
}
