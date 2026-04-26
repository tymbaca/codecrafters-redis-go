package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func New() *Service {
	return &Service{
		txs:  make(map[string][]command.Command),
		data: make(map[string]entry),
		wal:  noopWal{},
	}
}

type Service struct {
	wal wal

	dataMu sync.RWMutex
	data   map[string]entry
	txsMu  sync.Mutex
	txs    map[string][]command.Command
}

func (s *Service) prelude(ctx context.Context, cmd command.Command, queue bool) (enc.Value, error) {
	if ctx.Value(ignorePreludeKey) != nil {
		return nil, nil
	}

	if err := s.wal.Append(ctx, cmd); err != nil {
		return nil, fmt.Errorf("append to WAL: %w", err)
	}

	if queue {
		if queued := s.txQueue(ctx, cmd); queued {
			return enc.Queued, nil
		}
	}

	return nil, nil
}

type wal interface {
	Append(ctx context.Context, cmds ...command.Command) error
}

type noopWal struct{}

func (w noopWal) Append(ctx context.Context, cmds ...command.Command) error { return nil }

type entry struct {
	val       string
	expireSet bool
	expire    time.Time
}
