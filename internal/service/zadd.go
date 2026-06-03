package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/internal/service/value/sorted"
	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) zadd(_ context.Context, cmd command.ZAdd) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var set *sorted.Set

	entry, ok := s.data[cmd.Key]
	if !ok {
		set = sorted.New()
	} else {
		set, ok = entry.val.(*sorted.Set)
		if !ok {
			return nil, enc.ErrWrongType
		}
	}

	inserted := set.Add(cmd.Score, cmd.Member)
	if inserted {
		return enc.Integer(1), nil
	}

	return enc.Integer(0), nil
}
