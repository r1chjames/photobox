//go:build integration
// +build integration

package integration

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database/repository"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/test/testutil"
	"gorm.io/datatypes"
)

// TestAlbumRepository_CreateAndGet_Integration tests album CRUD operations
func TestAlbumRepository_CreateAndGet_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	repo := repository.NewAlbumRepository(env)

	// Create an album
	album, err := repo.CreateAlbum("Test Album")
	assert.NoError(t, err)
	assert.NotNil(t, album)
	assert.NotEmpty(t, album.ID)
	assert.Equal(t, "Test Album", album.Name)

	// Retrieve the album
	retrieved, err := repo.GetAlbumById(album.ID)
	assert.NoError(t, err)
	assert.Equal(t, album.ID, retrieved.ID)
	assert.Equal(t, "Test Album", retrieved.Name)
}

// TestAlbumRepository_ListAllAlbums_Integration tests pagination
func TestAlbumRepository_ListAllAlbums_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	repo := repository.NewAlbumRepository(env)

	// Create multiple albums
	for i := 0; i < 15; i++ {
		_, err := repo.CreateAlbum(string(rune('A' + i)))
		assert.NoError(t, err)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	// Test pagination - first page
	albums, err := repo.ListAllAlbums("", 10)
	assert.NoError(t, err)
	assert.Len(t, albums, 10, "Should return 10 albums on first page")

	// Test pagination - second page. ListAllAlbums expects the cursor to be the
	// base64-encoded created_epoch of the last item of the previous page.
	if len(albums) > 0 {
		cursor := base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(albums[len(albums)-1].CreatedEpoch, 10)))
		albums2, err := repo.ListAllAlbums(cursor, 10)
		assert.NoError(t, err)
		assert.Len(t, albums2, 5, "Should return remaining 5 albums")
	}
}

// TestPhotoRepository_CreateAndList_Integration tests photo operations
func TestPhotoRepository_CreateAndList_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	// Seed test data
	testutil.SeedTestData(t, env.Db)

	photoRepo := repository.NewPhotoRepository(env)

	// Create test photos
	photos := []domain.Photo{
		{
			ID:             "photo1",
			AlbumId:        "album1",
			FileHash:       "hash1",
			FilesystemPath: "/photos/album1/photo1.jpg",
			Name:           "photo1.jpg",
			MediaType:      "image/jpeg",
			Metadata:       datatypes.JSON([]byte("{}")),
		},
		{
			ID:             "photo2",
			AlbumId:        "album1",
			FileHash:       "hash2",
			FilesystemPath: "/photos/album1/photo2.jpg",
			Name:           "photo2.jpg",
			MediaType:      "image/jpeg",
			Metadata:       datatypes.JSON([]byte("{}")),
		},
		{
			ID:             "photo3",
			AlbumId:        "album2",
			FileHash:       "hash3",
			FilesystemPath: "/photos/album2/photo3.jpg",
			Name:           "photo3.jpg",
			MediaType:      "image/jpeg",
			Metadata:       datatypes.JSON([]byte("{}")),
		},
	}

	for _, photo := range photos {
		err := photoRepo.CreatePhotoInfo(photo)
		assert.NoError(t, err)
	}

	// List all photos
	allPhotos, err := photoRepo.ListAllPhotos("", 10, false, "", "", "")
	assert.NoError(t, err)
	assert.Len(t, allPhotos, 3, "Should return all 3 photos")

	// List photos in specific album
	albumPhotos, err := photoRepo.ListAllPhotosInAlbum("album1", "", 10, false, "", "", "")
	assert.NoError(t, err)
	assert.Len(t, albumPhotos, 2, "Should return 2 photos from album1")

	// Verify all returned photos belong to album1
	for _, photo := range albumPhotos {
		assert.Equal(t, "album1", photo.AlbumId)
	}
}

// TestPhotoRepository_GetPhotosInAlbumCount_Integration tests photo counting
func TestPhotoRepository_GetPhotosInAlbumCount_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	testutil.SeedTestData(t, env.Db)

	photoRepo := repository.NewPhotoRepository(env)

	// Create photos in album1
	for i := 0; i < 5; i++ {
		photo := domain.Photo{
			ID:             string(rune('A' + i)),
			AlbumId:        "album1",
			FileHash:       string(rune('h' + i)),
			FilesystemPath: "/test/path",
			Name:           "test.jpg",
			MediaType:      "image/jpeg",
			Metadata:       datatypes.JSON([]byte("{}")),
		}
		err := photoRepo.CreatePhotoInfo(photo)
		assert.NoError(t, err)
	}

	// Count photos in album1
	count, err := photoRepo.GetPhotosInAlbumCount("album1")
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count, "Should count 5 photos in album1")

	// Count photos in empty album
	count2, err := photoRepo.GetPhotosInAlbumCount("album2")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count2, "Should count 0 photos in album2")
}

