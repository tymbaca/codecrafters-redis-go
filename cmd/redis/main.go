package redis

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/internal/service"
	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

var (
	dirFlag            = flag.String("dir", "", "current working directory")
	appendOnlyFlag     = flag.String("appendonly", "no", "enables appendonly mode")
	appendDirnameFlag  = flag.String("appenddirname", "", "appendonly directory name")
	appendFilenameFlag = flag.String("appendfilename", "", "appendonly file name")
	appendFsyncFlag    = flag.String("appendfsync", "", "appendonly fsync frequency")
)

func Run() error {
	flag.Parse()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// Level: slog.LevelDebug,
		Level: slog.LevelInfo,
	})))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	cancelAfter := context.AfterFunc(ctx, func() {
		_ = l.Close()
	})
	defer cancelAfter()

	service, err := service.New(ctx, service.Options{
		Dir:         *dirFlag,
		AppendOnly:  *appendOnlyFlag == "yes",
		AppendDir:   *appendDirnameFlag,
		AppendFile:  *appendFilenameFlag,
		AppendFsync: *appendFsyncFlag,
	})
	if err != nil {
		return err
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

			err := handleConn(ctx, conn, service)
			if err != nil {
				slog.Error("handle conn failed", "err", err)
			}
		}()
	}
}

func handleConn(ctx context.Context, conn io.ReadWriter, svc *service.Service) error {
	cmdCtx := command.NewContext(conn)
	defer svc.CloseConn(ctx, cmdCtx)

	for {
		command, ok, err := readCommand(conn)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		err = handleCommand(ctx, conn, cmdCtx, svc, command)
		if err != nil {
			return err
		}
	}
}

func readCommand(conn io.Reader) (enc.Value, bool, error) {
	val, err := enc.Decode(conn)
	if errors.Is(err, io.EOF) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return val, true, nil
}

func handleCommand(ctx context.Context, conn io.Writer, cmdCtx command.Context, svc *service.Service, commandVal enc.Value) error {
	err := svc.Intercept(ctx, cmdCtx, commandVal)
	if err != nil {
		return fmt.Errorf("intercept command value: %w", err)
	}

	cmd, err := command.Parse(cmdCtx, commandVal)
	if err != nil {
		return replyError(conn, err)
	}

	reply, err := svc.Exec(ctx, cmd)
	if err != nil {
		return fmt.Errorf("exec GET: %w", err)
	}

	if reply == nil {
		return nil
	}

	return reply.Encode(conn)
}

func replyError(conn io.Writer, err error) error {
	reply := enc.SimpleError(err.Error())

	return reply.Encode(conn)
}
