package service

import (
	"context"
	"log/slog"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"golang.org/x/sync/errgroup"
)

func (s *Service) publish(_ context.Context, cmd command.Publish) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()
	slog.Debug("publishing to channel", "chan", cmd.Chan)

	subscribers := s.channels[cmd.Chan]

	var wg errgroup.Group
	for connID, conn := range subscribers {
		wg.Go(func() error {
			slog.Debug("publish: sending to conn", "connID", connID)
			return conn.Send(enc.Array{
				enc.Bulk("message"),
				enc.Bulk(cmd.Chan),
				enc.Bulk(cmd.Val),
			})
		})
	}

	err := wg.Wait()
	if err != nil {
		return errValue(err), nil
	}

	return enc.Integer(len(subscribers)), nil
}