// TestPhotoRepository_TrashRetention_Integration tests the age-based trash
// purge used by the retention job.
func TestPhotoRepository_TrashRetention_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	testutil.SeedTestData(t, env.Db)

	photoRepo := repository.NewPhotoRepository(env)

	// Create three photos
	for i := 1; i <= 3; i++ {
		photo := domain.Photo{
			ID:             fmt.Sprintf("trash%d", i),
			AlbumId:        "album1",
			FileHash:       fmt.Sprintf("hash%d", i),
			FilesystemPath: "/test/path",
			Name:           fmt.Sprintf("trash%d.jpg", i),
			MediaType:      "image/jpeg",
			Metadata:       datatypes.JSON([]byte("{}")),
		}
		err := photoRepo.CreatePhotoInfo(photo)
		assert.NoError(t, err)
	}

	// Soft-delete all three
	for i := 1; i <= 3; i++ {
		_, err := photoRepo.SoftDeletePhoto(fmt.Sprintf("trash%d", i))
		assert.NoError(t, err)
	}

	// Backdate two of them beyond the retention window (60 days)
	oldTime := time.Now().AddDate(0, 0, -90)
	err := env.Db.Model(&domain.Photo{}).Where("id IN ?", []string{"trash1", "trash2"}).Update("deleted_at", oldTime).Error
	assert.NoError(t, err)

	// ListExpiredTrashPhotos should return only the two old ones
	cutoff := time.Now().AddDate(0, 0, -60)
	expired, err := photoRepo.ListExpiredTrashPhotos(cutoff)
	assert.NoError(t, err)
	assert.Len(t, expired, 2, "Should return the 2 photos deleted before the cutoff")

	// PermanentlyDeletePhotos removes only those rows
	err = photoRepo.PermanentlyDeletePhotos([]string{expired[0].ID, expired[1].ID})
	assert.NoError(t, err)

	remaining, err := photoRepo.ListTrashPhotos("", 100, false)
	assert.NoError(t, err)
	assert.Len(t, remaining, 1, "One photo should remain in trash")
}

// TestUserRepository_CreateAndAuth_Integration tests user authentication flow
func TestUserRepository_CreateAndAuth_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	userRepo := repository.NewUserRepository(env)

	// Create a user
	user := &domain.User{
		ID:       "user1",
		Username: "testuser",
		Password: "$2a$10$hashedpassword",
		Role:     domain.VIEWER,
	}

	created, err := userRepo.CreateUser(user)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)

	// Get user by username
	retrieved, err := userRepo.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", retrieved.Username)
	assert.Equal(t, domain.VIEWER, retrieved.Role)

	// Get user by ID
	byID, err := userRepo.GetUserById(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, byID.ID)
}

// TestUserRepository_ListUsers_Integration tests user listing with pagination
func TestUserRepository_ListUsers_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	userRepo := repository.NewUserRepository(env)

	// Create multiple users. Each user needs a unique ID and email — CreateUser
	// uses OnConflict DoNothing, and users with an empty email all collide on
	// the unique email index and would silently be dropped.
	for i := 0; i < 12; i++ {
		user := &domain.User{
			ID:       "listuser-" + strconv.Itoa(i),
			Username: string(rune('A'+i)) + "user",
			Email:    "user" + strconv.Itoa(i) + "@example.com",
			Password: "hashedpass",
			Role:     domain.VIEWER,
		}
		_, err := userRepo.CreateUser(user)
		assert.NoError(t, err)
		time.Sleep(5 * time.Millisecond)
	}

	// List first page
	users, err := userRepo.ListUsers(1, 10)
	assert.NoError(t, err)
	assert.Len(t, users, 10, "Should return 10 users on first page")

	// List second page
	users2, err := userRepo.ListUsers(2, 10)
	assert.NoError(t, err)
	assert.Len(t, users2, 2, "Should return remaining 2 users")
}

// TestJobRepository_StatusUpdates_Integration tests job status management
func TestJobRepository_StatusUpdates_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	testutil.SeedTestData(t, env.Db)

	jobRepo := repository.NewJobRepository(env)

	// Check job is not running
	isRunning, err := jobRepo.IsJobRunning("photo_index")
	assert.NoError(t, err)
	assert.False(t, isRunning)

	// Start job
	err = jobRepo.UpdateJobStatus("photo_index", "RUNNING")
	assert.NoError(t, err)

	// Verify job is running
	isRunning, err = jobRepo.IsJobRunning("photo_index")
	assert.NoError(t, err)
	assert.True(t, isRunning)

	// Stop job
	err = jobRepo.UpdateJobStatus("photo_index", "NOT_RUNNING")
	assert.NoError(t, err)

	// Verify job is stopped
	isRunning, err = jobRepo.IsJobRunning("photo_index")
	assert.NoError(t, err)
	assert.False(t, isRunning)
}

