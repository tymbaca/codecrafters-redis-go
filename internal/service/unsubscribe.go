package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) unsubscribe(_ context.Context, cmd command.Unubscribe) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	for _, ch := range cmd.Chans {
		// delete conn from channel
		subscribers := s.channels[ch]
		delete(subscribers, cmd.ConnID)
		s.channels[ch] = subscribers

		// delete channel from conn metadata
		if s.subscribersMeta[cmd.ConnID].channels != nil {
			delete(s.subscribersMeta[cmd.ConnID].channels, ch)
			if len(s.subscribersMeta[cmd.ConnID].channels) == 0 {
				// if no more channels left, delete metadata, so
				// this conn no more in subscriber mode
				delete(s.subscribersMeta, cmd.ConnID)
			}
		}

		cmd.Conn.Send(enc.Array{
			enc.Bulk("unsubscribe"),
			enc.Bulk(ch),
			enc.Integer(len(s.subscribersMeta[cmd.ConnID].channels)),
		})
	}

	return nil, nil
}
