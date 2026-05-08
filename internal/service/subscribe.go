package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) subscribe(ctx context.Context, cmd command.Subscribe) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	for _, ch := range cmd.Chans {
		s.registerSubscriber(ctx, cmd.Context, ch)
	}

	return nil, nil
}

func (s *Service) registerSubscriber(ctx context.Context, cmdCtx command.Context, ch string) error {
	subscribers := s.channels[ch] // map[command.Context]struct{}
	if subscribers == nil {
		subscribers = make(subscriberSet)
	}

	_, ok := subscribers[cmdCtx.ConnID]
	if !ok {
		// assign subscriber
		subscribers[cmdCtx.ConnID] = cmdCtx.Conn
		s.channels[ch] = subscribers

		// increment this subscribers channel count
		subMeta := s.subscribersMeta[cmdCtx.ConnID]
		if subMeta.channels == nil {
			subMeta.channels = make(map[string]struct{})
		}
		subMeta.channels[ch] = struct{}{}
		s.subscribersMeta[cmdCtx.ConnID] = subMeta
	}

	return cmdCtx.Conn.Send(enc.Array{
		enc.Bulk("subscribe"),
		enc.Bulk(ch),
		enc.Integer(len(s.subscribersMeta[cmdCtx.ConnID].channels)),
	})
}
