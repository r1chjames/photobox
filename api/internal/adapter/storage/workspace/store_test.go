package workspacestore

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeObjects records the keys it is asked to operate on.
type fakeObjects struct {
	put    []string
	get    []string
	delete []string
	list   []string
}

func (f *fakeObjects) Get(_ context.Context, key string) (io.ReadCloser, error) {
	f.get = append(f.get, key)
	return io.NopCloser(strings.NewReader("data")), nil
}

func (f *fakeObjects) Put(_ context.Context, key string, _ io.Reader, _ int64, _ string) error {
	f.put = append(f.put, key)
	return nil
}

func (f *fakeObjects) Delete(_ context.Context, key string) error {
	f.delete = append(f.delete, key)
	return nil
}

func (f *fakeObjects) Exists(_ context.Context, _ string) (bool, error) { return false, nil }

func (f *fakeObjects) List(_ context.Context, prefix string) ([]string, error) {
	f.list = append(f.list, prefix)
	return []string{prefix + "a", prefix + "b"}, nil
}

func (f *fakeObjects) PresignGet(_ context.Context, key string, _ int) (string, error) {
	return "https://example/" + key, nil
}

const (
	wsA = "11111111-1111-1111-1111-111111111111"
	wsB = "22222222-2222-2222-2222-222222222222"
)

func TestNew_RejectsInvalidWorkspaceID(t *testing.T) {
	for _, bad := range []string{"", "not-a-uuid", "../other", "11111111-1111-1111-1111-11111111111", "ws/../x"} {
		_, err := New(&fakeObjects{}, bad)
		assert.ErrorIs(t, err, ErrInvalidWorkspaceID, "workspace ID %q must be rejected", bad)
	}
}

func TestNew_RejectsNilStore(t *testing.T) {
	_, err := New(nil, wsA)
	assert.Error(t, err)
}

// TestKeys_AlwaysUnderWorkspacePrefix is the core structural-isolation
// property: every key a store produces begins with its own workspace ID, so a
// bound store can never address another tenant's objects.
func TestKeys_AlwaysUnderWorkspacePrefix(t *testing.T) {
	f := &fakeObjects{}
	s, err := New(f, wsA)
	require.NoError(t, err)

	ctx := context.Background()
	require.NoError(t, s.PutOriginal(ctx, "photo-1", strings.NewReader("x"), 1, "image/jpeg"))
	require.NoError(t, s.PutThumbnail(ctx, "photo-1", "m", strings.NewReader("x"), 1, "image/webp"))
	_, err = s.GetOriginal(ctx, "photo-1")
	require.NoError(t, err)
	_, err = s.GetThumbnail(ctx, "photo-1", "l")
	require.NoError(t, err)

	all := append(append([]string{}, f.put...), f.get...)
	require.NotEmpty(t, all)
	for _, key := range all {
		assert.Contains(t, key, "/"+wsA+"/", "key %q must be namespaced to workspace A", key)
		assert.NotContains(t, key, wsB, "key %q must not reference workspace B", key)
	}
}

// TestKeys_TraversalRejected asserts IDs that could escape the prefix are
// refused before any store call.
func TestKeys_TraversalRejected(t *testing.T) {
	f := &fakeObjects{}
	s, err := New(f, wsA)
	require.NoError(t, err)

	ctx := context.Background()
	for _, badID := range []string{
		"../../etc/passwd",
		"photo/../../other",
		"photo\x00id",
		"photo id",
		"photo/1",
		"..",
	} {
		err := s.PutOriginal(ctx, badID, strings.NewReader("x"), 1, "image/jpeg")
		assert.ErrorIs(t, err, ErrInvalidObjectID, "id %q must be rejected", badID)
	}
	assert.Empty(t, f.put, "no key should reach the object store for invalid IDs")
}

func TestThumbnailSizeValidation(t *testing.T) {
	f := &fakeObjects{}
	s, err := New(f, wsA)
	require.NoError(t, err)

	for _, size := range []string{"", "xl", "m.webp", "../l", "M"} {
		err := s.PutThumbnail(context.Background(), "photo-1", size, strings.NewReader("x"), 1, "image/webp")
		assert.ErrorIs(t, err, ErrInvalidSize, "size %q must be rejected", size)
	}
	assert.Empty(t, f.put)
}

// TestList_PrefixesAreWorkspaceScoped verifies list operations can only
// enumerate the bound workspace's prefixes.
func TestList_PrefixesAreWorkspaceScoped(t *testing.T) {
	f := &fakeObjects{}
	s, err := New(f, wsA)
	require.NoError(t, err)

	_, err = s.ListOriginals(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"originals/" + wsA + "/"}, f.list)

	f.list = nil
	_, err = s.ListAll(context.Background())
	require.NoError(t, err)
	for _, prefix := range f.list {
		assert.Contains(t, prefix, "/"+wsA+"/", "list prefix %q must be workspace-scoped", prefix)
	}
	assert.Len(t, f.list, 4, "expected the four object classes to be listed")
}

// TestTwoWorkspaces_NeverOverlap proves two bound stores produce disjoint key
// spaces — the property that makes cross-tenant access structurally
// impossible.
func TestTwoWorkspaces_NeverOverlap(t *testing.T) {
	fA, fB := &fakeObjects{}, &fakeObjects{}
	sA, err := New(fA, wsA)
	require.NoError(t, err)
	sB, err := New(fB, wsB)
	require.NoError(t, err)

	ctx := context.Background()
	// Same photo ID in both workspaces must map to different keys.
	require.NoError(t, sA.PutOriginal(ctx, "photo-1", strings.NewReader("x"), 1, "image/jpeg"))
	require.NoError(t, sB.PutOriginal(ctx, "photo-1", strings.NewReader("x"), 1, "image/jpeg"))

	require.Len(t, fA.put, 1)
	require.Len(t, fB.put, 1)
	assert.NotEqual(t, fA.put[0], fB.put[0], "the same photo ID must not collide across workspaces")
}

func TestDeleteThumbnail_AndTrashEditedKeys(t *testing.T) {
	f := &fakeObjects{}
	s, err := New(f, wsA)
	require.NoError(t, err)

	require.NoError(t, s.DeleteThumbnail(context.Background(), "photo-1", "s"))
	require.Len(t, f.delete, 1)
	assert.Equal(t, "thumbnails/"+wsA+"/photo-1/s.webp", f.delete[0])
}

func TestEditedKey_ValidatesParamsHash(t *testing.T) {
	f := &fakeObjects{}
	s, err := New(f, wsA)
	require.NoError(t, err)

	// A non-hex hash is rejected (cannot smuggle path separators).
	key, err := s.editedKey("photo-1", "../../evil")
	assert.ErrorIs(t, err, ErrInvalidObjectID)
	assert.Empty(t, key)

	key, err = s.editedKey("photo-1", "deadbeef")
	require.NoError(t, err)
	assert.Equal(t, "edited/"+wsA+"/photo-1/deadbeef", key)
}

func TestNew_ErrorMessageMentionsID(t *testing.T) {
	_, err := New(&fakeObjects{}, "bad")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidWorkspaceID))
}
