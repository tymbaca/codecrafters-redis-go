package service

import (
	"context"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/value/sorted"
)

func (s *Service) zrange(_ context.Context, cmd command.ZRange) (enc.Value, error) {
	s.dataMu.Lock()
	defer s.dataMu.Unlock()

	var resultArray enc.Array

	entr, ok := s.data[cmd.Key]
	if !ok {
		return resultArray, nil
	}

	set, ok := entr.val.(*sorted.Set)
	if !ok {
		return nil, enc.ErrWrongType
	}

	result := set.Range(cmd.Min, cmd.Max)
	for _, v := range result {
		resultArray = append(resultArray, enc.Bulk(v))
	}

	return resultArray, nil
}
