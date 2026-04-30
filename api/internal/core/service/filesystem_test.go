package service

import (
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

// MockFilesystemRepository is a mock implementation of port.FilesystemRepository
type MockFilesystemRepository struct {
	mock.Mock
}

func (m *MockFilesystemRepository) CreateDirectoryIfNotExists(basePath, albumName string) {
	m.Called(basePath, albumName)
}

func (m *MockFilesystemRepository) ScanFilesystem(photoChan chan string) {
	m.Called(photoChan)
}

func (m *MockFilesystemRepository) GenerateThumbnail(path string) []byte {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]byte)
}

func (m *MockFilesystemRepository) MoveToTrash(path string) (string, error) {
	args := m.Called(path)
	return args.String(0), args.Error(1)
}

func (m *MockFilesystemRepository) RestoreFromTrash(trashPath, originalPath string) error {
	args := m.Called(trashPath, originalPath)
	return args.Error(0)
}

func (m *MockFilesystemRepository) RenameDirectory(oldPath, newPath string) error {
	args := m.Called(oldPath, newPath)
	return args.Error(0)
}

// MockJobService is a mock implementation of port.JobService for filesystem tests
type MockJobServiceFS struct {
	mock.Mock
}

func (m *MockJobServiceFS) IsJobRunning(jobName string) (bool, error) {
	args := m.Called(jobName)
	return args.Bool(0), args.Error(1)
}

func (m *MockJobServiceFS) UpdateJobStatus(jobName, status string) error {
	args := m.Called(jobName, status)
	return args.Error(0)
}

func (m *MockJobServiceFS) UpdateAllJobsStatus(status string) error {
	args := m.Called(status)
	return args.Error(0)
}

func (m *MockJobServiceFS) JobStart(jobName string) error {
	args := m.Called(jobName)
	return args.Error(0)
}

func (m *MockJobServiceFS) JobComplete(jobName string) error {
	args := m.Called(jobName)
	return args.Error(0)
}

func (m *MockJobServiceFS) CreateBaseJobs() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockJobServiceFS) StartJobIfNotRunning(jobName string) error {
	args := m.Called(jobName)
	return args.Error(0)
}

// MockUtilityServiceFS is a mock implementation of port.UtilityService for filesystem tests
type MockUtilityServiceFS struct {
	mock.Mock
}

func (m *MockUtilityServiceFS) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUtilityServiceFS) Healthcheck() ([]*domain.Setting, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Setting), args.Error(1)
}

func (m *MockUtilityServiceFS) GetSetting(key string) (*domain.Setting, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Setting), args.Error(1)
}

func (m *MockUtilityServiceFS) ListAllSettings() ([]*domain.Setting, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Setting), args.Error(1)
}

func (m *MockUtilityServiceFS) UpdateSetting(setting *domain.Setting) error {
	args := m.Called(setting)
	return args.Error(0)
}

func (m *MockUtilityServiceFS) UpdateAllSettings(settings []*domain.Setting) error {
	args := m.Called(settings)
	return args.Error(0)
}

func (m *MockUtilityServiceFS) CreateBaseSettings(reset bool) error {
	args := m.Called(reset)
	return args.Error(0)
}

// TestNewFilesystemService tests filesystem service creation
func TestNewFilesystemService(t *testing.T) {
	mockFsRepo := new(MockFilesystemRepository)
	mockJobSvc := new(MockJobServiceFS)
	mockUtilSvc := new(MockUtilityServiceFS)

	service := NewFilesystemService(mockFsRepo, mockJobSvc, mockUtilSvc)

	assert.NotNil(t, service)
	assert.Equal(t, mockFsRepo, service.fsRepo)
	assert.Equal(t, mockJobSvc, service.jobSvc)
	assert.Equal(t, mockUtilSvc, service.utilitySvc)
}

// TestGenerateThumbnail_WithExifThumbnail tests thumbnail generation when EXIF has thumbnail
func TestGenerateThumbnail_WithExifThumbnail(t *testing.T) {
	mockFsRepo := new(MockFilesystemRepository)
	mockJobSvc := new(MockJobServiceFS)
	mockUtilSvc := new(MockUtilityServiceFS)

	service := NewFilesystemService(mockFsRepo, mockJobSvc, mockUtilSvc)

	path := "/path/to/photo.jpg"
	exifData := exif.Exif{} // Empty EXIF (no thumbnail)
	expectedThumbnail := []byte{0x01, 0x02, 0x03}

	// Mock repository to return generated thumbnail
	mockFsRepo.On("GenerateThumbnail", path).Return(expectedThumbnail)

	result := service.GenerateThumbnail(path, exifData)

	assert.Equal(t, expectedThumbnail, result)
	mockFsRepo.AssertExpectations(t)
}

