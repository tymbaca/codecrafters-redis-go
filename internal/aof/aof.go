package aof

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/redis-starter-go/pkg/counting"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/manifest"
)

const limit = 100 * 1024 * 1024

func New(cwd, aofDirname, aofFilename string) (*AOF, error) {
	path := filepath.Join(cwd, aofDirname, aofFilename+".1.incr.aof")

	f, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}

	return &AOF{
		cwd:            cwd,
		aofDirname:     aofDirname,
		aofFilename:    aofFilename,
		currentFileNum: 1,
		currentFile:    f,
		offset:         0,
	}, nil
}

type AOF struct {
	cwd, aofDirname, aofFilename string

	currentFileNum int
	currentFile    *os.File
	offset         int

	manifest     manifest.Manifest
	manifestFile *os.File
}

func (a *AOF) Append(ctx context.Context, cmd enc.Value) error {
	w := counting.NewWriter(a.currentFile)

	err := cmd.Encode(w)
	a.offset += w.N
	if err != nil {
		return err
	}

	if a.offset >= limit {
		err = createNewFile(a)
		if err != nil {
			return err
		}
	}

	return nil
}

func createNewFile(a *AOF) error {
	a.offset = 0
	a.currentFileNum++

	err := a.currentFile.Close()
	if err != nil {
		return err
	}

	filename := filename(a)
	path := filepath.Join(a.cwd, a.aofDirname, filename)

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	a.currentFile = f

	a.manifest.Records = append(a.manifest.Records, manifest.Record{
		File: filename,
		Seq:  a.currentFileNum,
		Type: "i",
	})

	err = a.manifestFile.Truncate(0)
	if err != nil {
		return err
	}

	err = a.manifest.Encode(a.manifestFile)
	if err != nil {
		return err
	}

	return nil
}

func filename(a *AOF) string {
	return fmt.Sprintf("%s.%d.incr.aof", a.aofFilename, a.currentFileNum)
}

func (a *AOF) Close() error {
	return a.currentFile.Close()
}
