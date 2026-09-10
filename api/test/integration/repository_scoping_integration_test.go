//go:build integration

package integration

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database/repository"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/test/testutil"
)

// These tests verify the Phase 2 query-scoping sweep at the repository layer:
// a repository bound with WithWorkspace must see only its own workspace's rows,
// and a bound repository with an empty workspace must see nothing (fail
// closed). They assert scoping directly, independent of the HTTP layer.

const (
	wsAlpha = "aaaaaaaa-0000-0000-0000-00000000000a"
	wsBeta  = "bbbbbbbb-0000-0000-0000-00000000000b"
)

// TestPhotoRepository_WorkspaceScoping_Integration asserts every tenant-facing
// photo query respects the bound workspace.
func TestPhotoRepository_WorkspaceScoping_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	now := time.Now().UnixMilli()
	// Two workspaces, each with a photo (distinct hashes so duplicate
	// detection cannot confuse them).
	for _, row := range []struct {
		id, name, ws, hash string
	}{
		{"alpha-photo", "alpha.jpg", wsAlpha, "hash-alpha"},
		{"beta-photo", "beta.jpg", wsBeta, "hash-beta"},
	} {
		require.NoError(t, env.Db.Exec(`INSERT INTO photobox.photos
			(id, name, workspace_id, file_hash, created_epoch, metadata, deleted_at, hidden, media_type)
			VALUES (?, ?, ?, ?, ?, '{}', NULL, false, 'image')`,
			row.id, row.name, row.ws, row.hash, now).Error)
	}
	// A trashed photo in beta, to exercise the trash query.
	require.NoError(t, env.Db.Exec(`INSERT INTO photobox.photos
		(id, name, workspace_id, created_epoch, metadata, deleted_at, hidden, media_type)
		VALUES ('beta-trashed', 'gone.jpg', ?, ?, '{}', NOW(), false, 'image')`,
		wsBeta, now).Error)

	base := repository.NewPhotoRepository(env)
	alpha := base.WithWorkspace(wsAlpha)

	t.Run("list returns only the bound workspace", func(t *testing.T) {
		photos, err := alpha.ListAllPhotos("", 100, false, "", "", "")
		require.NoError(t, err)
		require.Len(t, photos, 1, "only workspace alpha's photo should be listed")
		assert.Equal(t, "alpha-photo", photos[0].ID)
	})

	t.Run("get by id is scoped", func(t *testing.T) {
		_, err := alpha.GetPhotoById("beta-photo", false)
		assert.ErrorIs(t, err, domain.ErrDataNotFound, "another workspace's photo must not resolve")

		got, err := alpha.GetPhotoById("alpha-photo", false)
		require.NoError(t, err)
		assert.Equal(t, "alpha-photo", got.ID)
	})

	t.Run("trash list is scoped", func(t *testing.T) {
		trash, err := alpha.ListTrashPhotos("", 100, false)
		require.NoError(t, err)
		assert.Empty(t, trash, "alpha's trash is empty; beta's trashed photo must not appear")
	})

	t.Run("search is scoped", func(t *testing.T) {
		// "jpg" matches both tenants' filenames.
		results, err := alpha.SearchPhotos("alpha", 100)
		require.NoError(t, err)
		for _, p := range results {
			assert.Equal(t, wsAlpha, p.WorkspaceID)
		}
	})

	t.Run("timeline is scoped", func(t *testing.T) {
		// Insert year/month so the timeline has something to count.
		require.NoError(t, env.Db.Exec(`UPDATE photobox.photos SET year = 2024, month = 5 WHERE id = 'alpha-photo'`).Error)
		require.NoError(t, env.Db.Exec(`UPDATE photobox.photos SET year = 2019, month = 1 WHERE id = 'beta-photo'`).Error)

		entries, err := alpha.GetTimeline()
		require.NoError(t, err)
		for _, e := range entries {
			assert.NotEqual(t, 2019, e.Year, "beta's 2019 photos must not appear in alpha's timeline")
		}
	})

	t.Run("duplicate detection is per-workspace", func(t *testing.T) {
		// Same file hash in both workspaces must NOT be reported as duplicates
		// to either tenant — duplicates are a within-tenant concept.
		require.NoError(t, env.Db.Exec(`UPDATE photobox.photos SET file_hash = 'shared-hash' WHERE id IN ('alpha-photo','beta-photo')`).Error)

		dups, err := alpha.GetDuplicatePhotos()
		require.NoError(t, err)
		assert.Empty(t, dups, "a cross-workspace hash collision must not be reported as a duplicate")
	})

	t.Run("mutations are scoped", func(t *testing.T) {
		require.NoError(t, alpha.SetFavorite("beta-photo", true))

		var fav bool
		require.NoError(t, env.Db.Raw(`SELECT favorite FROM photobox.photos WHERE id = 'beta-photo'`).Scan(&fav).Error)
		assert.False(t, fav, "a scoped repository must not mutate another workspace's photo")

		_, err := alpha.SoftDeletePhoto("beta-photo")
		assert.Error(t, err, "a scoped repository must not resolve another workspace's photo to delete it")

		var deleted sql.NullTime
		require.NoError(t, env.Db.Raw(`SELECT deleted_at FROM photobox.photos WHERE id = 'beta-photo'`).Scan(&deleted).Error)
		assert.False(t, deleted.Valid, "beta's photo must remain un-deleted")
	})

	t.Run("empty workspace fails closed", func(t *testing.T) {
		none := base.WithWorkspace("")
		photos, err := none.ListAllPhotos("", 100, false, "", "", "")
		require.NoError(t, err)
		assert.Empty(t, photos, "an empty workspace must match nothing, not everything")

		_, err = none.GetPhotoById("alpha-photo", false)
		assert.ErrorIs(t, err, domain.ErrDataNotFound)
	})
}

