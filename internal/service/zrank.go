package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/value/sorted"
)

func (s *Service) zrank(_ context.Context, cmd command.ZRank) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var set *sorted.Set

	entr, ok := s.data[cmd.Key]
	if !ok {
		return enc.Null(), nil
	}

	set, ok = entr.val.(*sorted.Set)
	if !ok {
		return nil, enc.ErrWrongType
	}

	rank, ok := set.Rank(cmd.Member)
	if !ok {
		return enc.Null(), nil
	}

	return enc.Integer(rank), nil
}
