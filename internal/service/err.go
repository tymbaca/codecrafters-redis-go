package service

import (
	"errors"
	"fmt"
)

var (
	ErrNotInteger          = errors.New("ERR value is not an integer or out of range")
	ErrNestedMulti         = errors.New("ERR MULTI calls can not be nested")
	ErrExecWithoutMulti    = errors.New("ERR EXEC without MULTI")
	ErrDiscardWithoutMulti = errors.New("ERR DISCARD without MULTI")
	ErrUnknownCommand      = func(name string) error { return fmt.Errorf("ERR unknown command '%s'", name) }
)
