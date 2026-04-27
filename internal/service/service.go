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

func (s *Service) Exec(ctx context.Context, cmd command.Command) (enc.Value, error) {
	queue := true
	switch cmd.(type) {
	case command.Multi, command.Exec, command.Discard:
		// no need to queue this commands in transaction
		queue = false
	}

	pre, err := s.prelude(ctx, cmd, queue)
	if err != nil || pre != nil {
		return pre, err
	}

	return s.execCmd(ctx, cmd)
}

func (s *Service) execCmd(ctx context.Context, cmd command.Command) (enc.Value, error) {
	switch cmd := cmd.(type) {
	case command.Multi:
		return s.multi(ctx, cmd)
	case command.Exec:
		return s.exec(ctx, cmd)
	case command.Discard:
		return s.discard(ctx, cmd)
	case command.Get:
		return s.get(ctx, cmd)
	case command.Set:
		return s.set(ctx, cmd)
	case command.Incr:
		return s.incr(ctx, cmd)
	}

	panic("unreachable")
}

func (s *Service) execCmds(ctx context.Context, cmds []command.Command) (enc.Value, error) {
	var arr enc.Array
	for _, cmd := range cmds {
		val, err := s.execCmd(ctx, cmd)
		if err != nil {
			return nil, err
		}

		arr = append(arr, val)
	}

	return arr, nil
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
