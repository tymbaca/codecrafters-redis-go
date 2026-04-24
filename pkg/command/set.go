package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/pkg/iter"
)

type Set struct {
	Key       string
	Val       string
	Exists    ExistsKind
	ExpireSet bool
	Expire    time.Duration
	GetOld    bool
	Time      time.Time
}

func (se Set) isCommand() {}

type ExistsKind int

const (
	ExistsKindNone ExistsKind = iota
	ExistsKindNX              // if not exists
	ExistsKindXX              // if exists
)

func ParseSet(args []string) (Set, error) {
	cmd := Set{}
	cmd.Time = time.Now()

	iter := iter.Iter(args)
	key, ok := iter.Next()
	if !ok {
		return Set{}, fmt.Errorf("expected SET key")
	}
	cmd.Key = key

	val, ok := iter.Next()
	if !ok {
		return Set{}, fmt.Errorf("expected SET value")
	}
	cmd.Val = val

	for {
		arg, ok := iter.Next()
		if !ok {
			break
		}

		switch strings.ToUpper(arg) {
		case "EX":
			secStr, ok := iter.Next()
			if !ok {
				return Set{}, fmt.Errorf("expected EX duration")
			}

			sec, err := strconv.Atoi(secStr)
			if err != nil {
				return Set{}, fmt.Errorf("expected EX duration to be an integer: %w", err)
			}

			cmd.ExpireSet = true
			cmd.Expire = time.Duration(sec) * time.Second

		case "PX":
			msStr, ok := iter.Next()
			if !ok {
				return Set{}, fmt.Errorf("expected PX duration")
			}

			ms, err := strconv.Atoi(msStr)
			if err != nil {
				return Set{}, fmt.Errorf("expected PX duration to be an integer: %w", err)
			}

			cmd.ExpireSet = true
			cmd.Expire = time.Duration(ms) * time.Millisecond

		case "NX":
			cmd.Exists = ExistsKindNX
		case "XX":
			cmd.Exists = ExistsKindXX

		case "GET":
			cmd.GetOld = true
		}
	}

	return cmd, nil
}
