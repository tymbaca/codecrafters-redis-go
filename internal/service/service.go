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
	}
}

type Service struct {
	mu   sync.RWMutex
	data map[string]entry
	wal  wal
}

type wal interface {
	Append(ctx context.Context, cmds ...command.Command) error
}

type entry struct {
	val       string
	expireSet bool
	expire    time.Time
}
