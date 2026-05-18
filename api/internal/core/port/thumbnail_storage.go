package port

import "context"

// ThumbnailStorage abstracts thumbnail persistence across backends
// (filesystem, S3/MinIO, etc.). Photo IDs are opaque strings; the
// adapter handles key/path construction internally.
type ThumbnailStorage interface {
	// Get returns thumbnail bytes for a photo at the given size ("s", "m", "l").
	Get(ctx context.Context, photoId, size string) ([]byte, error)

	// Put stores thumbnail bytes for a photo at the given size.
	Put(ctx context.Context, photoId, size string, data []byte) error

	// Exists reports whether a thumbnail exists for the given photo and size.
	Exists(ctx context.Context, photoId, size string) (bool, error)

	// Delete removes all thumbnail sizes for a photo.
	Delete(ctx context.Context, photoId string) error
}
