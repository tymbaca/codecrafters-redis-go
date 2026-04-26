package service

import (
	"context"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/option"
)

func (s *Service) Set(ctx context.Context, cmd command.Set) (enc.Value, error) {
	if err := s.wal.Append(ctx, cmd); err != nil {
		return nil, fmt.Errorf("append to WAL: %w", err)
	}

	old, setOk, err := s.set(ctx, cmd)
	if err != nil {
		return nil, err
	}

	if cmd.GetOld {
		oldVal, oldSet := old.Get()
		return enc.BulkString{Val: oldVal, Null: !oldSet}, nil
	}

	if !setOk {
		return enc.BulkString{Null: true}, nil
	}

	return enc.OK, nil
}

func (s *Service) set(ctx context.Context, cmd command.Set) (option.Value[string], bool, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	oldValue, exists := s.data[cmd.Key]
	oldValueOption := option.Wrap(oldValue.val, exists)

	if cmd.Exists == command.ExistsKindXX && !exists {
		return oldValueOption, false, nil
	}
	if cmd.Exists == command.ExistsKindNX && exists {
		return oldValueOption, false, nil
	}

	s.data[cmd.Key] = entry{
		val:       cmd.Val,
		expireSet: cmd.ExpireSet,
		expire:    cmd.Time.Add(cmd.Expire),
	}

	return oldValueOption, true, nil
}
