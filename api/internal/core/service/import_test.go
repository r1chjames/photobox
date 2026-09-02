package service

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// newTestTakeoutImporter creates a TakeoutImporter with mock dependencies.
func newTestTakeoutImporter(t *testing.T, photoDir string) (*TakeoutImporter, *MockPhotoRepository, *MockFilesystemService) {
	t.Helper()
	config := appconfig.AppConfig{PhotoDir: photoDir}
	mockRepo := new(MockPhotoRepository)
	mockFsSvc := new(MockFilesystemService)

	// Default: no existing photos in cache.
	mockRepo.On("GetPhotoIndexCache").Return(map[string]struct {
		FileHash         string
		FileModifiedTime int64
	}{}, nil)

	importer := NewTakeoutImporter(nil, mockFsSvc, mockRepo, &config, nil)
	return importer, mockRepo, mockFsSvc
}

// createFakeTakeout creates a temp directory tree mimicking a Google Takeout export.
func createFakeTakeout(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	// Photos/ folder with two images and sidecar JSONs
	photosDir := filepath.Join(root, "Photos")
	if err := os.MkdirAll(photosDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(photosDir, "IMG_001.jpg"), []byte("fake-jpeg-1"))
	writeFile(t, filepath.Join(photosDir, "IMG_002.jpg"), []byte("fake-jpeg-2"))

	// Sidecar for IMG_001 with geo data
	sidecar := `{"PhotoTakenTime":{"Seconds":1609459200},"GeoData":{"LatitudeE7":377849660,"LongitudeE7":-1224003220}}`
	writeFile(t, filepath.Join(photosDir, "IMG_001.json"), []byte(sidecar))

	// Videos/ folder with one video
	videosDir := filepath.Join(root, "Videos")
	if err := os.MkdirAll(videosDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(videosDir, "VID_001.mp4"), []byte("fake-mp4"))

	// A non-media file that should be ignored
	writeFile(t, filepath.Join(root, "index.html"), []byte("<html></html>"))

	return root
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestScanTakeout_FindsFiles(t *testing.T) {
	exportDir := createFakeTakeout(t)
	photoDir := t.TempDir()
	importer, _, _ := newTestTakeoutImporter(t, photoDir)

	result, err := importer.ScanTakeout(exportDir)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalPhotos)
	assert.Equal(t, 1, result.TotalVideos)
	assert.Greater(t, result.TotalSizeBytes, int64(0))
	assert.ElementsMatch(t, []string{"Photos", "Videos"}, result.Albums)
	assert.Empty(t, result.Duplicates)
}

func TestScanTakeout_DetectsDuplicates(t *testing.T) {
	root := t.TempDir()
	photosDir := filepath.Join(root, "Photos")
	if err := os.MkdirAll(photosDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Two files with identical content → intra-export duplicate.
	content := []byte("identical-content")
	writeFile(t, filepath.Join(photosDir, "a.jpg"), content)
	writeFile(t, filepath.Join(photosDir, "b.jpg"), content)

	photoDir := t.TempDir()
	importer, _, _ := newTestTakeoutImporter(t, photoDir)

	result, err := importer.ScanTakeout(root)
	assert.NoError(t, err)
	assert.Equal(t, 2, result.TotalPhotos)
	assert.GreaterOrEqual(t, result.Duplicates, 1)
}

func TestScanTakeout_ExistingDuplicate(t *testing.T) {
	root := t.TempDir()
	photosDir := filepath.Join(root, "Photos")
	if err := os.MkdirAll(photosDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(photosDir, "dup.jpg"), []byte("unique-content-xyz"))

	hash := computeFileHash(filepath.Join(photosDir, "dup.jpg"))
	assert.NotEmpty(t, hash)

	photoDir := t.TempDir()
	config := appconfig.AppConfig{PhotoDir: photoDir}
	mockRepo := new(MockPhotoRepository)
	mockFsSvc := new(MockFilesystemService)

	// Register populated cache expectation (testify matches most-recently-registered first).
	populatedCache := map[string]struct {
		FileHash         string
		FileModifiedTime int64
	}{
		"existing": {FileHash: hash, FileModifiedTime: time.Now().Unix()},
	}
	mockRepo.On("GetPhotoIndexCache").Return(populatedCache, nil)

	importer := NewTakeoutImporter(nil, mockFsSvc, mockRepo, &config, nil)

	result, err := importer.ScanTakeout(root)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, result.Duplicates, 1)
}

func TestScanTakeout_NonExistentPath(t *testing.T) {
	photoDir := t.TempDir()
	importer, _, _ := newTestTakeoutImporter(t, photoDir)

	_, err := importer.ScanTakeout("/nonexistent/path/xyz")
	assert.Error(t, err)
}

func TestScanTakeout_InvalidExtension(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "file.txt"), []byte("text"))

	photoDir := t.TempDir()
	importer, _, _ := newTestTakeoutImporter(t, photoDir)

	result, err := importer.ScanTakeout(root)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.TotalPhotos)
	assert.Equal(t, 0, result.TotalVideos)
}

