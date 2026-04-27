package service

import (
	"context"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) incr(ctx context.Context, cmd command.Incr) (enc.Value, error) {
	valStr, set, err := s.getInner(ctx, command.Get(cmd))
	if err != nil {
		return nil, err
	}

	if !set {
		valStr = "0"
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return errValue(ErrNotInteger), nil
	}

	val++

	_, _, err = s.setInner(ctx, command.Set{
		Key:  cmd.Key,
		Val:  strconv.Itoa(val),
		Time: cmd.Time,
	})
	if err != nil {
		return nil, err
	}

	return enc.Integer(val), nil
}

func errValue(err error) enc.Value {
	return enc.SimpleError(err.Error())
}
