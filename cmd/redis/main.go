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

	"github.com/codecrafters-io/redis-starter-go/pkg/assert"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func Run() error {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	context.AfterFunc(ctx, func() {
		_ = l.Close()
	})

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("accept conn failed", "err", err)
		}

		go func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			context.AfterFunc(ctx, func() {
				_ = conn.Close()
			})

			err := handleConn(conn)
			if err != nil {
				slog.Error("handle conn failed", "err", err)
			}
		}()
	}
}

func handleConn(conn io.ReadWriter) error {
	for {
		command, ok, err := readCommand(conn)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		err = handleCommand(conn, command)
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

func handleCommand(storage Storage, conn io.Writer, command enc.Value) error {
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

		responseVal := arr[1]
		return responseVal.Encode(conn)

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
