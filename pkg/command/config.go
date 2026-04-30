package command

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type Config struct {
	Kind ConfigKind
	Keys []string
	Vals []string // for [ConfigSet]
	Context
}

type ConfigKind int

const (
	ConfigGet ConfigKind = iota
	ConfigSet
)

func ParseConfig(ctx Context, args []string) (cmd Config, err error) {
	cmd.Context = ctx

	iter := iter.New(args)
	kind, ok := iter.Next()
	if !ok {
		return cmd, fmt.Errorf("expected GET or SET")
	}

	switch strings.ToUpper(kind) {
	case "GET":
		cmd.Kind = ConfigGet
	case "SET":
		cmd.Kind = ConfigSet
	default:
		return cmd, fmt.Errorf("expected GET or SET")
	}
	key, ok := iter.Next()
	if !ok {
		return cmd, enc.ErrSyntaxError
	}
	cmd.Keys = append(cmd.Keys, key)

	if cmd.Kind == ConfigSet {
		val, ok := iter.Next()
		if !ok {
			return cmd, enc.ErrSyntaxError
		}
		cmd.Vals = append(cmd.Vals, val)
	}

	for {
		key, ok := iter.Next()
		if !ok {
			break
		}
		cmd.Keys = append(cmd.Keys, key)

		if cmd.Kind == ConfigSet {
			val, ok := iter.Next()
			if !ok {
				return cmd, enc.ErrSyntaxError
			}
			cmd.Vals = append(cmd.Vals, val)
		}
	}

	return cmd, nil
}
