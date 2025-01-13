package service

import (
	b64 "encoding/base64"
	"fmt"
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

	photoChan := make(chan string, runtime.NumCPU())
	defer close(photoChan)
	go func(photoChan chan string) {
		for path := range photoChan {
			log.Printf("Processing photo: %s", path)
			photoFile, err := os.Lstat(path)
			if err != nil {
				log.Printf("Unable to process photo at path %s", err)
			}
			photo := fss.getMetaData(path, photoFile.Name(), photoFile.Size())
			err = save(photo)
			if err != nil {
				log.Printf("Unable to save photo %s", err)
			}
		}
	}(photoChan)
	fss.fsRepo.ScanFilesystem(photoChan)
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
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	exifData := utils.GetExifData(file)
	thumbnail := fss.GenerateThumbnail(path, exifData)
	return PhotoFile{
		MD5:       utils.GetSum(file),
		Path:      path,
		Directory: photoDirectory,
		Size:      size,
		Extension: utils.GetExtension(path),
		Name:      name,
		Exif:      exifData,
		Mime:      utils.GetFileType(path),
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