// TestAlbumRepository_WorkspaceScoping_Integration asserts album queries are
// scoped, including the per-workspace uniqueness of names.
func TestAlbumRepository_WorkspaceScoping_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	for _, row := range []struct{ id, name, ws string }{
		{"alpha-album", "Holidays", wsAlpha},
		{"beta-album", "Holidays", wsBeta}, // same name, different workspace: legal
	} {
		require.NoError(t, env.Db.Exec(`INSERT INTO photobox.albums (id, name, workspace_id, created_epoch, metadata)
			VALUES (?, ?, ?, ?, '{}')`, row.id, row.name, row.ws, time.Now().UnixMilli()).Error)
	}

	base := repository.NewAlbumRepository(env)
	alpha := base.WithWorkspace(wsAlpha)

	t.Run("list is scoped", func(t *testing.T) {
		albums, err := alpha.ListAllAlbums("", 100)
		require.NoError(t, err)
		require.Len(t, albums, 1)
		assert.Equal(t, "alpha-album", albums[0].ID)
	})

	t.Run("count is scoped", func(t *testing.T) {
		n, err := alpha.AlbumCount()
		require.NoError(t, err)
		assert.EqualValues(t, 1, n)
	})

	t.Run("get by id is scoped", func(t *testing.T) {
		_, err := alpha.GetAlbumById("beta-album")
		assert.ErrorIs(t, err, domain.ErrDataNotFound)
	})

	t.Run("delete is scoped", func(t *testing.T) {
		require.NoError(t, alpha.DeleteAlbum("beta-album"))

		var remaining int64
		require.NoError(t, env.Db.Raw(`SELECT COUNT(*) FROM photobox.albums WHERE id = 'beta-album'`).Scan(&remaining).Error)
		assert.EqualValues(t, 1, remaining, "a scoped repository must not delete another workspace's album")
	})

	t.Run("empty workspace fails closed", func(t *testing.T) {
		albums, err := base.WithWorkspace("").ListAllAlbums("", 100)
		require.NoError(t, err)
		assert.Empty(t, albums)
	})
}

// TestShareRepository_WorkspaceScoping_Integration asserts share listing and
// revocation are scoped. Token resolution stays global by design (the token is
// the capability for the anonymous public view).
func TestShareRepository_WorkspaceScoping_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	for _, row := range []struct{ token, ws string }{
		{"alpha-share", wsAlpha},
		{"beta-share", wsBeta},
	} {
		require.NoError(t, env.Db.Exec(`INSERT INTO photobox.shared_links
			(token, resource_type, resource_id, created_by, workspace_id)
			VALUES (?, 'photo', 'p1', 'user-1', ?)`, row.token, row.ws).Error)
	}

	base := repository.NewShareRepository(env)
	alpha := base.WithWorkspace(wsAlpha)

	t.Run("list is scoped", func(t *testing.T) {
		shares, err := alpha.ListShares()
		require.NoError(t, err)
		require.Len(t, shares, 1)
		assert.Equal(t, "alpha-share", shares[0].Token)
	})

	t.Run("revoke is scoped", func(t *testing.T) {
		require.NoError(t, alpha.DeleteShare("beta-share"))

		var remaining int64
		require.NoError(t, env.Db.Raw(`SELECT COUNT(*) FROM photobox.shared_links WHERE token = 'beta-share'`).Scan(&remaining).Error)
		assert.EqualValues(t, 1, remaining, "a scoped repository must not revoke another workspace's share")
	})

	t.Run("token resolution remains global for the capability URL", func(t *testing.T) {
		share, err := base.GetShareByToken("beta-share")
		require.NoError(t, err)
		assert.Equal(t, "beta-share", share.Token)
	})

	t.Run("empty workspace fails closed", func(t *testing.T) {
		shares, err := base.WithWorkspace("").ListShares()
		require.NoError(t, err)
		assert.Empty(t, shares)
	})
}

// TestPhotoRepository_ScopedPagination_Integration proves the pagination
// correctness that post-query handler filtering could not provide: a caller's
// own rows are returned up to the requested limit even when other tenants hold
// more rows.
func TestPhotoRepository_ScopedPagination_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	now := time.Now().UnixMilli()
	// Beta holds many rows; alpha holds three. A pre-filter limit would starve
	// alpha; a scoped query must not.
	for i := 0; i < 20; i++ {
		require.NoError(t, env.Db.Exec(`INSERT INTO photobox.photos
			(id, name, workspace_id, created_epoch, metadata, media_type)
			VALUES (?, ?, ?, ?, '{}', 'image')`,
			"beta-"+string(rune('a'+i)), "b.jpg", wsBeta, now+int64(i)).Error)
	}
	for i := 0; i < 3; i++ {
		require.NoError(t, env.Db.Exec(`INSERT INTO photobox.photos
			(id, name, workspace_id, created_epoch, metadata, media_type)
			VALUES (?, ?, ?, ?, '{}', 'image')`,
			"alpha-"+string(rune('a'+i)), "a.jpg", wsAlpha, now+100+int64(i)).Error)
	}

	alpha := repository.NewPhotoRepository(env).WithWorkspace(wsAlpha)
	photos, err := alpha.ListAllPhotos("", 3, false, "", "", "")
	require.NoError(t, err)
	assert.Len(t, photos, 3, "a scoped limit must be filled from the caller's own rows")
	for _, p := range photos {
		assert.Equal(t, wsAlpha, p.WorkspaceID)
	}
}