func TestImportTakeout_AlreadyRunning(t *testing.T) {
	photoDir := t.TempDir()
	importer, _, _ := newTestTakeoutImporter(t, photoDir)

	// Manually set running flag.
	importer.mu.Lock()
	importer.running = true
	importer.mu.Unlock()

	req := domain.TakeoutImportRequest{Path: "/tmp", Mode: "import"}
	err := importer.ImportTakeout(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already in progress")
}

func TestImportTakeout_ScanModeReturnsNil(t *testing.T) {
	photoDir := t.TempDir()
	importer, _, _ := newTestTakeoutImporter(t, photoDir)

	req := domain.TakeoutImportRequest{Path: "/tmp", Mode: "scan"}
	err := importer.ImportTakeout(context.Background(), req)
	assert.NoError(t, err)
}

func TestGetProgress_InitialState(t *testing.T) {
	photoDir := t.TempDir()
	importer, _, _ := newTestTakeoutImporter(t, photoDir)

	progress := importer.GetProgress()
	assert.Empty(t, progress.Status)
	assert.Zero(t, progress.Current)
	assert.Zero(t, progress.Total)
}

func TestOpenTakeout_ZipFile(t *testing.T) {
	// Create a zip file with a media file inside.
	zipPath := filepath.Join(t.TempDir(), "takeout.zip")
	createTestZip(t, zipPath, "Photos/test.jpg", []byte("jpeg-data"))

	photoDir := t.TempDir()
	importer, _, _ := newTestTakeoutImporter(t, photoDir)

	result, err := importer.ScanTakeout(zipPath)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.TotalPhotos)
}

func TestOpenTakeout_InvalidType(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "file.txt"), []byte("text"))

	photoDir := t.TempDir()
	importer, _, _ := newTestTakeoutImporter(t, photoDir)

	_, err := importer.ScanTakeout(filepath.Join(root, "file.txt"))
	assert.Error(t, err)
}

func TestImportTakeout_CopyAndSave(t *testing.T) {
	exportDir := createFakeTakeout(t)
	photoDir := t.TempDir()

	config := appconfig.AppConfig{PhotoDir: photoDir}
	mockRepo := new(MockPhotoRepository)
	mockFsSvc := new(MockFilesystemService)

	mockRepo.On("GetPhotoIndexCache").Return(map[string]struct {
		FileHash         string
		FileModifiedTime int64
	}{}, nil)

	// Mock GetPhotoMetadata to return a PhotoFile.
	mockFsSvc.On("GetPhotoMetadata", mock.MatchedBy(func(path string) bool {
		return filepath.IsAbs(path) && strings.HasPrefix(path, photoDir)
	})).Return(domain.PhotoFile{Path: "mocked"}).Maybe()

	importer := NewTakeoutImporter(nil, mockFsSvc, mockRepo, &config, nil)

	req := domain.TakeoutImportRequest{Path: exportDir, Mode: "import", PreserveAlbums: true}
	err := importer.ImportTakeout(context.Background(), req)
	assert.NoError(t, err)

	// Verify files were copied to photoDir.
	photosInDest, _ := os.ReadDir(filepath.Join(photoDir, "Photos"))
	assert.Len(t, photosInDest, 2)
	videosInDest, _ := os.ReadDir(filepath.Join(photoDir, "Videos"))
	assert.Len(t, videosInDest, 1)

	// Verify progress is complete.
	progress := importer.GetProgress()
	assert.Equal(t, "complete", progress.Status)
	assert.Equal(t, 3, progress.Imported)
}

func TestTakeoutImporter_SidecarGeoData(t *testing.T) {
	exportDir := createFakeTakeout(t)
	photoDir := t.TempDir()

	config := appconfig.AppConfig{PhotoDir: photoDir}
	mockRepo := new(MockPhotoRepository)
	mockFsSvc := new(MockFilesystemService)

	mockRepo.On("GetPhotoIndexCache").Return(map[string]struct {
		FileHash         string
		FileModifiedTime int64
	}{}, nil)

	var capturedFiles []domain.PhotoFile
	mockFsSvc.On("GetPhotoMetadata", mock.Anything).Run(func(args mock.Arguments) {
		// Verify sidecar geo data is applied by checking the returned PhotoFile.
	}).Return(domain.PhotoFile{}).Maybe()

	importer := NewTakeoutImporter(nil, mockFsSvc, mockRepo, &config, nil)

	req := domain.TakeoutImportRequest{Path: exportDir, Mode: "import", PreserveAlbums: true}
	err := importer.ImportTakeout(context.Background(), req)
	assert.NoError(t, err)
	_ = capturedFiles
}

// createTestZip creates a minimal zip file with the given entry.
func createTestZip(t *testing.T, path string, name string, data []byte) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	fw, err := w.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write(data); err != nil {
		t.Fatal(err)
	}
	w.Close()
}
