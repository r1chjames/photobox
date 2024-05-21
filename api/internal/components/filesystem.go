package components

import (
	"bytes"
	"crypto/md5"
	b64 "encoding/base64"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"gitlab.com/r1chjames/photobox/api/internal/database"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gitlab.com/r1chjames/photobox/api/internal/utils"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func PerformPhotoIndex(appConfig AppConfig, dbEnv *database.Env) {
	_ = dbEnv.JobStarting("Photo_index")
	log.Print("Starting photo index")

	defer func(dbEnv *database.Env, jobName string) {
		_ = dbEnv.JobCompleted(jobName)
		log.Print("Finished photo index")
	}(dbEnv, "Photo_index")

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
			dbEnv.SavePhotoRecordToDatabase(photo)
		}
	}(photoChan)
	ScanFilesystem(appConfig, photoChan)
}

func createDirectoryIfNotExists(basePhotoPath string, directoryName string) {
	fullPath := fmt.Sprintf("%s/%s", basePhotoPath, directoryName)
	err := os.Mkdir(fullPath, os.ModePerm) //TODO check if exists, swallow error if so
	if err != nil {
		log.Printf("Unable to create album folder. Check the value of setting default_new_albums_dir exists and is writable, %s", err)
	}
}

func WriteFileToFilesystem(dbEnv *database.Env, photo PhotoUpload) PhotoFile {

	if !isImageFile(photo.Name) {
	}

	basePath, _ := dbEnv.GetSetting("default_new_albums_dir")
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

func ScanFilesystem(appConfig AppConfig, photoChan chan string) {

	photosRoot := appConfig.PhotoDir

	wg.Add(1)
	walkDir(photosRoot, photoChan)
	wg.Wait()
}

func walkDir(dir string, photoChan chan string) {
	defer wg.Done()

	visit := func(path string, d os.DirEntry, err error) error {
		if d.IsDir() && path != dir {
			log.Printf("Processing directory: %s", d.Name())
			wg.Add(1)
			go walkDir(path, photoChan)
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

func isImageFile(fileName string) bool {
	imageFileTypes := []string{"JPG", "JPEG", "PNG", "TIFF", "DNG", "RAW"}
	fileType := strings.ToUpper(getFileExtension(fileName))
	return utils.Exists(imageFileTypes, fileType)
}

func getMetaData(path string, name string, size int64) PhotoFile {
	slashIndices := utils.AllIndicesOfChar(path, "/")
	photoDirectory := path[slashIndices[len(slashIndices)-2]+1 : slashIndices[len(slashIndices)-1]]
	exifData, thumbnail := getExifDataAndThumbnail(path)
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

func getExifDataAndThumbnail(path string) (exif.Exif, []byte) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to open file, %s", err)
		return exif.Exif{}, []byte{}
	}
	var exifData *exif.Exif
	var parsedThumbnail []byte

	exifData, err = exif.Decode(file)
	if err != nil {
		exifData = &exif.Exif{}
	}

	parsedThumbnail = getThumbnail(exifData)
	if len(parsedThumbnail) == 0 {
		parsedThumbnail = generateMissingThumbnail(path)
	}
	return *exifData, parsedThumbnail
}

func getThumbnail(exif *exif.Exif) []byte {
	thumbnail, _ := exif.JpegThumbnail()
	return thumbnail
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

func getFileExtension(path string) string {
	dotIndex := strings.LastIndex(path, ".")
	return path[dotIndex+1:]
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
