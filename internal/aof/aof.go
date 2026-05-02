package aof

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/redis-starter-go/pkg/counting"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
)

func New(cwd, aofDirname, aofFilename string) (*AOF, error) {
	path := filepath.Join(cwd, aofDirname, aofFilename+".1.incr.aof")

	f, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}

	return &AOF{
		cwd:         cwd,
		aofDirname:  aofDirname,
		aofFilename: aofFilename,
		currentFile: f,
		offset:      0,
	}, nil
}

type AOF struct {
	cwd, aofDirname, aofFilename string

	currentFile io.WriteCloser
	offset      int
}

func (a *AOF) Append(ctx context.Context, cmd enc.Value) error {
	w := counting.NewWriter(a.currentFile)

	err := cmd.Encode(w)
	a.offset += w.N

	return err
}

func (a *AOF) Close() error {
	return a.currentFile.Close()
}
