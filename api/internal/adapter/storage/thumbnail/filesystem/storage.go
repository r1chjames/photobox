package filesystem

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Storage stores thumbnails on the local filesystem under a configurable
// base directory, organised as baseDir/{size}/{safeId}.webp.
type Storage struct {
	baseDir string
}

// New creates a filesystem thumbnail storage.
//
// Thumbnails are stored under baseDir. The caller is responsible for
// ensuring the directory exists and is writable.  A separate writable
// volume (e.g. a PVC or host path) is recommended in containerised
// deployments, distinct from the read-only photo mount.
func New(baseDir string) *Storage {
	return &Storage{baseDir: baseDir}
}

func (s *Storage) Get(_ context.Context, photoId, size string) ([]byte, error) {
	path := s.path(photoId, size)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil // not-found is not an error for the caller
		}
		return nil, err
	}
	return data, nil
}

func (s *Storage) Put(_ context.Context, photoId, size string, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	path := s.path(photoId, size)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *Storage) Exists(_ context.Context, photoId, size string) (bool, error) {
	_, err := os.Stat(s.path(photoId, size))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Storage) Delete(_ context.Context, photoId string) error {
	// Remove all sizes for this photo
	for _, size := range []string{"s", "m", "l"} {
		path := s.path(photoId, size)
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}

// path returns the filesystem path for a photo's thumbnail at a given size.
func (s *Storage) path(photoId, size string) string {
	safeId := strings.ReplaceAll(photoId, "/", "_")
	safeId = strings.ReplaceAll(safeId, "+", "-")
	safeId = strings.ReplaceAll(safeId, "=", "")
	return filepath.Join(s.baseDir, size, safeId+".webp")
}
