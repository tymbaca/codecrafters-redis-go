package aof

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/codecrafters-io/redis-starter-go/pkg/command"
	"github.com/codecrafters-io/redis-starter-go/pkg/counting"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/manifest"
	"github.com/samber/lo"
)

const limit = 100 * 1024 * 1024

func New(ctx context.Context, cwd, aofDirname, aofFilename string, replay func(ctx context.Context, cmd command.Command) error) (*AOF, error) {
	aof := &AOF{
		cwd:            cwd,
		aofDirname:     aofDirname,
		aofFilename:    aofFilename,
		currentFileNum: 0,
		currentFile:    nil,
		offset:         0,
	}

	root, err := os.OpenRoot(cwd)
	if err != nil {
		return nil, err
	}
	aof.root = root

	err = root.MkdirAll(aofDirname, 0o777)
	if err != nil {
		return nil, err
	}

	slog.Debug("open or create manifest file", "path", filepath.Join(cwd, aofDirname, aofFilename+".manifest"))
	mft, mftFile, err := readOrCreateManifest(root, aofFilename+".manifest")
	if err != nil {
		return nil, err
	}
	aof.manifest = mft
	aof.manifestFile = mftFile

	if len(aof.manifest.Records) == 0 {
		slog.Debug("no records in manifest, creating first AOF file")
		err := createNewFile(aof)
		if err != nil {
			return nil, err
		}
	} else {
		slog.Debug("found records in manifest", "records", aof.manifest.Records)

		toReadRecords := lo.Filter(aof.manifest.Records, func(r manifest.Record, _ int) bool { return r.Type == "i" })
		slices.SortFunc(toReadRecords, func(a, b manifest.Record) int { return cmp.Compare(a.Seq, b.Seq) })

		cmdCtx := command.Context{}

		for i, rec := range toReadRecords {

			f, err := os.OpenFile(filepath.Join(cwd, aof.aofDirname, rec.File), os.O_RDWR|os.O_APPEND|os.O_SYNC, 0o666)
			if err != nil {
				return nil, err
			}

			if replay != nil {
				slog.Debug("replaying file", "file", rec.File, "seq", rec.Seq)

				for {
					cmdVal, ok, err := readCommand(f)
					if err != nil {
						return nil, fmt.Errorf("read command: %w", err)
					}

					if !ok {
						break
					}

					cmd, err := command.Parse(cmdCtx, cmdVal)
					if err != nil {
						return nil, fmt.Errorf("parse command: %w", err)
					}

					slog.Debug("replaying command", "cmd", cmd, "file", rec.File, "seq", rec.Seq)
					err = replay(ctx, cmd)
					if err != nil {
						return nil, fmt.Errorf("replay command: %w", err)
					}
				}
			} else {
				slog.Debug("skipping replay of file", "file", rec.File, "seq", rec.Seq)
			}

			// set last file as current
			if i == len(toReadRecords)-1 {
				slog.Debug("file set as current", "file", rec.File, "seq", rec.Seq)
				aof.currentFileNum = rec.Seq
				aof.currentFile = f
			} else {
				err := f.Close()
				if err != nil {
					return nil, fmt.Errorf("close file: %w", err)
				}
			}
		}
	}

	return aof, nil
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

func readOrCreateManifest(root *os.Root, name string) (manifest.Manifest, *os.File, error) {
	f, err := root.OpenFile(name, os.O_RDWR|os.O_SYNC, 0o666)
	if os.IsNotExist(err) {
		f, err = root.Create(name)
	}
	if err != nil {
		return manifest.Manifest{}, nil, err
	}

	mft, err := manifest.Decode(f)
	if err != nil {
		return manifest.Manifest{}, nil, err
	}

	return mft, f, nil
}

type AOF struct {
	root *os.Root

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

func (a *AOF) Close() error {
	return errors.Join(
		a.currentFile.Close(),
		a.manifestFile.Close(),
	)
}

func createNewFile(a *AOF) error {
	a.offset = 0
	a.currentFileNum++

	if a.currentFile != nil {
		err := a.currentFile.Close()
		if err != nil {
			return err
		}
	}

	filename := filename(a)
	path := filepath.Join(a.cwd, a.aofDirname, filename)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0o666)
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
