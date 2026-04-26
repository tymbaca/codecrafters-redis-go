package service

import "errors"

var (
	ErrNotInteger          = errors.New("ERR value is not an integer or out of range")
	ErrNestedMulti         = errors.New("ERR MULTI calls can not be nested")
	ErrExecWithoutMulti    = errors.New("ERR EXEC without MULTI")
	ErrDiscardWithoutMulti = errors.New("ERR DISCARD without MULTI")
)
