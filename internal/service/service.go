package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/aof"
	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

type Service struct {
	aof *aof.AOF
	wal wal

	dataMu sync.RWMutex
	data   map[string]entry
	txsMu  sync.Mutex
	txs    map[string][]command.Command

	cfgMu       sync.RWMutex
	dir         string
	appendOnly  bool
	appendDir   string
	appendFile  string
	appendFsync string
}

type Options struct {
	Dir         string
	AppendOnly  bool
	AppendDir   string
	AppendFile  string
	AppendFsync string
}

func New(ctx context.Context, opts Options) (*Service, error) {
	svc := &Service{
		txs:         make(map[string][]command.Command),
		data:        make(map[string]entry),
		wal:         noopWal{},
		dir:         opts.Dir,
		appendOnly:  opts.AppendOnly,
		appendDir:   opts.AppendDir,
		appendFile:  opts.AppendFile,
		appendFsync: opts.AppendFsync,
	}

	if svc.dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("get current dir: %w", err)
		}
		svc.dir = cwd
	}
	if svc.appendDir == "" {
		svc.appendDir = "appendonlydir"
	}
	if svc.appendFile == "" {
		svc.appendFile = "appendonly.aof"
	}
	if svc.appendFsync == "" {
		svc.appendFsync = "everysec"
	}

	if svc.appendOnly {
		slog.Info("HELLO AOF")
		aof, err := aof.New(ctx, svc.dir, svc.appendDir, svc.appendFile, func(ctx context.Context, cmd command.Command) error {
			_, err := svc.Exec(ctx, cmd)
			return err
		})
		if err != nil {
			return nil, err
		}

		svc.aof = aof
	}

	return svc, nil
}

func (s *Service) Intercept(ctx context.Context, cmdCtx command.Context, cmdVal enc.Value) error {
	s.cfgMu.Lock()
	defer s.cfgMu.Unlock()

	if s.aof == nil {
		return nil
	}

	cmd, err := command.Parse(cmdCtx, cmdVal)
	if err != nil {
		return nil
	}

	switch cmd.(type) {
	default:
		return nil
	case command.Config:
	case command.Discard:
	case command.Exec:
	case command.Incr:
	case command.Multi:
	case command.Set:
	}

	return s.aof.Append(ctx, cmdVal)
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
	case command.Config:
		return s.config(ctx, cmd)
	case command.Echo:
		return cmd.Val, nil
	case command.Ping:
		return enc.Pong, nil
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
