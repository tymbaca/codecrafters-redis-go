package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) Multi(ctx context.Context, cmd command.Multi) (enc.Value, error) {
	s.txsMu.Lock()
	defer s.txsMu.Unlock()

	_, alreadyExists := s.txs[cmd.ConnID]
	if alreadyExists {
		return errValue(ErrNestedMulti), nil
	}

	s.txs[cmd.ConnID] = make([]command.Command, 0)

	return enc.OK, nil
}
