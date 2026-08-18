package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
)

// TestNewFilesystemRepository tests repository creation
func TestNewFilesystemRepository(t *testing.T) {
	timezone, _ := time.LoadLocation("UTC")
	config := appconfig.AppConfig{
		PhotoDir: "/test/photos",
		Timezone: timezone,
	}

	var jobService *service.JobService = nil

	repo := NewFilesystemRepository(config, jobService)

	assert.NotNil(t, repo)
	assert.Equal(t, config, repo.config)
	assert.Equal(t, jobService, repo.jobSvc)
}

// TestFilesystemRepository_StructureValidation tests the repository struct is properly defined
func TestFilesystemRepository_StructureValidation(t *testing.T) {
	timezone, _ := time.LoadLocation("UTC")
	config := appconfig.AppConfig{
		PhotoDir: "/test/photos",
		Timezone: timezone,
	}

	var jobService *service.JobService = nil

	repo := NewFilesystemRepository(config, jobService)

	// Verify the repository has the expected fields
	assert.NotNil(t, repo)
	assert.Equal(t, "/test/photos", repo.config.PhotoDir)
}

// NOTE: The following tests are documented as requiring integration testing
// due to their dependency on actual file system operations and image processing:
//
// 1. TestCreateDirectoryIfNotExists - Requires actual filesystem operations
//    - Needs temporary directory creation
//    - Tests os.Mkdir and error handling
//    - Should be tested in integration test suite
//
// 2. TestScanFilesystem - Requires actual filesystem scanning
//    - Needs temporary directory structure with test image files
//    - Tests concurrent goroutines and channels
//    - Depends on filepath.WalkDir and sync.WaitGroup
//    - Should be tested in integration test suite
//
// 3. TestWalkDir - Requires actual filesystem walking
//    - Needs real directory structure
//    - Tests recursive directory traversal with goroutines
//    - Depends on filepath.WalkDir and utils.IsImageFile
//    - Should be tested in integration test suite
//
// 4. TestGenerateThumbnail - Requires actual image processing
//    - Needs real image files
//    - Tests imaging library (disintegration/imaging)
//    - Depends on imaging.Open, imaging.Thumbnail, imaging.Encode
//    - Should be tested in integration test suite
//
// These methods involve complex file I/O, concurrent operations with goroutines,
// channels, and external library calls (image processing) that are better tested
// with actual files in an integration test environment.

// TestFFmpegPoolBounds verifies the ffmpeg semaphore caps concurrent process
// spawns to the configured pool size (issue #41).
func TestFFmpegPoolBounds(t *testing.T) {
	timezone, _ := time.LoadLocation("UTC")

	t.Run("semaphore capacity equals pool size", func(t *testing.T) {
		config := appconfig.AppConfig{PhotoDir: t.TempDir(), Timezone: timezone, FFmpegPoolSize: 3}
		repo := NewFilesystemRepository(config, nil)
		assert.Equal(t, 3, cap(repo.ffmpegSem), "semaphore should hold the pool size")
	})

	t.Run("default pool size is 2 when unset", func(t *testing.T) {
		config := appconfig.AppConfig{PhotoDir: t.TempDir(), Timezone: timezone}
		repo := NewFilesystemRepository(config, nil)
		assert.Equal(t, 2, cap(repo.ffmpegSem), "unset pool should default to 2")
	})

	t.Run("pool of 1 serialises video thumbnail calls", func(t *testing.T) {
		config := appconfig.AppConfig{PhotoDir: t.TempDir(), Timezone: timezone, FFmpegPoolSize: 1}
		repo := NewFilesystemRepository(config, nil)
		assert.Equal(t, 1, cap(repo.ffmpegSem))
		// Acquiring both slots: the second send must block until the first is
		// released, proving the semaphore enforces the cap.
		repo.ffmpegSem <- struct{}{}
		select {
		case repo.ffmpegSem <- struct{}{}:
			t.Fatal("second acquisition should block with pool size 1")
		default:
			// Expected: channel is full.
		}
		<-repo.ffmpegSem
	})
}
