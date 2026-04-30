package service

import (
	"context"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func (s *Service) config(ctx context.Context, cmd command.Config) (enc.Value, error) {
	s.cfgMu.Lock()
	defer s.cfgMu.Unlock()

	switch cmd.Kind {
	case command.ConfigGet:
		return s.configGet(ctx, cmd)
	case command.ConfigSet:
		return s.configSet(ctx, cmd)
	default:
		panic(fmt.Sprintf("unexpected command.ConfigKind: %#v", cmd.Kind))
	}
}

func (s *Service) configGet(_ context.Context, cmd command.Config) (enc.Value, error) {
	var res enc.Array

	for i := range cmd.Keys {
		key := cmd.Keys[i]
		var val string

		switch key {
		case "dir":
			val = s.dir
		case "appendonly":
			switch s.appendOnly {
			case true:
				val = "yes"
			case false:
				val = "no"
			}
		case "appenddirname":
			val = s.appendDir
		case "appendfilename":
			val = s.appendFilename
		case "appendfsync":
			val = s.appendFsync
		default:
			continue
		}

		res = append(res, enc.BulkString{Val: key})
		res = append(res, enc.BulkString{Val: val})
	}

	return res, nil
}

func (s *Service) configSet(_ context.Context, cmd command.Config) (enc.Value, error) {
	for i := range cmd.Keys {
		key := cmd.Keys[i]
		val := cmd.Vals[i]

		switch key {
		case "dir":
			return errValue(enc.ErrConfigSetFailed(key, "can't set protected config")), nil
		case "appendonly":
			switch val {
			case "yes":
				s.appendOnly = true
			case "no":
				s.appendOnly = false
			}
		case "appenddirname":
			s.appendDir = val
		case "appendfilename":
			s.appendFilename = val
		case "appendfsync":
			s.appendFsync = val
		default:
			return errValue(enc.ErrConfigSetUnknownOption(key)), nil
		}
	}

	return enc.OK, nil
}
