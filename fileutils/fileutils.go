package fileutils

import (
	"io"
	"os"
	"path/filepath"

	errors "github.com/pkg/errors"
)

// CreateDirIfMissing makes sure that the dir exists and returns whether the dir is empty
func CreateDirIfMissing(dirPath string) (bool, error) {
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return false, errors.Wrapf(err, "error while creating dir: %s", dirPath)
	}
	if err := SyncParentDir(dirPath); err != nil {
		return false, err
	}
	return DirEmpty(dirPath)
}

// SyncParentDir fsyncs the parent dir of the given path
func SyncParentDir(path string) error {
	return SyncDir(filepath.Dir(path))
}

// SyncDir fsyncs the given dir
func SyncDir(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return errors.Wrapf(err, "error while opening dir:%s", dirPath)
	}
	if err := dir.Sync(); err != nil {
		dir.Close()
		return errors.Wrapf(err, "error while synching dir:%s", dirPath)
	}
	if err := dir.Close(); err != nil {
		return errors.Wrapf(err, "error while closing dir:%s", dirPath)
	}
	return err
}

// DirEmpty returns true if the dir at dirPath is empty
func DirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, errors.Wrapf(err, "error opening dir [%s]", dirPath)
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	err = errors.Wrapf(err, "error checking if dir [%s] is empty", dirPath)
	return false, err
}
