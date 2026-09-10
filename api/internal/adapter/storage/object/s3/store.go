// Package s3 implements port.ObjectStore against any S3-compatible object
// store (Wasabi, AWS S3, MinIO, Backblaze B2 via the S3 API) using the MinIO
// Go SDK. Keys are always supplied by a workspace-bound accessor, never built
// here (issue #74 §4.3).
package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds connection parameters for an S3-compatible endpoint.
type Config struct {
	// Endpoint is the host[:port] of the S3 API, e.g.
	// "s3.eu-west-1.wasabisys.com" or "minio:9000".
	Endpoint string
	// Region is the signing region (required by Wasabi/AWS; MinIO ignores it).
	Region string
	// AccessKey / SecretKey are the S3 credentials.
	AccessKey string
	SecretKey string
	// Bucket is the bucket all objects live in.
	Bucket string
	// UseSSL selects https.
	UseSSL bool
}

// Store is an S3-backed ObjectStore.
type Store struct {
	client *minio.Client
	bucket string
}

// New connects to the configured endpoint and verifies the bucket exists.
func New(cfg Config) (*Store, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("s3 object store: connect: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("s3 object store: bucket check: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("s3 object store: bucket %q does not exist", cfg.Bucket)
	}

	return &Store{client: client, bucket: cfg.Bucket}, nil
}

// Get returns a reader for the object at key. A missing object is reported as
// a minio ErrorResponse with status 404 from the reader's first read.
func (s *Store) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// Put stores an object.
func (s *Store) Put(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Delete removes an object.
func (s *Store) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

// Exists reports whether an object exists.
func (s *Store) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		var resp minio.ErrorResponse
		if errors.As(err, &resp) && resp.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// List returns object keys under prefix.
func (s *Store) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	for obj := range s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		keys = append(keys, obj.Key)
	}
	return keys, nil
}

// PresignGet returns a time-limited GET URL for the object.
func (s *Store) PresignGet(ctx context.Context, key string, expirySeconds int) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, s.bucket, key, time.Duration(expirySeconds)*time.Second, url.Values{})
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
