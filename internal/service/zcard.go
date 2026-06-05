package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/internal/service/value/sorted"
	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) zcard(_ context.Context, cmd command.ZCard) (enc.Value, error) {
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

	return enc.Integer(set.Length()), nil
}
