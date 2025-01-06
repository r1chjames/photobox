package repository

import (
	"bytes"
	"crypto/md5"
	b64 "encoding/base64"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

func (fs *FilesystemRepository) createDirectoryIfNotExists(basePhotoPath string, directoryName string) {
	fullPath := fmt.Sprintf("%s/%s", basePhotoPath, directoryName)
	err := os.Mkdir(fullPath, os.ModePerm) //TODO check if exists, swallow error if so
	if err != nil {
		log.Printf("Unable to create album folder. Check the value of setting default_new_albums_dir exists and is writable, %s", err)
	}
}

func (fs *FilesystemRepository) WriteFileToFilesystem(photo PhotoUpload) PhotoFile {

	if !isImageFile(photo.Name) {
	}

	basePath, _ := fs.GetSetting("default_new_albums_dir")
	fileSavePath := fmt.Sprintf("%s/%s/%s", basePath.Value, photo.AlbumName, photo.Name)
	log.Printf("Saving photo to: %s", fileSavePath)

	createDirectoryIfNotExists(basePath.Value, photo.AlbumName)
	value := strings.Split(photo.BinaryContent, ",")

	decodedData, err := b64.StdEncoding.DecodeString(value[1])
	err = os.WriteFile(fileSavePath, decodedData, 0644)
	if err != nil {
		log.Print("Unable to save photo from upload")
	}

	fileInfo, _ := os.Lstat(fileSavePath)

	return getMetaData(fileSavePath, fileInfo.Name(), fileInfo.Size())

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

		if d.Type().IsRegular() && isImageFile(d.Name()) {
			photoChan <- path
		}
		return nil
	}

	err := filepath.WalkDir(dir, visit)
	if err != nil {
		log.Print(err)
	}
}

func getExifData(path string) exif.Exif {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to open file, %s", err)
		return exif.Exif{}
	}
	var exifData *exif.Exif

	exifData, err = exif.Decode(file)
	if err != nil {
		exifData = &exif.Exif{}
	}

	return *exifData
}

func generateThumbnail(path string, exifData exif.Exif) []byte {
	parsedThumbnail := getThumbnail(&exifData)
	if len(parsedThumbnail) == 0 {
		parsedThumbnail = generateMissingThumbnail(path)
	}
	return parsedThumbnail
}

func getThumbnail(exif *exif.Exif) []byte {
	thumbnail, _ := exif.JpegThumbnail()
	return thumbnail
}

func getFileExtension(path string) string {
	dotIndex := strings.LastIndex(path, ".")
	return path[dotIndex+1:]
}

func isImageFile(fileName string) bool {
	imageFileTypes := []string{"JPG", "JPEG", "PNG", "TIFF", "DNG", "RAW"}
	fileType := strings.ToUpper(getFileExtension(fileName))
	return utils.Exists(imageFileTypes, fileType)
}

func generateMissingThumbnail(path string) []byte {
	extension := getFileExtension(path)
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

func getFileType(path string) string {
	out, err := exec.Command("file", "--brief", "--mime-type", path).Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

func getExtension(path string) string {
	return filepath.Ext(path)
}

func getSize(info os.FileInfo) int64 {
	return info.Size()
}

func getSum(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to generate checksum, %s", err)
		return ""
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)) // TODO Optimize this, if possible
}
