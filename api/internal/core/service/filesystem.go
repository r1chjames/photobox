package service

import (
	b64 "encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/rwcarlsen/goexif/exif"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
)

/**
 * FilesystemService implements port.FilesystemService interface
 * and provides an access to the utility repository
 */
type FilesystemService struct {
	fsRepo     port.FilesystemRepository
	jobSvc     port.JobService
	utilitySvc port.UtilityService
}

// NewFilesystemService creates a new filesystem service instance
func NewFilesystemService(fsRepo port.FilesystemRepository, jobSvc port.JobService, utilitySvc port.UtilityService) *FilesystemService {
	return &FilesystemService{
		fsRepo,
		jobSvc,
		utilitySvc,
	}
}

func (fss *FilesystemService) PerformPhotoIndex(save func(domain.PhotoFile) error) {
	_ = fss.jobSvc.JobStart("Photo_index")
	slog.Info("Starting photo index")

	defer func(jobName string) {
		_ = fss.jobSvc.JobComplete(jobName)
		slog.Info("Finished photo index")
	}("Photo_index")

	// Create buffered channel for photo paths
	numWorkers := runtime.NumCPU()
	photoChan := make(chan string, numWorkers*2)

	// Create worker pool
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			slog.Info("Worker started", "workerID", workerID)
			for path := range photoChan {
				slog.Debug("Worker processing photo", "workerID", workerID, "path", path)
				photoFile, err := os.Lstat(path)
				if err != nil {
					slog.Error("Unable to process photo", "workerID", workerID, "path", path, "error", err)
					continue
				}
				photo := fss.getMetaData(path, photoFile.Name(), photoFile.Size())
				err = save(photo)
				if err != nil {
					slog.Error("Unable to save photo", "workerID", workerID, "photo", photo.Name, "error", err)
				}
			}
			slog.Info("Worker finished", "workerID", workerID)
		}(i)
	}

	// Scan filesystem and send paths to workers
	fss.fsRepo.ScanFilesystem(photoChan)

	// Close channel to signal workers to stop
	close(photoChan)

	// Wait for all workers to finish
	wg.Wait()
}

func (fss *FilesystemService) WriteFileToFilesystem(photo domain.PhotoUpload) domain.PhotoFile {

	if !utils.IsImageFile(photo.Name) {
	}

	basePath, _ := fss.utilitySvc.GetSetting("default_new_albums_dir")
	fileSavePath := fmt.Sprintf("%s/%s/%s", basePath.Value, photo.AlbumName, photo.Name)
	slog.Info("Saving photo", "path", fileSavePath)

	fss.fsRepo.CreateDirectoryIfNotExists(basePath.Value, photo.AlbumName)
	value := strings.Split(photo.BinaryContent, ",")

	decodedData, err := b64.StdEncoding.DecodeString(value[1])
	err = os.WriteFile(fileSavePath, decodedData, 0644)
	if err != nil {
		slog.Error("Unable to save photo from upload", "error", err)
	}

	fileInfo, _ := os.Lstat(fileSavePath)

	return fss.getMetaData(fileSavePath, fileInfo.Name(), fileInfo.Size())

}

func (fss *FilesystemService) getMetaData(path string, name string, size int64) domain.PhotoFile {
	slashIndices := utils.AllIndicesOfChar(path, "/")
	photoDirectory := path[slashIndices[len(slashIndices)-2]+1 : slashIndices[len(slashIndices)-1]]

	file, err := utils.OpenFile(path)
	if err != nil {
		slog.Error("Unable to open file", "error", err)
		// Return empty domain.PhotoFile if file cannot be opened
		return domain.PhotoFile{
			Path:      path,
			Directory: photoDirectory,
			Size:      size,
			Extension: utils.GetExtension(path),
			Name:      name,
		}
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	md5Sum, err := utils.GetSum(file)
	if err != nil {
		slog.Error("Unable to calculate MD5 sum", "path", path, "error", err)
		md5Sum = ""
	}

	mimeType, err := utils.GetFileType(path)
	if err != nil {
		slog.Error("Unable to get file type", "path", path, "error", err)
		mimeType = ""
	}

	exifData := utils.GetExifData(file)
	thumbnail := fss.GenerateThumbnail(path, exifData)
	return domain.PhotoFile{
		MD5:       md5Sum,
		Path:      path,
		Directory: photoDirectory,
		Size:      size,
		Extension: utils.GetExtension(path),
		Name:      name,
		Exif:      exifData,
		Mime:      mimeType,
		Thumbnail: thumbnail,
	}
}

func (fss *FilesystemService) GenerateThumbnail(path string, exifData exif.Exif) []byte {
	parsedThumbnail, _ := exifData.JpegThumbnail()
	if len(parsedThumbnail) == 0 {
		parsedThumbnail = fss.fsRepo.GenerateThumbnail(path)
	}
	return parsedThumbnail
}
