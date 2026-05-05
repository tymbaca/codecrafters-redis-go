package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ensureManifestFile(svc *Service) error {
	path := filepath.Join(svc.dir, svc.appendDir, svc.appendFile+".manifest")

	_, _ = os.Create(path)
	f, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		_, err = fmt.Fprintf(f, "file %s.1.incr.aof seq 1 type i", svc.appendFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func ensureDirCreated(dir string) {
	_ = os.MkdirAll(dir, 0o755)
}

func ensureFileCreated(path string) {
	_, _ = os.Create(path)
}
