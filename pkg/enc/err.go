package enc

import (
	"errors"
	"fmt"
)

var (
	ErrWrongType           = newError(errors.New("WRONGTYPE Operation against a key holding the wrong kind of value"))
	ErrSyntaxError         = newError(errors.New("ERR syntax error"))
	ErrNotInteger          = newError(errors.New("ERR value is not an integer or out of range"))
	ErrNotFloat            = newError(errors.New("ERR value is not a valid float"))
	ErrNestedMulti         = newError(errors.New("ERR MULTI calls can not be nested"))
	ErrExecWithoutMulti    = newError(errors.New("ERR EXEC without MULTI"))
	ErrDiscardWithoutMulti = newError(errors.New("ERR DISCARD without MULTI"))
	ErrUnknownCommand      = func(name string) error { return newError(fmt.Errorf("ERR unknown command '%s'", name)) }

	ErrWrongNumberOfArgument = func(cmd string) error {
		return newError(fmt.Errorf("ERR wrong number of arguments for '%s' command"))
	}
	ErrConfigSetFailed = func(key string, desc string) error {
		return newError(fmt.Errorf("ERR CONFIG SET failed (possibly related to argument '%s') - %s", key, desc))
	}
	ErrConfigSetUnknownOption = func(key string) error {
		return newError(fmt.Errorf("ERR Unknown option or number of arguments for CONFIG SET - '%s'", key))
	}
	ErrCantExecInSubcriberMode = func(cmd string) error {
		return newError(fmt.Errorf("ERR Can't execute '%s': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context", cmd))
	}
)

type Error struct {
	err error
}

func (e Error) Error() string {
	return e.err.Error()
}

func IsError(err error) bool {
	var encErr Error
	return errors.As(err, &encErr)
}

func newError(err error) error {
	return Error{err: err}
}
