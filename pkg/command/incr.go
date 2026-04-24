package command

import (
	"fmt"
	"time"
)

type Incr struct {
	Key  string
	Time time.Time
}

func (in *Incr) isCommand() {}

func ParseIncr(args []string) (Incr, error) {
	if len(args) < 1 {
		return Incr{}, fmt.Errorf("INCR must have a key")
	}

	return Incr{
		Key:  args[0],
		Time: time.Now(),
	}, nil
}
