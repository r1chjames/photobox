package components

import (
	"crypto/md5"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gitlab.com/r1chjames/photobox/api/internal/utils"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func ScanFilesystem(appConfig AppConfig) []PhotoFile {

	photosRoot := appConfig.PhotoDir

	var photos []PhotoFile
	photoChan := make(chan PhotoFile, 10)

	wg.Add(1)
	go func(photoChan chan PhotoFile) {
		for photo := range photoChan {
			photos = append(photos, photo)
		}
	}(photoChan)

	walkDir(photosRoot, photoChan)
	wg.Wait()

	close(photoChan)

	return photos
}

func walkDir(dir string, photoChan chan PhotoFile) {
	defer wg.Done()

	visit := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && path != dir {
			log.Printf("Processing directory: %s", f.Name())
			wg.Add(1)
			go walkDir(path, photoChan)
			return filepath.SkipDir
		}

		if f.Mode().IsRegular() && isImageFile(f.Name()) {
			log.Printf("Processing file: %s", path)
			data := getMetaData(path, f)
			photoChan <- data
		}
		return nil
	}

	err := filepath.Walk(dir, visit)
	if err != nil {
		log.Print(err)
	}
}

func isImageFile(fileName string) bool {
	imageFileTypes := []string{"JPG", "JPEG", "PNG", "TIFF", "DNG", "RAW"}
	fileType := strings.ToUpper(getFileExtension(fileName))
	return utils.Exists(imageFileTypes, fileType)
}

func getMetaData(path string, info os.FileInfo) PhotoFile {
	slashIndices := utils.AllIndicesOfChar(path, "/")
	photoDirectory := path[slashIndices[len(slashIndices)-2]+1:slashIndices[len(slashIndices)-1]]
	exifData, thumbnail := getExifDataAndThumbnail(path)
	return PhotoFile{
		MD5:       getSum(path),
		Path:      path,
		Directory: photoDirectory,
		Size:      getSize(info),
		Extension: getExtension(path),
		Name:      info.Name(), //GetFileName(path)
		Exif:      exifData,
		Mime:      getFileType(path),
		Thumbnail: thumbnail,
	}
}

func getExifDataAndThumbnail(path string) (exif.Exif, []byte) {
	f, err := os.Open(path)
	if err != nil {
		return exif.Exif{}, []byte{}
	}
	return getExifData(f), getThumbnail(f)
}

func getThumbnail(file *os.File) []byte {
	x, err := exif.Decode(file)
	if err != nil {
		return []byte{}
	}
	thumbnail, _ := x.JpegThumbnail()
	return thumbnail
}

func getExifData(file *os.File) exif.Exif {
	x, err := exif.Decode(file)
	if err != nil {
		return exif.Exif{}
	}
	return *x
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
		log.Fatal(err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)) // TODO Optimize this, if possible
}