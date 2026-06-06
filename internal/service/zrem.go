package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/value/sorted"
)

func (s *Service) zrem(_ context.Context, cmd command.ZRem) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	entr, ok := s.data[cmd.Key]
	if !ok {
		return enc.Integer(0), nil
	}

	set, ok := entr.val.(*sorted.Set)
	if !ok {
		return nil, enc.ErrWrongType
	}

	removed := set.Remove(cmd.Member)
	if !removed {
		return enc.Integer(0), nil
	}

	return enc.Integer(1), nil
}
