package testutil

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// CreateTestPhotoDir creates a temporary directory with test photos
func CreateTestPhotoDir(t *testing.T) string {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "photobox-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	return tempDir
}

// CleanupTestPhotoDir removes the test photo directory
func CleanupTestPhotoDir(t *testing.T, dir string) {
	t.Helper()
	os.RemoveAll(dir)
}

// CreateTestImage creates a test image file
func CreateTestImage(t *testing.T, path string, width, height int, format string) {
	t.Helper()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	// Create a simple colored image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := color.RGBA{
				R: uint8((x * 255) / width),
				G: uint8((y * 255) / height),
				B: 128,
				A: 255,
			}
			img.Set(x, y, c)
		}
	}

	// Create file
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create image file %s: %v", path, err)
	}
	defer file.Close()

	// Encode based on format
	switch format {
	case "jpg", "jpeg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(file, img)
	default:
		t.Fatalf("Unsupported image format: %s", format)
	}

	if err != nil {
		t.Fatalf("Failed to encode image %s: %v", path, err)
	}
}

// CreateTestImageWithEXIF creates a test JPEG image with EXIF data
// Note: For full EXIF support, you'd need to use a library like github.com/rwcarlsen/goexif
func CreateTestImageWithEXIF(t *testing.T, path string, width, height int) {
	t.Helper()
	// For now, create a basic JPEG
	// In a real implementation, you'd add EXIF metadata here
	CreateTestImage(t, path, width, height, "jpg")
}

// CreateTestPhotoStructure creates a realistic directory structure with test photos
func CreateTestPhotoStructure(t *testing.T, baseDir string) {
	t.Helper()

	// Create directory structure
	subdirs := []string{
		"Album1",
		"Album2",
		"Vacation/2024",
		"Vacation/2023",
	}

	for _, subdir := range subdirs {
		path := filepath.Join(baseDir, subdir)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", path, err)
		}
	}

	// Create test images in different directories
	testImages := []struct {
		path   string
		width  int
		height int
		format string
	}{
		{"Album1/photo1.jpg", 1920, 1080, "jpg"},
		{"Album1/photo2.jpg", 1920, 1080, "jpg"},
		{"Album1/photo3.png", 1024, 768, "png"},
		{"Album2/sunset.jpg", 3840, 2160, "jpg"},
		{"Album2/landscape.jpg", 2560, 1440, "jpg"},
		{"Vacation/2024/beach.jpg", 1920, 1080, "jpg"},
		{"Vacation/2024/mountain.jpg", 1920, 1080, "jpg"},
		{"Vacation/2023/city.jpg", 1280, 720, "jpg"},
		{"random_photo.jpg", 800, 600, "jpg"},
	}

	for _, img := range testImages {
		fullPath := filepath.Join(baseDir, img.path)
		CreateTestImage(t, fullPath, img.width, img.height, img.format)
	}
}

// CountFilesInDir recursively counts files in a directory
func CountFilesInDir(t *testing.T, dir string) int {
	t.Helper()

	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to count files in %s: %v", dir, err)
	}

	return count
}
