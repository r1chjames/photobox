package port

//go:generate mockgen -source=filesystem.go -destination=mock/filesystem.go -package=mock

// FilesystemRepository is an interface for interacting with filesystem-related business logic
type FilesystemRepository interface {
	// WriteFileToFilesystem writes a photo to the filesystem
	WriteFileToFilesystem(photo PhotoUpload) error
	// ScanFilesystem scans a directory on the filesystem
	ScanFilesystem(photoChan chan string) error
}

// FilesystemService is an interface for interacting with filesystem-related business logic
type FilesystemService interface {
	// PerformPhotoIndex initiates an index of image files on filesystem
	PerformPhotoIndex() error
}
