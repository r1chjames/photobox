package port

import (
	"github.com/rwcarlsen/goexif/exif"
	. "gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

//go:generate mockgen -source=filesystem.go -destination=mock/filesystem.go -package=mock

// FilesystemRepository is an interface for interacting with filesystem-related business logic
type FilesystemRepository interface {
	// ScanFilesystem scans a directory on the filesystem
	ScanFilesystem(photoChan chan string)
	GenerateThumbnail(path string) []byte
	CreateDirectoryIfNotExists(basePhotoPath string, directoryName string)
}

// FilesystemService is an interface for interacting with filesystem-related business logic
type FilesystemService interface {
	// PerformPhotoIndex initiates an index of image files on filesystem
	PerformPhotoIndex(save func(PhotoFile) error)
	// WriteFileToFilesystem writes a photo to the filesystem
	WriteFileToFilesystem(photo PhotoUpload) PhotoFile
	GenerateThumbnail(path string, exifData exif.Exif) []byte
}
