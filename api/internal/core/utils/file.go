package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
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

func GetFileType(path string) (string, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file for type detection %s: %w", path, err)
	}
	defer file.Close()

	// Read the first 512 bytes for content type detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file for type detection %s: %w", path, err)
	}

	// Detect content type using Go's native detector
	contentType := http.DetectContentType(buffer[:n])
	return contentType, nil
}

func GetExtension(path string) string {
	return filepath.Ext(path)
}

func GetSize(info os.FileInfo) int64 {
	return info.Size()
}

func GetSum(file *os.File) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("failed to calculate MD5 sum: %w", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil // TODO Optimize this, if possible
}
