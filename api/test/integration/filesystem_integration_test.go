//go:build integration
// +build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/filesystem/repository"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/test/testutil"
)

// TestFilesystemRepository_ScanFilesystem_Integration tests real filesystem scanning
func TestFilesystemRepository_ScanFilesystem_Integration(t *testing.T) {
	// Create test directory with photos
	testDir := testutil.CreateTestPhotoDir(t)
	defer testutil.CleanupTestPhotoDir(t, testDir)

	// Create test photo structure
	testutil.CreateTestPhotoStructure(t, testDir)

	// Verify test photos were created
	fileCount := testutil.CountFilesInDir(t, testDir)
	assert.Equal(t, 9, fileCount, "Should have created 9 test images")

	// Setup filesystem repository
	timezone, _ := time.LoadLocation("UTC")
	config := appconfig.AppConfig{
		PhotoDir: testDir,
		Timezone: timezone,
	}

	repo := repository.NewFilesystemRepository(config, nil)

	// Scan filesystem
	photoChan := make(chan string, 100)
	go func() {
		repo.ScanFilesystem(photoChan)
		close(photoChan)
	}()

	// Collect scanned photos
	var scannedPhotos []string
	for photo := range photoChan {
		scannedPhotos = append(scannedPhotos, photo)
	}

	// Verify all photos were scanned
	assert.Len(t, scannedPhotos, 9, "Should have scanned all 9 photos")

	// Verify scanned paths are absolute
	for _, photo := range scannedPhotos {
		assert.True(t, filepath.IsAbs(photo), "Photo path should be absolute: %s", photo)
	}
}

// TestFilesystemRepository_GenerateThumbnail_Integration tests real thumbnail generation
func TestFilesystemRepository_GenerateThumbnail_Integration(t *testing.T) {
	// Create test directory
	testDir := testutil.CreateTestPhotoDir(t)
	defer testutil.CleanupTestPhotoDir(t, testDir)

	// Create a test image
	testImagePath := filepath.Join(testDir, "test.jpg")
	testutil.CreateTestImage(t, testImagePath, 1920, 1080, "jpg")

	// Setup repository
	timezone, _ := time.LoadLocation("UTC")
	config := appconfig.AppConfig{
		PhotoDir: testDir,
		Timezone: timezone,
	}

	repo := repository.NewFilesystemRepository(config, nil)

	// Generate thumbnail
	thumbnail := repo.GenerateThumbnail(testImagePath, 200, 200)

	// Verify thumbnail was generated
	assert.NotNil(t, thumbnail, "Thumbnail should be generated")
	assert.Greater(t, len(thumbnail), 0, "Thumbnail should have data")
	assert.Less(t, len(thumbnail), 100000, "Thumbnail should be smaller than 100KB")
}

// TestFilesystemRepository_GenerateThumbnail_PNG_Integration tests PNG thumbnail generation
func TestFilesystemRepository_GenerateThumbnail_PNG_Integration(t *testing.T) {
	testDir := testutil.CreateTestPhotoDir(t)
	defer testutil.CleanupTestPhotoDir(t, testDir)

	// Create a PNG test image
	testImagePath := filepath.Join(testDir, "test.png")
	testutil.CreateTestImage(t, testImagePath, 1024, 768, "png")

	timezone, _ := time.LoadLocation("UTC")
	config := appconfig.AppConfig{
		PhotoDir: testDir,
		Timezone: timezone,
	}

	repo := repository.NewFilesystemRepository(config, nil)

	// Generate thumbnail
	thumbnail := repo.GenerateThumbnail(testImagePath, 200, 200)

	assert.NotNil(t, thumbnail)
	assert.Greater(t, len(thumbnail), 0)
}

// TestFilesystemRepository_GenerateThumbnail_InvalidFile_Integration tests error handling
func TestFilesystemRepository_GenerateThumbnail_InvalidFile_Integration(t *testing.T) {
	testDir := testutil.CreateTestPhotoDir(t)
	defer testutil.CleanupTestPhotoDir(t, testDir)

	timezone, _ := time.LoadLocation("UTC")
	config := appconfig.AppConfig{
		PhotoDir: testDir,
		Timezone: timezone,
	}

	repo := repository.NewFilesystemRepository(config, nil)

	// Try to generate thumbnail for non-existent file
	thumbnail := repo.GenerateThumbnail("/nonexistent/file.jpg", 200, 200)

	// Should return nil for invalid file
	assert.Nil(t, thumbnail, "Should return nil for non-existent file")
}

