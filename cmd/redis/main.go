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

func handleCommand(conn io.Writer, command enc.Value) error {
	if command.Type() != enc.TypeArray {
		return fmt.Errorf("exected array, got: %s (of type %s)", command.String(), command.Type())
	}

	arr := command.(enc.Array)
	assert.True(len(arr) > 0)
	assert.True(arr[0].Type() == enc.TypeBulkString)

	commandStr := arr[0].(enc.BulkString)

	if commandStr == "PING" {
		responseVal := enc.SimpleString("PONG")
		return responseVal.Encode(conn)
	}

	panic("not implemented")
}
