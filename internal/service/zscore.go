package service

import (
	"context"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/service/value/sorted"
	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) zscore(_ context.Context, cmd command.ZScore) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	entr, ok := s.data[cmd.Key]
	if !ok {
		return enc.Null(), nil
	}

	set, ok := entr.val.(*sorted.Set)
	if !ok {
		return nil, enc.ErrWrongType
	}

	score, ok := set.Score(cmd.Member)
	if !ok {
		return enc.Null(), nil
	}

	return enc.Bulk(fmt.Sprint(score)), nil
}
