package service

import (
	b64 "encoding/base64"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
)

/**
 * FilesystemService implements port.FilesystemService interface
 * and provides an access to the utility repository
 */
type FilesystemService struct {
	fsRepo           port.FilesystemRepository
	jobSvc           port.JobService
	utilitySvc       port.UtilityService
	indexWorkers     int
	thumbnailStorage string
}

// NewFilesystemService creates a new filesystem service instance
func NewFilesystemService(fsRepo port.FilesystemRepository, jobSvc port.JobService, utilitySvc port.UtilityService, indexWorkers int, thumbnailStorage string) *FilesystemService {
	return &FilesystemService{
		fsRepo,
		jobSvc,
		utilitySvc,
		indexWorkers,
		thumbnailStorage,
	}
}

func (fss *FilesystemService) PerformPhotoIndex(ctx context.Context, save func([]domain.PhotoFile) error, indexCache map[string]struct{ FileHash string; FileModifiedTime int64 }) {
	_ = fss.jobSvc.JobStart("Photo_index")
	slog.Info("Starting photo index")

	defer func(jobName string) {
		_ = fss.jobSvc.JobComplete(jobName)
		slog.Info("Finished photo index")
	}("Photo_index")

	// Create buffered channel for photo paths
	numWorkers := fss.indexWorkers
	if numWorkers < 1 {
		numWorkers = 1
	}
	photoChan := make(chan string, numWorkers*2)

	// Create worker pool
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			slog.Info("Worker started", "workerID", workerID)
			batch := make([]domain.PhotoFile, 0, 100)
			for path := range photoChan {
				select {
				case <-ctx.Done():
					slog.Info("Worker cancelled", "workerID", workerID)
					// Drain remaining items without processing
					for range photoChan {
					}
					return
				default:
				}
				slog.Debug("Worker processing photo", "workerID", workerID, "path", path)
				photoFileInfo, err := os.Lstat(path)
				if err != nil {
					slog.Error("Unable to process photo", "workerID", workerID, "path", path, "error", err)
					continue
				}

			// Check index cache to skip unchanged files (unless valkey mode —
			// thumbnails may need generation even when file content is unchanged)
			photoId := b64.StdEncoding.EncodeToString([]byte(path))
			if fss.thumbnailStorage != "valkey" {
				if cached, exists := indexCache[photoId]; exists {
					if photoFileInfo.ModTime().Unix() == cached.FileModifiedTime {
						slog.Debug("Skipping unchanged photo", "path", path)
						continue
					}
				}
			}

				photo := fss.getMetaData(path, photoFileInfo.Name(), photoFileInfo)
				batch = append(batch, photo)
				if len(batch) >= 100 {
					if err := save(batch); err != nil {
						slog.Error("Failed to save batch", "error", err)
					}
					batch = batch[:0]
				}
			}
			// Flush remaining items
			if len(batch) > 0 {
				if err := save(batch); err != nil {
					slog.Error("Failed to save final batch", "error", err)
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

	basePath, _ := fss.utilitySvc.GetSetting("default_new_albums_dir")
	fileSavePath := fmt.Sprintf("%s/%s/%s", basePath.Value, photo.AlbumName, photo.Name)
	slog.Info("Saving photo", "path", fileSavePath)

	fss.fsRepo.CreateDirectoryIfNotExists(basePath.Value, photo.AlbumName)
	value := strings.Split(photo.BinaryContent, ",")

	decodedData, _ := b64.StdEncoding.DecodeString(value[1])
	err := os.WriteFile(fileSavePath, decodedData, 0644)
	if err != nil {
		slog.Error("Unable to save photo from upload", "error", err)
	}

	fileInfo, _ := os.Lstat(fileSavePath)

	return fss.getMetaData(fileSavePath, fileInfo.Name(), fileInfo)

}

// GetPhotoMetadata extracts metadata for an existing file on disk, reusing the
// same extraction logic as the indexer. Used by third-party importers (e.g.
// Google Takeout) that stage files into PhotoDir before saving them to the DB.
func (fss *FilesystemService) GetPhotoMetadata(path string) domain.PhotoFile {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		slog.Warn("GetPhotoMetadata: unable to stat file", "path", path, "error", err)
		return domain.PhotoFile{Path: path}
	}
	return fss.getMetaData(path, fileInfo.Name(), fileInfo)
}

func (fss *FilesystemService) getMetaData(path string, name string, fileInfo os.FileInfo) domain.PhotoFile {
	slashIndices := utils.AllIndicesOfChar(path, "/")
	photoDirectory := path[slashIndices[len(slashIndices)-2]+1 : slashIndices[len(slashIndices)-1]]

	mediaType := "image"
	duration := 0
	width := 0
	height := 0
	if utils.IsVideoFile(name) {
		mediaType = "video"
		dur, w, h, _, err := utils.GetVideoMetadata(path)
		if err == nil {
			duration = dur
			width = w
			height = h
		}
	}

	// Open file once for all operations
	file, err := os.Open(path)
	if err != nil {
		slog.Error("Unable to open file", "error", err)
		// Return minimal PhotoFile if file cannot be opened
		return domain.PhotoFile{
			Path:         path,
			Directory:    photoDirectory,
			Size:         fileInfo.Size(),
			Extension:    utils.GetExtension(path),
			Name:         name,
			MediaType:    mediaType,
			Duration:     duration,
			Width:        width,
			Height:       height,
			ModifiedTime: fileInfo.ModTime().Unix(),
			Latitude:     0,
			Longitude:    0,
		}
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	// For images: get dimensions from the open file
	if mediaType == "image" {
		cfg, _, err := image.DecodeConfig(file)
		if err == nil {
			width = cfg.Width
			height = cfg.Height
		}
		// Reset file position for subsequent reads
		_, _ = file.Seek(0, 0)
	}

	// Get MIME type by reading first 512 bytes
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)
	mimeType := http.DetectContentType(buffer[:n])
	_, _ = file.Seek(0, 0)

	// Get MD5 sum
	md5Sum, err := utils.GetSum(file)
	if err != nil {
		slog.Error("Unable to calculate MD5 sum", "path", path, "error", err)
		md5Sum = ""
	}
	_, _ = file.Seek(0, 0)

	// Get EXIF data
	exifData := utils.GetExifData(file)

	// Extract GPS coordinates from EXIF
	lat, lng := 0.0, 0.0
	if latTag, err := exifData.Get(exif.GPSLatitude); err == nil {
		if lngTag, err := exifData.Get(exif.GPSLongitude); err == nil {
			lat = convertGPSCoordinate(latTag)
			lng = convertGPSCoordinate(lngTag)
			// Check GPSLatitudeRef and GPSLongitudeRef for negative values
			if ref, err := exifData.Get(exif.GPSLatitudeRef); err == nil {
				if refStr, _ := ref.StringVal(); refStr == "S" {
					lat = -lat
				}
			}
			if ref, err := exifData.Get(exif.GPSLongitudeRef); err == nil {
				if refStr, _ := ref.StringVal(); refStr == "W" {
					lng = -lng
				}
			}
		}
	}

	return domain.PhotoFile{
		MD5:          md5Sum,
		Path:         path,
		Directory:    photoDirectory,
		Size:         fileInfo.Size(),
		Extension:    utils.GetExtension(path),
		Name:         name,
		Exif:         exifData,
		Mime:         mimeType,
		MediaType:    mediaType,
		Duration:     duration,
		Width:        width,
		Height:       height,
		ModifiedTime: fileInfo.ModTime().Unix(),
		Latitude:     lat,
		Longitude:    lng,
	}
}

func (fss *FilesystemService) GenerateThumbnail(path string, exifData exif.Exif, width, height int) []byte {
	if utils.IsVideoFile(path) {
		return fss.fsRepo.GenerateThumbnail(path, width, height)
	}
	parsedThumbnail, _ := exifData.JpegThumbnail()
	if len(parsedThumbnail) == 0 {
		parsedThumbnail = fss.fsRepo.GenerateThumbnail(path, width, height)
	}
	return parsedThumbnail
}

func (fss *FilesystemService) MoveToTrash(path string) (string, error) {
	return fss.fsRepo.MoveToTrash(path)
}

func (fss *FilesystemService) RestoreFromTrash(trashPath, originalPath string) error {
	return fss.fsRepo.RestoreFromTrash(trashPath, originalPath)
}

func (fss *FilesystemService) RenameDirectory(oldPath, newPath string) error {
	return fss.fsRepo.RenameDirectory(oldPath, newPath)
}

func (fss *FilesystemService) PermanentlyDelete(path string) error {
	return fss.fsRepo.PermanentlyDelete(path)
}

func (fss *FilesystemService) PermanentlyDeleteTrashFile(originalPath string) error {
	return fss.fsRepo.PermanentlyDeleteTrashFile(originalPath)
}

// convertGPSCoordinate converts an EXIF GPS rational tag to decimal degrees
func convertGPSCoordinate(tag *tiff.Tag) float64 {
	// GPS coordinates are stored as 3 rational values: degrees, minutes, seconds
	num, den, err := tag.Rat2(0)
	if err != nil {
		return 0
	}
	degrees := float64(num) / float64(den)

	num, den, err = tag.Rat2(1)
	if err != nil {
		return 0
	}
	minutes := float64(num) / float64(den)

	num, den, err = tag.Rat2(2)
	if err != nil {
		return 0
	}
	seconds := float64(num) / float64(den)

	return degrees + minutes/60 + seconds/3600
}
