package service

import (
	"context"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
)

func New() *Service {
	return &Service{
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
