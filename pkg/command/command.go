// Package command provides structures and parse functions for redis commands.
package command

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/google/uuid"
)

type Command interface {
	Ctx() Context
	isCommand()
}

func NewContext(conn io.ReadWriter) Context {
	return Context{
		ConnID: uuid.NewString(),
		Conn:   newConn(conn),
	}
}

type Context struct {
	ConnID ConnID
	Conn   *Conn
}

type ConnID = string

func (c Context) Ctx() Context { return c }
func (c Context) isCommand()   {}

func Parse(ctx Context, cmd enc.Value) (Command, error) {
	arr, ok := cmd.(enc.Array)
	if !ok {
		return nil, fmt.Errorf("cmd must be array of values")
	}

	commandStr := arr[0].(enc.BulkString)
	args, err := argsToString(arr[1:])
	if err != nil {
		return nil, fmt.Errorf("parse command args (%#v): %w", arr, err)
	}

	switch strings.ToUpper(commandStr.Val) {
	case "PING":
		return Ping{Context: ctx}, nil
	case "ECHO":
		if len(arr) != 2 {
			return nil, fmt.Errorf("invalid ECHO command: must be 1 arg, got: %#v", args)
		}

		return Echo{Val: arr[1], Context: ctx}, nil
	case "GET":
		if len(args) < 1 {
			return nil, fmt.Errorf("invalid GET command: must be 1 or more args, got: %#v", args)
		}

		return ParseGet(ctx, args)
	case "SET":
		if len(args) < 2 {
			return nil, fmt.Errorf("invalid SET command: must be 3 or more args, got: %#v", args)
		}

		return ParseSet(ctx, args)
	case "INCR":
		if len(args) < 1 {
			return nil, fmt.Errorf("invalid INCR command: must be 1 arg, got: %#v", args)
		}

		return ParseIncr(ctx, args)
	case "MULTI":
		return Multi{Context: ctx}, nil
	case "EXEC":
		return Exec{Context: ctx}, nil
	case "DISCARD":
		return Discard{Context: ctx}, nil
	case "CONFIG":
		return ParseConfig(ctx, args)
	case "SUBSCRIBE":
		return ParseSubscribe(ctx, args)
	case "UNSUBSCRIBE":
		return ParseUnubscribe(ctx, args)
	case "PUBLISH":
		return ParsePublish(ctx, args)
	case "ZADD":
		return ParseZAdd(ctx, args)
	default:
		return nil, enc.ErrUnknownCommand(commandStr.Val)
	}
}

type Conn struct {
	mu sync.Mutex
	rw io.ReadWriter
}

func newConn(rw io.ReadWriter) *Conn {
	return &Conn{
		rw: rw,
	}
}

func (c *Conn) Send(val enc.Value) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return val.Encode(c.rw)
}

func (c *Conn) Recv() (enc.Value, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return enc.Decode(c.rw)
}

func argsToString(array enc.Array) ([]string, error) {
	res := make([]string, 0, len(array))
	for _, el := range array {
		bs, ok := el.(enc.BulkString)
		if !ok {
			return nil, fmt.Errorf("got non-bulk-string in request arguments: %#v", el)
		}

		res = append(res, bs.Val)
	}

	return res, nil
}

func GetName(cmd Command) string {
	return strings.ToUpper(reflect.ValueOf(cmd).Type().Name())
}
