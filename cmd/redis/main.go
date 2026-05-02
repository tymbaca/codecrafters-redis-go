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
	"strings"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/internal/service"
	"github.com/codecrafters-io/redis-starter-go/pkg/assert"
	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/google/uuid"
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

	service, err := service.New(service.Options{
		Dir:            *dirFlag,
		AppendOnly:     *appendOnlyFlag == "yes",
		AppendDir:      *appendDirnameFlag,
		AppendFilename: *appendFilenameFlag,
		AppendFsync:    *appendFsyncFlag,
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
	connCtx := command.Context{
		ConnID: uuid.NewString(),
	}

	defer svc.Exec(ctx, command.Discard{Context: connCtx})

	for {
		command, ok, err := readCommand(conn)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		err = handleCommand(ctx, conn, connCtx, svc, command)
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

func handleCommand(ctx context.Context, conn io.Writer, connCtx command.Context, svc *service.Service, commandVal enc.Value) error {
	if commandVal.Type() != enc.TypeArray {
		return fmt.Errorf("exected array, got: %s (of type %s)", commandVal.String(), commandVal.Type())
	}

	arr := commandVal.(enc.Array)
	assert.True(len(arr) > 0)
	assert.True(arr[0].Type() == enc.TypeBulkString)
	commandStr := arr[0].(enc.BulkString)

	args, err := argsToString(arr[1:])
	if err != nil {
		return fmt.Errorf("parse command args (%#v): %w", arr, err)
	}

	switch strings.ToUpper(commandStr.Val) {
	case "PING":
		responseVal := enc.SimpleString("PONG")
		return responseVal.Encode(conn)

	case "ECHO":
		if len(arr) != 2 {
			return fmt.Errorf("invalid ECHO command: must be 1 arg, got: %#v", args)
		}

		reply := arr[1]
		return reply.Encode(conn)

	case "GET":
		if len(args) < 1 {
			return fmt.Errorf("invalid GET command: must be 1 or more args, got: %#v", args)
		}

		cmd, err := command.ParseGet(connCtx, args)
		if err != nil {
			return replyError(conn, fmt.Errorf("parse GET: %w", err))
		}

		reply, err := svc.Exec(ctx, cmd)
		if err != nil {
			return fmt.Errorf("exec GET: %w", err)
		}

		return reply.Encode(conn)

	case "SET":
		if len(args) < 2 {
			return fmt.Errorf("invalid SET command: must be 3 or more args, got: %#v", args)
		}

		cmd, err := command.ParseSet(connCtx, args)
		if err != nil {
			return replyError(conn, fmt.Errorf("parse SET: %w", err))
		}

		reply, err := svc.Exec(ctx, cmd)
		if err != nil {
			return fmt.Errorf("exec SET: %w", err)
		}

		return reply.Encode(conn)

	case "INCR":
		if len(args) < 1 {
			return fmt.Errorf("invalid INCR command: must be 1 arg, got: %#v", args)
		}

		cmd, err := command.ParseIncr(connCtx, args)
		if err != nil {
			return replyError(conn, fmt.Errorf("parse INCR: %w", err))
		}

		reply, err := svc.Exec(ctx, cmd)
		if err != nil {
			return fmt.Errorf("exec INCR: %w", err)
		}

		return reply.Encode(conn)

	case "MULTI":
		reply, err := svc.Exec(ctx, command.Multi{Context: connCtx})
		if err != nil {
			return fmt.Errorf("exec MULTI: %w", err)
		}

		return reply.Encode(conn)

	case "EXEC":
		reply, err := svc.Exec(ctx, command.Exec{Context: connCtx})
		if err != nil {
			return fmt.Errorf("exec EXEC: %w", err)
		}

		return reply.Encode(conn)

	case "DISCARD":
		reply, err := svc.Exec(ctx, command.Discard{Context: connCtx})
		if err != nil {
			return fmt.Errorf("exec DISCARD: %w", err)
		}

		return reply.Encode(conn)

	case "CONFIG":
		cmd, err := command.ParseConfig(connCtx, args)
		if err != nil {
			return replyError(conn, fmt.Errorf("parse CONFIG: %w", err))
		}

		reply, err := svc.Exec(ctx, cmd)
		if err != nil {
			return fmt.Errorf("exec CONFIG: %w", err)
		}

		return reply.Encode(conn)

	default:
		return replyError(conn, enc.ErrUnknownCommand(commandStr.Val))
	}
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

func replyError(conn io.Writer, err error) error {
	reply := enc.SimpleError(err.Error())

	return reply.Encode(conn)
}
