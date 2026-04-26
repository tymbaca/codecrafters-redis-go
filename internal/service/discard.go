package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) Discard(ctx context.Context, cmd command.Discard) (enc.Value, error) {
	s.txsMu.Lock()
	defer s.txsMu.Unlock()

	if _, exists := s.txs[cmd.ConnID]; !exists {
		return errValue(ErrDiscardWithoutMulti), nil
	}

	delete(s.txs, cmd.ConnID)

	return enc.OK, nil
}