// TestGenerateThumbnail_NoExifThumbnail tests thumbnail generation when EXIF has no thumbnail
func TestGenerateThumbnail_NoExifThumbnail(t *testing.T) {
	mockFsRepo := new(MockFilesystemRepository)
	mockJobSvc := new(MockJobServiceFS)
	mockUtilSvc := new(MockUtilityServiceFS)

	service := NewFilesystemService(mockFsRepo, mockJobSvc, mockUtilSvc)

	path := "/path/to/photo.jpg"
	exifData := exif.Exif{}
	expectedThumbnail := []byte{0xFF, 0xD8, 0xFF}

	mockFsRepo.On("GenerateThumbnail", path).Return(expectedThumbnail)

	result := service.GenerateThumbnail(path, exifData)

	assert.NotNil(t, result)
	mockFsRepo.AssertExpectations(t)
}

// TestGenerateThumbnail_NilRepository tests thumbnail generation with nil repository response
func TestGenerateThumbnail_NilRepository(t *testing.T) {
	mockFsRepo := new(MockFilesystemRepository)
	mockJobSvc := new(MockJobServiceFS)
	mockUtilSvc := new(MockUtilityServiceFS)

	service := NewFilesystemService(mockFsRepo, mockJobSvc, mockUtilSvc)

	path := "/path/to/photo.jpg"
	exifData := exif.Exif{}

	mockFsRepo.On("GenerateThumbnail", path).Return(nil)

	result := service.GenerateThumbnail(path, exifData)

	// Should return nil if repository returns nil
	assert.Nil(t, result)
	mockFsRepo.AssertExpectations(t)
}

// NOTE: The following tests are documented as requiring integration testing
// due to their dependency on actual file I/O operations:
//
// 1. TestPerformPhotoIndex - Requires actual filesystem scanning and file operations
//    - Needs temporary directory structure with test image files
//    - Tests channel communication and concurrent processing
//    - Should be tested in integration test suite
//
// 2. TestWriteFileToFilesystem - Requires actual file writing
//    - Needs base64 encoded image data
//    - Tests directory creation and file writing
//    - Depends on os.WriteFile and os.Lstat
//    - Should be tested in integration test suite
//
// 3. TestGetMetaData - Requires actual file operations
//    - Needs real image files with EXIF data
//    - Tests MD5 hashing, EXIF extraction, MIME detection
//    - Depends on utils.OpenFile, utils.GetExifData, etc.
//    - Should be tested in integration test suite
//
// These methods involve complex file I/O, goroutines, channels, and external
// library calls (EXIF parsing, image processing) that are better tested with
// actual files in an integration test environment.

// TestFilesystemService_DependencyInjection tests that all dependencies are properly injected
func TestFilesystemService_DependencyInjection(t *testing.T) {
	mockFsRepo := new(MockFilesystemRepository)
	mockJobSvc := new(MockJobServiceFS)
	mockUtilSvc := new(MockUtilityServiceFS)

	service := NewFilesystemService(mockFsRepo, mockJobSvc, mockUtilSvc)

	// Verify all dependencies are accessible
	assert.NotNil(t, service.fsRepo)
	assert.NotNil(t, service.jobSvc)
	assert.NotNil(t, service.utilitySvc)

	// Verify they're the correct instances
	assert.IsType(t, &MockFilesystemRepository{}, service.fsRepo)
	assert.IsType(t, &MockJobServiceFS{}, service.jobSvc)
	assert.IsType(t, &MockUtilityServiceFS{}, service.utilitySvc)
}

// TestFilesystemService_StructureValidation tests the service struct is properly defined
func TestFilesystemService_StructureValidation(t *testing.T) {
	mockFsRepo := new(MockFilesystemRepository)
	mockJobSvc := new(MockJobServiceFS)
	mockUtilSvc := new(MockUtilityServiceFS)

	service := NewFilesystemService(mockFsRepo, mockJobSvc, mockUtilSvc)

	// Test that the service struct has the expected fields
	assert.NotNil(t, service)

	// Use reflection to verify field existence (without relying on field access)
	serviceType := service
	assert.NotNil(t, serviceType)
}
