package aof

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/redis-starter-go/pkg/counting"
	"github.com/codecrafters-io/redis-starter-go/pkg/enc"
	"github.com/codecrafters-io/redis-starter-go/pkg/manifest"
)

const limit = 100 * 1024 * 1024

func New(cwd, aofDirname, aofFilename string) (*AOF, error) {
	aof := &AOF{
		cwd:            cwd,
		aofDirname:     aofDirname,
		aofFilename:    aofFilename,
		currentFileNum: 0,
		currentFile:    nil,
		offset:         0,
	}

	manifestPath := filepath.Join(cwd, aofDirname, aofFilename+".manifest")
	manifest, manifestFile, err := readOrCreateManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	aof.manifest = manifest
	aof.manifestFile = manifestFile

	if len(aof.manifest.Records) == 0 {
		err := createNewFile(aof)
		if err != nil {
			return nil, err
		}
	} else {
		
	}

	return aof, nil
}

func readOrCreateManifest(path string) (manifest.Manifest, *os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if os.IsNotExist(err) {
		f, err = os.Create(path)
		if err != nil {
			return manifest.Manifest{}, nil, err
		}
	} else if err != nil {
		return err
	}
	defer f.Close()

	mft, err := manifest.Decode(f)
	if err != nil {
		return manifest.Manifest{}, nil, err
	}

	return nil
}

func ensureAofFiles(cwd, aofDirname, aofFilename string) error {
	ensureDirCreated(filepath.Join(cwd, aofDirname))
	ensureFileCreated())
	if err := ensureManifestFile(svc); err != nil {
		return err
	}

	if svc.aof == nil {
		aof, err := aof.New(svc.dir, svc.appendDir, svc.appendFile)
		if err != nil {
			return fmt.Errorf("init AOF: %w", err)
		}
		svc.aof = aof
	}

	return replayAof(svc)
}

func replayAof(svc *Service) error {
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

func (a *AOF) Close() error {
	return a.currentFile.Close()
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

func ensureManifestFile() error {
}

func ensureDirCreated(dir string) {
	_ = os.MkdirAll(dir, 0o755)
}

func ensureFileCreated(path string) {
	_, _ = os.Create(path)
}
