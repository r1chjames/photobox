package service

import (
	b64 "encoding/base64"
	"fmt"
	"sync"

	"github.com/rwcarlsen/goexif/exif"
	. "gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
	"log"
	"os"
	"runtime"
	"strings"
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

func (fss *FilesystemService) PerformPhotoIndex(save func(PhotoFile) error) {
	_ = fss.jobSvc.JobStart("Photo_index")
	log.Print("Starting photo index")

	defer func(jobName string) {
		_ = fss.jobSvc.JobComplete(jobName)
		log.Print("Finished photo index")
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
			log.Printf("Worker %d started", workerID)
			for path := range photoChan {
				log.Printf("Worker %d processing photo: %s", workerID, path)
				photoFile, err := os.Lstat(path)
				if err != nil {
					log.Printf("Worker %d unable to process photo at path %s: %v", workerID, path, err)
					continue
				}
				photo := fss.getMetaData(path, photoFile.Name(), photoFile.Size())
				err = save(photo)
				if err != nil {
					log.Printf("Worker %d unable to save photo %s: %v", workerID, photo.Name, err)
				}
			}
			log.Printf("Worker %d finished", workerID)
		}(i)
	}

	// Scan filesystem and send paths to workers
	fss.fsRepo.ScanFilesystem(photoChan)

	// Close channel to signal workers to stop
	close(photoChan)

	// Wait for all workers to finish
	wg.Wait()
}

func (fss *FilesystemService) WriteFileToFilesystem(photo PhotoUpload) PhotoFile {

	if !utils.IsImageFile(photo.Name) {
	}

	basePath, _ := fss.utilitySvc.GetSetting("default_new_albums_dir")
	fileSavePath := fmt.Sprintf("%s/%s/%s", basePath.Value, photo.AlbumName, photo.Name)
	log.Printf("Saving photo to: %s", fileSavePath)

	fss.fsRepo.CreateDirectoryIfNotExists(basePath.Value, photo.AlbumName)
	value := strings.Split(photo.BinaryContent, ",")

	decodedData, err := b64.StdEncoding.DecodeString(value[1])
	err = os.WriteFile(fileSavePath, decodedData, 0644)
	if err != nil {
		log.Print("Unable to save photo from upload")
	}

	fileInfo, _ := os.Lstat(fileSavePath)

	return fss.getMetaData(fileSavePath, fileInfo.Name(), fileInfo.Size())

}

func (fss *FilesystemService) getMetaData(path string, name string, size int64) PhotoFile {
	slashIndices := utils.AllIndicesOfChar(path, "/")
	photoDirectory := path[slashIndices[len(slashIndices)-2]+1 : slashIndices[len(slashIndices)-1]]

	file, err := utils.OpenFile(path)
	if err != nil {
		log.Printf("Unable to open file, %s", err)
		// Return empty PhotoFile if file cannot be opened
		return PhotoFile{
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
		log.Printf("Unable to calculate MD5 sum for %s: %s", path, err)
		md5Sum = ""
	}

	mimeType, err := utils.GetFileType(path)
	if err != nil {
		log.Printf("Unable to get file type for %s: %s", path, err)
		mimeType = ""
	}

	exifData := utils.GetExifData(file)
	thumbnail := fss.GenerateThumbnail(path, exifData)
	return PhotoFile{
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
