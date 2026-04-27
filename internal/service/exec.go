package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) exec(ctx context.Context, cmd command.Exec) (enc.Value, error) {
	s.txsMu.Lock()
	defer s.txsMu.Unlock()

	queue, exists := s.txs[cmd.ConnID]
	if !exists {
		return errValue(ErrExecWithoutMulti), nil
	}

	ctx = context.WithValue(ctx, ignorePreludeKey, ignorePreludeKey)
	execVals, err := s.execCmds(ctx, queue)
	if err != nil {
		return nil, err
	}

	delete(s.txs, cmd.ConnID)

	return execVals, nil
}

type ignorePrelude struct{}

var ignorePreludeKey ignorePrelude = ignorePrelude{}

func (s *Service) txQueue(ctx context.Context, cmd command.Command) bool {
	s.txsMu.Lock()
	defer s.txsMu.Unlock()

	if queue, ok := s.txs[cmd.Ctx().ConnID]; ok {
		queue = append(queue, cmd)
		s.txs[cmd.Ctx().ConnID] = queue

		return true
	}

	return false
}
