package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds the connection parameters for an S3-compatible object store
// (MinIO, AWS S3, Cloudflare R2, etc.).
type Config struct {
	Endpoint  string // e.g. "minio:9000" or "s3.amazonaws.com"
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// Storage stores thumbnails in an S3-compatible object store using the
// MinIO Go SDK.  Object keys follow the pattern:
//
//	thumbnails/{safeId}/{size}.jpg
type Storage struct {
	client *minio.Client
	bucket string
}

// New creates an S3/MinIO thumbnail storage.  It connects to the
// configured endpoint and verifies the bucket exists, creating it if
// necessary (MinIO auto-creates; AWS requires the bucket to exist).
func New(cfg Config) (*Storage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("s3 thumbnail storage: connect: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("s3 thumbnail storage: bucket check: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("s3 thumbnail storage: create bucket: %w", err)
		}
	}

	return &Storage{client: client, bucket: cfg.Bucket}, nil
}

func (s *Storage) Get(ctx context.Context, photoId, size string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, s.key(photoId, size), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, obj); err != nil {
		// If the object doesn't exist, minio-go returns an error response;
		// treat that as a clean miss.
		var resp minio.ErrorResponse
		if errors.As(err, &resp) && resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Storage) Put(ctx context.Context, photoId, size string, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(ctx, s.bucket, s.key(photoId, size), reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	return err
}

func (s *Storage) Exists(ctx context.Context, photoId, size string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, s.key(photoId, size), minio.StatObjectOptions{})
	if err != nil {
		var resp minio.ErrorResponse
		if errors.As(err, &resp) && resp.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Storage) Delete(ctx context.Context, photoId string) error {
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for _, size := range []string{"s", "m", "l"} {
			objectsCh <- minio.ObjectInfo{
				Key: s.key(photoId, size),
			}
		}
	}()

	for err := range s.client.RemoveObjects(ctx, s.bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		if err.Err != nil {
			return err.Err
		}
	}
	return nil
}

// key returns the S3 object key for a photo's thumbnail at a given size.
func (s *Storage) key(photoId, size string) string {
	safeId := strings.ReplaceAll(photoId, "/", "_")
	safeId = strings.ReplaceAll(safeId, "+", "-")
	safeId = strings.ReplaceAll(safeId, "=", "")
	return fmt.Sprintf("thumbnails/%s/%s.jpg", safeId, size)
}