// TestFilesystemRepository_CreateDirectoryIfNotExists_Integration tests directory creation
func TestFilesystemRepository_CreateDirectoryIfNotExists_Integration(t *testing.T) {
	testDir := testutil.CreateTestPhotoDir(t)
	defer testutil.CleanupTestPhotoDir(t, testDir)

	timezone, _ := time.LoadLocation("UTC")
	config := appconfig.AppConfig{
		PhotoDir: testDir,
		Timezone: timezone,
	}

	repo := repository.NewFilesystemRepository(config, nil)

	// Create a new album directory
	albumName := "NewAlbum"
	repo.CreateDirectoryIfNotExists(testDir, albumName)

	// Verify directory was created
	albumPath := filepath.Join(testDir, albumName)
	info, err := filepath.Glob(albumPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, info, "Album directory should exist")
}

// TestFilesystemScanConcurrency_Integration tests concurrent filesystem operations
func TestFilesystemScanConcurrency_Integration(t *testing.T) {
	testDir := testutil.CreateTestPhotoDir(t)
	defer testutil.CleanupTestPhotoDir(t, testDir)

	// Create a larger directory structure
	testutil.CreateTestPhotoStructure(t, testDir)

	timezone, _ := time.LoadLocation("UTC")
	config := appconfig.AppConfig{
		PhotoDir: testDir,
		Timezone: timezone,
	}

	// Run scan multiple times concurrently to test for race conditions.
	// Each goroutine uses its own repository instance because ScanFilesystem
	// owns mutable state (WaitGroup + dir semaphore) that is not safe to share.
	done := make(chan bool)

	for i := 0; i < 3; i++ {
		go func() {
			r := repository.NewFilesystemRepository(config, nil)
			photoChan := make(chan string, 100)
			r.ScanFilesystem(photoChan)
			close(photoChan)

			count := 0
			for range photoChan {
				count++
			}
			assert.Equal(t, 9, count)
			done <- true
		}()
	}

	// Wait for all scans to complete
	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestFilesystemRepository_MoveAndRestoreFromTrash_Integration verifies the
// trash round-trip: a file moved to .trash can be restored to its original
// path, and both locations are consistent.
func TestFilesystemRepository_MoveAndRestoreFromTrash_Integration(t *testing.T) {
	testDir := testutil.CreateTestPhotoDir(t)
	defer testutil.CleanupTestPhotoDir(t, testDir)

	testutil.CreateTestPhotoStructure(t, testDir)

	timezone, _ := time.LoadLocation("UTC")
	config := appconfig.AppConfig{
		PhotoDir: testDir,
		Timezone: timezone,
	}
	repo := repository.NewFilesystemRepository(config, nil)

	// Pick the first photo and note its original path
	photoChan := make(chan string, 10)
	go func() {
		repo.ScanFilesystem(photoChan)
		close(photoChan)
	}()
	originalPath := ""
	for path := range photoChan {
		originalPath = path
		break
	}
	assert.NotEmpty(t, originalPath, "should find at least one photo")
	_, err := os.Stat(originalPath)
	assert.NoError(t, err, "original should exist before move")

	// Move to trash
	trashPath, err := repo.MoveToTrash(originalPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, trashPath)
	assert.NotEqual(t, originalPath, trashPath)

	// Original gone, trash present
	_, err = os.Stat(originalPath)
	assert.Error(t, err, "original should be gone after move")
	_, err = os.Stat(trashPath)
	assert.NoError(t, err, "trash file should exist")

	// Restore
	err = repo.RestoreFromTrash(trashPath, originalPath)
	assert.NoError(t, err)

	// Original back, trash gone
	_, err = os.Stat(originalPath)
	assert.NoError(t, err, "original should be restored")
	_, err = os.Stat(trashPath)
	assert.Error(t, err, "trash file should be gone after restore")
}
