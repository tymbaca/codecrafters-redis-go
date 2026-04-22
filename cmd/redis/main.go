package redis

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codecrafters-io/redis-starter-go/pkg/assert"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func Run() error {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-time.After(2 * time.Second)
		cancel()
	}()

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	context.AfterFunc(ctx, func() {
		_ = l.Close()
	})

	storage := &Storage{
		data: make(map[string]enc.Value),
	}

	for {
		conn, err := l.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		if err != nil {
			slog.Error("accept conn failed", "err", err)
			continue
		}

		go func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			context.AfterFunc(ctx, func() {
				_ = conn.Close()
			})

			err := handleConn(conn, storage)
			if err != nil {
				slog.Error("handle conn failed", "err", err)
			}
		}()
	}
}

func handleConn(conn io.ReadWriter, storage *Storage) error {
	for {
		command, ok, err := readCommand(conn)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		err = handleCommand(conn, storage, command)
		if err != nil {
			return err
		}
	}
}

func readCommand(conn io.Reader) (enc.Value, bool, error) {
	val, err := enc.ReadValue(conn)
	if errors.Is(err, io.EOF) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return val, true, nil
}

func handleCommand(conn io.Writer, storage *Storage, command enc.Value) error {
	if command.Type() != enc.TypeArray {
		return fmt.Errorf("exected array, got: %s (of type %s)", command.String(), command.Type())
	}

	arr := command.(enc.Array)
	assert.True(len(arr) > 0)
	assert.True(arr[0].Type() == enc.TypeBulkString)

	commandStr := arr[0].(enc.BulkString)

	switch commandStr.Val {
	case "PING":
		responseVal := enc.SimpleString("PONG")
		return responseVal.Encode(conn)

	case "ECHO":
		if len(arr) != 2 {
			return fmt.Errorf("invalid ECHO command: must be 2 elements, got: %#v", arr)
		}

		reply := arr[1]
		return reply.Encode(conn)

	case "SET":
		if len(arr) != 3 {
			return fmt.Errorf("invalid SET command: must be 3 elements, got: %#v", arr)
		}

		if arr[1].Type() != enc.TypeBulkString || arr[2].Type() != enc.TypeBulkString {
			return fmt.Errorf("invalid SET command: key and value must be bulk strings, got: %#v", arr)
		}

		key := arr[1].(enc.BulkString).Val
		val := arr[2]

		storage.Set(key, val)
		return replyOK(conn)

	case "GET":
		if len(arr) != 2 {
			return fmt.Errorf("invalid GET command: must be 2 elements, got: %#v", arr)
		}

		if arr[1].Type() != enc.TypeBulkString {
			return fmt.Errorf("invalid GET command: keymust be a bulk string, got: %#v", arr)
		}

		key := arr[1].(enc.BulkString).Val
		val, ok := storage.Get(key)

		reply := enc.BulkString{Val: val.(enc.BulkString).Val}
		if !ok {
			reply.Null = true
		}

		return reply.Encode(conn)

	default:
		return fmt.Errorf("command not implemented: %s", commandStr)
	}
}

func replyOK(w io.Writer) error {
	responseVal := enc.SimpleString("OK")
	return responseVal.Encode(w)
}

type Storage struct {
	data map[string]enc.Value
}

func (s *Storage) Get(key string) (enc.Value, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *Storage) Set(key string, val enc.Value) {
	s.data[key] = val
}
