package service

import (
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
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
	fsRepo port.FilesystemRepository
	jobSvc port.JobService
}

// NewFilesystemService creates a new filesystem service instance
func NewFilesystemService(fsRepo port.FilesystemRepository, jobSvc port.JobService) *FilesystemService {
	return &FilesystemService{
		fsRepo,
		jobSvc,
	}
}

func (fss *FilesystemService) PerformPhotoIndex() error {
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
			photo := getMetaData(path, photoFile.Name(), photoFile.Size())
			fss.fsRepo.WriteFileToFilesystem(photo)
		}
	}(photoChan)
	fss.fsRepo.ScanFilesystem(photoChan)
}

func getMetaData(path string, name string, size int64) PhotoFile {
	slashIndices := utils.AllIndicesOfChar(path, "/")
	photoDirectory := path[slashIndices[len(slashIndices)-2]+1 : slashIndices[len(slashIndices)-1]]
	exifData := getExifData(path)
	thumbnail := generateThumbnail(path, exifData)
	return PhotoFile{
		MD5:       getSum(path),
		Path:      path,
		Directory: photoDirectory,
		Size:      size,
		Extension: getExtension(path),
		Name:      name, //GetFileName(path)
		Exif:      exifData,
		Mime:      getFileType(path),
		Thumbnail: thumbnail,
	}
}