// TestSettingsRepository_CRUD_Integration tests settings operations
func TestSettingsRepository_CRUD_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	settingsRepo := repository.NewUtilityRepository(env)

	// Get non-existent setting — the repo returns an empty setting, no error
	setting, err := settingsRepo.GetSetting("nonexistent")
	assert.NoError(t, err)
	assert.Empty(t, setting.Value, "Non-existent setting should have an empty value")

	// Create base settings
	err = settingsRepo.CreateBaseSettings(false)
	assert.NoError(t, err)

	// Get all settings
	settings, err := settingsRepo.GetAllSettings()
	assert.NoError(t, err)
	assert.NotEmpty(t, settings, "Should have created base settings")

	// Update a setting
	if len(settings) > 0 {
		setting := settings[0]
		setting.Value = "updated_value"
		err = settingsRepo.UpdateSetting(setting)
		assert.NoError(t, err)

		// Verify update
		retrieved, err := settingsRepo.GetSetting(setting.Key)
		assert.NoError(t, err)
		assert.Equal(t, "updated_value", retrieved.Value)
	}
}

// TestPhotoRepository_AdvancedSearch_Integration tests combined filters
func TestPhotoRepository_AdvancedSearch_Integration(t *testing.T) {
	env := testutil.CreateTestEnv(t)
	defer testutil.CleanupTestEnv(t, env)

	testutil.SeedTestData(t, env.Db)

	photoRepo := repository.NewPhotoRepository(env)

	// Seed photos with different cameras, GPS presence, orientation
	now := time.Now()
	photos := []domain.Photo{
		{ID: "c1", Name: "canon-landscape-gps.jpg", AlbumId: "album1", FilesystemPath: "/photos/a.jpg",
			Metadata: datatypes.JSON([]byte(`{"exif":{"Make":"Canon","Model":"EOS R5"}}`)),
			Latitude: 52.0, Longitude: 0.5, Width: 2000, Height: 1000,
			CreatedEpoch: now.UnixMilli(), Year: now.Year(), Month: int(now.Month())},
		{ID: "c2", Name: "canon-portrait-nogps.jpg", AlbumId: "album1", FilesystemPath: "/photos/b.jpg",
			Metadata: datatypes.JSON([]byte(`{"exif":{"Make":"Canon","Model":"EOS R6"}}`)),
			Width: 1000, Height: 2000,
			CreatedEpoch: now.UnixMilli(), Year: now.Year(), Month: int(now.Month())},
		{ID: "i1", Name: "iphone-square-gps.jpg", AlbumId: "album1", FilesystemPath: "/photos/c.jpg",
			Metadata: datatypes.JSON([]byte(`{"exif":{"Make":"Apple","Model":"iPhone 15"}}`)),
			Latitude: 51.0, Longitude: -1.0, Width: 1000, Height: 1000,
			CreatedEpoch: now.UnixMilli(), Year: now.Year(), Month: int(now.Month())},
	}
	for _, p := range photos {
		assert.NoError(t, photoRepo.CreatePhotoInfo(p))
	}

	t.Run("camera filter", func(t *testing.T) {
		res, err := photoRepo.SearchPhotosWithFilters(domain.PhotoSearchFilters{Camera: "Canon"}, "", 100, false)
		assert.NoError(t, err)
		assert.Len(t, res, 2, "both Canon photos")
	})

	t.Run("camera + hasGps", func(t *testing.T) {
		res, err := photoRepo.SearchPhotosWithFilters(domain.PhotoSearchFilters{Camera: "Canon", HasGPS: true}, "", 100, false)
		assert.NoError(t, err)
		assert.Len(t, res, 1, "only the Canon photo with GPS")
		assert.Equal(t, "c1", res[0].ID)
	})

	t.Run("orientation portrait", func(t *testing.T) {
		res, err := photoRepo.SearchPhotosWithFilters(domain.PhotoSearchFilters{Orientation: "portrait"}, "", 100, false)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "c2", res[0].ID)
	})

	t.Run("hasGps only", func(t *testing.T) {
		res, err := photoRepo.SearchPhotosWithFilters(domain.PhotoSearchFilters{HasGPS: true}, "", 100, false)
		assert.NoError(t, err)
		assert.Len(t, res, 2, "c1 and i1 have GPS")
	})

	t.Run("no filters returns all", func(t *testing.T) {
		res, err := photoRepo.SearchPhotosWithFilters(domain.PhotoSearchFilters{}, "", 100, false)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(res), 3)
	})
}
