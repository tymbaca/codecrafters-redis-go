package redis

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

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
			defer conn.Close()

			err := handleConn(ctx, conn)
			if err != nil {
				slog.Error("handle conn failed", "err", err)
			}
		}()
	}
}

func handleConn(ctx context.Context, conn io.ReadWriter) error {
	for {
		command, ok, err := readCommand(ctx, conn)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		err = handleCommand(ctx, conn, command)
		if err != nil {
			return err
		}
	}
}

func readCommand(ctx context.Context, conn io.Reader) (Command, bool, error) {
	panic("not implemented")
}

func readArray(ctx context.Context, conn io.Reader) enc.Array {
	panic("not implemented")
}

func handleCommand(ctx context.Context, conn io.Reader, command Command) error {
	panic("not implemented")
}

type Command struct{}
