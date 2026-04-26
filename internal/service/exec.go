package service

import (
	"context"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) Exec(ctx context.Context, cmd command.Exec) (enc.Value, error) {
	pre, err := s.prelude(ctx, cmd, false)
	if err != nil || pre != nil {
		return pre, err
	}

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

func (s *Service) execCmds(ctx context.Context, cmds []command.Command) (enc.Value, error) {
	var arr enc.Array
	for _, cmd := range cmds {
		var val enc.Value
		var err error

		switch cmd := cmd.(type) {
		case command.Get:
			val, err = s.Get(ctx, cmd)
		case command.Set:
			val, err = s.Set(ctx, cmd)
		case command.Incr:
			val, err = s.Incr(ctx, cmd)
		default:
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("exec %#v: %w", cmd, err)
		}

		arr = append(arr, val)
	}

	return arr, nil
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
