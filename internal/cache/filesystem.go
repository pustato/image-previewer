package cache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var _ Filesystem = (*filesystem)(nil)

var ErrFileNotExists = errors.New("file is not exists")

const filePermission = 0o700

type filesystem struct {
	basePath string
}

func NewFilesystem(basePath string) (Filesystem, error) {
	if err := ensureDir(basePath); err != nil {
		return nil, fmt.Errorf("new filesystem: %w", err)
	}

	return &filesystem{
		basePath: basePath,
	}, nil
}

func (f *filesystem) WriteFile(name string, content []byte) error {
	path := filepath.Join(f.basePath, name)

	if err := os.WriteFile(path, content, filePermission); err != nil {
		return fmt.Errorf("filesystem write file %s: %w", path, err)
	}

	return nil
}

func (f *filesystem) ReadFile(name string) ([]byte, error) {
	path := filepath.Join(f.basePath, name)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotExists
		}

		return nil, fmt.Errorf("filesystem read file %s: %w", path, err)
	}

	return content, nil
}

func (f *filesystem) RemoveFile(name string) error {
	path := filepath.Join(f.basePath, name)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove file %s: %w", path, err)
	}

	return nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(dir, filePermission); err != nil {
		return fmt.Errorf("ensure dir: %w", err)
	}

	return nil
}
