// Package workspacestore provides a prefix-bound object accessor: every key
// is built inside the struct from a workspace ID bound at construction, so
// cross-tenant object access is structurally impossible in-process (issue
// #74 §4.3, the in-process analog of AWS Cognito scoped STS credentials).
package workspacestore

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"

	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

// ID validation. Workspace and photo IDs must be UUID-shaped: no slashes, no
// dots, nothing that can escape a prefix. This is the boundary check that
// keeps a bad ID from ever reaching the key builder.
var (
	uuidRe    = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	plainIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	// editHashRe allows the lowercase hex digest used in edited-object keys.
	editHashRe = regexp.MustCompile(`^[0-9a-f]{8,64}$`)
	sizeRe     = regexp.MustCompile(`^(s|m|l)$`)
)

var (
	// ErrInvalidWorkspaceID is returned when the bound workspace ID is not a
	// UUID. Construction fails closed.
	ErrInvalidWorkspaceID = errors.New("invalid workspace id")
	// ErrInvalidObjectID is returned when a photo/resource ID contains
	// characters that could escape the workspace prefix.
	ErrInvalidObjectID = errors.New("invalid object id")
	// ErrInvalidSize is returned for an unsupported thumbnail size.
	ErrInvalidSize = errors.New("invalid thumbnail size")
)

// Store is an object accessor bound to exactly one workspace. All keys it
// produces are under a prefix derived from that workspace ID.
type Store struct {
	objects port.ObjectStore
	wsID    string
}

// New binds an ObjectStore to a workspace. The workspace ID is validated as a
// UUID; an invalid ID fails closed rather than producing an unbound store.
func New(objects port.ObjectStore, workspaceID string) (*Store, error) {
	if !uuidRe.MatchString(workspaceID) {
		return nil, fmt.Errorf("%w: %q", ErrInvalidWorkspaceID, workspaceID)
	}
	if objects == nil {
		return nil, errors.New("workspacestore: nil object store")
	}
	return &Store{objects: objects, wsID: workspaceID}, nil
}

// WorkspaceID returns the bound workspace ID.
func (s *Store) WorkspaceID() string { return s.wsID }

// originalsKey builds originals/{ws}/{photoID}. photoID must be
// boundary-safe (UUID or base64url-ish id without path separators).
func (s *Store) originalsKey(photoID string) (string, error) {
	if !plainIDRe.MatchString(photoID) {
		return "", fmt.Errorf("%w: %q", ErrInvalidObjectID, photoID)
	}
	return fmt.Sprintf("originals/%s/%s", s.wsID, photoID), nil
}

// thumbnailsKey builds thumbnails/{ws}/{photoID}/{size}.webp.
func (s *Store) thumbnailsKey(photoID, size string) (string, error) {
	if !plainIDRe.MatchString(photoID) {
		return "", fmt.Errorf("%w: %q", ErrInvalidObjectID, photoID)
	}
	if !sizeRe.MatchString(size) {
		return "", fmt.Errorf("%w: %q", ErrInvalidSize, size)
	}
	return fmt.Sprintf("thumbnails/%s/%s/%s.webp", s.wsID, photoID, size), nil
}

// trashKey builds trash/{ws}/{photoID}.
func (s *Store) trashKey(photoID string) (string, error) {
	if !plainIDRe.MatchString(photoID) {
		return "", fmt.Errorf("%w: %q", ErrInvalidObjectID, photoID)
	}
	return fmt.Sprintf("trash/%s/%s", s.wsID, photoID), nil
}

// editedKey builds edited/{ws}/{photoID}/{paramsHash}.
func (s *Store) editedKey(photoID, paramsHash string) (string, error) {
	if !plainIDRe.MatchString(photoID) {
		return "", fmt.Errorf("%w: %q", ErrInvalidObjectID, photoID)
	}
	if !editHashRe.MatchString(paramsHash) {
		return "", fmt.Errorf("%w: %q", ErrInvalidObjectID, paramsHash)
	}
	return fmt.Sprintf("edited/%s/%s/%s", s.wsID, photoID, paramsHash), nil
}

// PutOriginal stores a photo's original bytes.
func (s *Store) PutOriginal(ctx context.Context, photoID string, reader io.Reader, size int64, contentType string) error {
	key, err := s.originalsKey(photoID)
	if err != nil {
		return err
	}
	return s.objects.Put(ctx, key, reader, size, contentType)
}

// GetOriginal returns a reader for a photo's original bytes.
func (s *Store) GetOriginal(ctx context.Context, photoID string) (io.ReadCloser, error) {
	key, err := s.originalsKey(photoID)
	if err != nil {
		return nil, err
	}
	return s.objects.Get(ctx, key)
}

// DeleteOriginal removes a photo's original bytes.
func (s *Store) DeleteOriginal(ctx context.Context, photoID string) error {
	key, err := s.originalsKey(photoID)
	if err != nil {
		return err
	}
	return s.objects.Delete(ctx, key)
}

// PutThumbnail stores a thumbnail at the given size.
func (s *Store) PutThumbnail(ctx context.Context, photoID, size string, reader io.Reader, sizeBytes int64, contentType string) error {
	key, err := s.thumbnailsKey(photoID, size)
	if err != nil {
		return err
	}
	return s.objects.Put(ctx, key, reader, sizeBytes, contentType)
}

// GetThumbnail returns a reader for a photo's thumbnail at the given size.
func (s *Store) GetThumbnail(ctx context.Context, photoID, size string) (io.ReadCloser, error) {
	key, err := s.thumbnailsKey(photoID, size)
	if err != nil {
		return nil, err
	}
	return s.objects.Get(ctx, key)
}

// DeleteThumbnail removes a photo's thumbnail at the given size.
func (s *Store) DeleteThumbnail(ctx context.Context, photoID, size string) error {
	key, err := s.thumbnailsKey(photoID, size)
	if err != nil {
		return err
	}
	return s.objects.Delete(ctx, key)
}

// PresignOriginal returns a time-limited direct-download URL for an original.
func (s *Store) PresignOriginal(ctx context.Context, photoID string, expirySeconds int) (string, error) {
	key, err := s.originalsKey(photoID)
	if err != nil {
		return "", err
	}
	return s.objects.PresignGet(ctx, key, expirySeconds)
}

// ListOriginals lists only this workspace's original keys. The prefix is
// derived from the bound workspace, so a sweep can never touch another
// tenant's objects.
func (s *Store) ListOriginals(ctx context.Context) ([]string, error) {
	return s.objects.List(ctx, fmt.Sprintf("originals/%s/", s.wsID))
}

// ListAll lists every object belonging to this workspace across the four
// top-level prefixes (originals, thumbnails, trash, edited). Objects live at
// `{class}/{workspace}/...`, so there is no single common prefix; each class
// is listed with the bound workspace ID. Used by the per-workspace orphan
// sweep and by workspace deletion.
func (s *Store) ListAll(ctx context.Context) ([]string, error) {
	var all []string
	for _, class := range []string{"originals", "thumbnails", "trash", "edited"} {
		keys, err := s.objects.List(ctx, fmt.Sprintf("%s/%s/", class, s.wsID))
		if err != nil {
			return nil, err
		}
		all = append(all, keys...)
	}
	return all, nil
}

// DeleteAll removes every object belonging to this workspace. Used when a
// workspace is deleted; the workspace prefix binding guarantees only this
// workspace's objects are ever enumerated.
func (s *Store) DeleteAll(ctx context.Context) error {
	keys, err := s.ListAll(ctx)
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := s.objects.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}
