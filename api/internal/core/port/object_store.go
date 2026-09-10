package port

import (
	"context"
	"io"
)

// ObjectStore is the process-level, full-access object storage client. It is
// deliberately low-level and its key-building is owned by WorkspaceStore:
// application code never builds object keys directly (issue #74 §4.3).
type ObjectStore interface {
	// Get returns a reader for an object. The caller closes it.
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Put stores an object of the given size.
	Put(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	// Delete removes an object.
	Delete(ctx context.Context, key string) error
	// Exists reports whether an object exists.
	Exists(ctx context.Context, key string) (bool, error)
	// List returns object keys under a prefix (used by the orphan sweep).
	List(ctx context.Context, prefix string) ([]string, error)
	// PresignGet returns a time-limited URL for direct download (the
	// config-flip escape hatch for originals, issue #74 D7).
	PresignGet(ctx context.Context, key string, expirySeconds int) (string, error)
}
