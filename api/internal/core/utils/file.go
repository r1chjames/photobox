package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetFileExtension(path string) string {
	dotIndex := strings.LastIndex(path, ".")
	return path[dotIndex+1:]
}

func IsImageFile(fileName string) bool {
	imageFileTypes := []string{"JPG", "JPEG", "PNG", "TIFF", "DNG", "RAW"}
	fileType := strings.ToUpper(GetFileExtension(fileName))
	return Exists(imageFileTypes, fileType)
}

func OpenFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetExifData(file *os.File) exif.Exif {
	var exifData *exif.Exif

	exifData, err := exif.Decode(file)
	if err != nil {
		exifData = &exif.Exif{}
	}

	return *exifData
}

func GetFileType(path string) string {
	out, err := exec.Command("file", "--brief", "--mime-type", path).Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

func GetExtension(path string) string {
	return filepath.Ext(path)
}

func GetSize(info os.FileInfo) int64 {
	return info.Size()
}

func GetSum(file *os.File) string {
	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)) // TODO Optimize this, if possible
}
