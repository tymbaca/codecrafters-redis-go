package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) subscribe(ctx context.Context, cmd command.Subscribe) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	for i, ch := range cmd.Chans {
		subscribers := s.channels[ch] // map[command.Context]struct{}
		if subscribers == nil {
			subscribers = make(subscriberSet)
		}

		subscribers[cmd.Context] = struct{}{}
		s.channels[ch] = subscribers

		reply := enc.Array{enc.Bulk("subscribe"), enc.Bulk(ch), enc.Integer(i + 1)}
		err := reply.Encode(cmd.Conn)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
