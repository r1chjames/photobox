package repository

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FilesystemRepository struct {
	wg     sync.WaitGroup
	jobSvc *service.JobService
	config appconfig.AppConfig
}

func NewFilesystemRepository(config appconfig.AppConfig, jobService *service.JobService) *FilesystemRepository {
	return &FilesystemRepository{
		wg:     sync.WaitGroup{},
		jobSvc: jobService,
		config: config,
	}
}

func (fs *FilesystemRepository) CreateDirectoryIfNotExists(basePhotoPath string, directoryName string) {
	fullPath := fmt.Sprintf("%s/%s", basePhotoPath, directoryName)
	err := os.Mkdir(fullPath, os.ModePerm) //TODO check if exists, swallow error if so
	if err != nil {
		log.Printf("Unable to create album folder. Check the value of setting default_new_albums_dir exists and is writable, %s", err)
	}
}

func (fs *FilesystemRepository) ScanFilesystem(photoChan chan string) {
	photosRoot := fs.config.PhotoDir

	fs.wg.Add(1)
	fs.walkDir(photosRoot, photoChan)
	fs.wg.Wait()
}

func (fs *FilesystemRepository) walkDir(dir string, photoChan chan string) {
	defer fs.wg.Done()

	visit := func(path string, d os.DirEntry, err error) error {
		if d.IsDir() && path != dir {
			log.Printf("Processing directory: %s", d.Name())
			fs.wg.Add(1)
			go fs.walkDir(path, photoChan)
			return filepath.SkipDir
		}

		if d.Type().IsRegular() && utils.IsImageFile(d.Name()) {
			photoChan <- path
		}
		return nil
	}

	err := filepath.WalkDir(dir, visit)
	if err != nil {
		log.Print(err)
	}
}

func (fs *FilesystemRepository) GenerateThumbnail(path string) []byte {
	extension := utils.GetFileExtension(path)
	img, err := imaging.Open(path)
	if err != nil {
		log.Printf("Unable to open file, %s", err)
		return nil
	}
	thumb := imaging.Thumbnail(img, 600, 600, imaging.CatmullRom)
	var buffer bytes.Buffer
	writer := io.MultiWriter(&buffer)
	format, _ := imaging.FormatFromExtension(extension)
	_ = imaging.Encode(writer, thumb, format)
	return buffer.Bytes()
}
