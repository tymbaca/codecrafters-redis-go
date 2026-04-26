package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) Get(ctx context.Context, cmd command.Get) (enc.Value, error) {
	pre, err := s.prelude(ctx, cmd, true)
	if err != nil || pre != nil {
		return pre, err
	}

	val, set, err := s.get(ctx, cmd)
	if err != nil {
		return nil, err
	}

	reply := enc.BulkString{Val: val, Null: !set}
	return reply, nil
}

func (s *Service) get(ctx context.Context, cmd command.Get) (string, bool, error) {
	s.dataMu.RLock()
	defer s.dataMu.RUnlock()

	entry, ok := s.data[cmd.Key]

	if entry.expireSet && entry.expire.Before(cmd.Time) {
		delete(s.data, cmd.Key)
		return "", false, nil
	}

	if !ok {
		return "", false, nil
	}

	return entry.val, true, nil
}
