package port

import (
	"github.com/rwcarlsen/goexif/exif"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
)

//go:generate mockgen -source=filesystem.go -destination=mock/filesystem.go -package=mock

// FilesystemRepository is an interface for interacting with filesystem-related business logic
type FilesystemRepository interface {
	// ScanFilesystem scans a directory on the filesystem
	ScanFilesystem(photoChan chan string)
	GenerateThumbnail(path string, width, height int) []byte
	CreateDirectoryIfNotExists(basePhotoPath string, directoryName string)
	MoveToTrash(path string) (string, error)
	RestoreFromTrash(trashPath, originalPath string) error
	RenameDirectory(oldPath, newPath string) error
}

// FilesystemService is an interface for interacting with filesystem-related business logic
type FilesystemService interface {
	// PerformPhotoIndex initiates an index of image files on filesystem
	PerformPhotoIndex(save func([]domain.PhotoFile) error, indexCache map[string]struct{FileHash string; FileModifiedTime int64})
	// WriteFileToFilesystem writes a photo to the filesystem
	WriteFileToFilesystem(photo domain.PhotoUpload) domain.PhotoFile
	GenerateThumbnail(path string, exifData exif.Exif, width, height int) []byte
	MoveToTrash(path string) (string, error)
	RestoreFromTrash(trashPath, originalPath string) error
	RenameDirectory(oldPath, newPath string) error
}
