package storage

import (
	"context"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/option"
)

func New() *Storage {
	return &Storage{
		data: make(map[string]entry),
	}
}

type Storage struct {
	mu   sync.RWMutex
	data map[string]entry
}

type entry struct {
	val       string
	expireSet bool
	expire    time.Time
}

func (s *Storage) Get(ctx context.Context, cmd command.Get) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

func (s *Storage) Set(ctx context.Context, cmd command.Set) (option.Value[string], bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
