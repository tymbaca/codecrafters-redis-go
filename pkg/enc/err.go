package enc

import (
	"errors"
	"fmt"
)

var (
	ErrSyntaxError         = errors.New("ERR syntax error")
	ErrNotInteger          = errors.New("ERR value is not an integer or out of range")
	ErrNestedMulti         = errors.New("ERR MULTI calls can not be nested")
	ErrExecWithoutMulti    = errors.New("ERR EXEC without MULTI")
	ErrDiscardWithoutMulti = errors.New("ERR DISCARD without MULTI")
	ErrUnknownCommand      = func(name string) error { return fmt.Errorf("ERR unknown command '%s'", name) }

	ErrWrongNumberOfArgument = func(cmd string) error {
		return fmt.Errorf("ERR wrong number of arguments for '%s' command")
	}
	ErrConfigSetFailed = func(key string, desc string) error {
		return fmt.Errorf("ERR CONFIG SET failed (possibly related to argument '%s') - %s", key, desc)
	}
	ErrConfigSetUnknownOption = func(key string) error {
		return fmt.Errorf("ERR Unknown option or number of arguments for CONFIG SET - '%s'", key)
	}
	ErrCantExecInSubcriberMode = func(cmd string) error {
		return fmt.Errorf("ERR Can't execute '%s': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context", cmd)
	}
)
