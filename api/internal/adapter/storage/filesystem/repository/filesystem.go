package repository

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/deepteams/webp"
	"github.com/disintegration/imaging"
	jpegscaled "github.com/m8rge/go-scaled-jpeg"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
)

type FilesystemRepository struct {
	wg           sync.WaitGroup
	jobSvc       *service.JobService
	config       appconfig.AppConfig
	dirSem       chan struct{} // Semaphore to limit concurrent directory walking
	ffmpegFound   bool
	exifToolFound bool
}

func NewFilesystemRepository(config appconfig.AppConfig, jobService *service.JobService) *FilesystemRepository {
	_, ffmpegErr := exec.LookPath("ffmpeg")
	hasFfmpeg := ffmpegErr == nil
	if !hasFfmpeg {
		slog.Warn("ffmpeg not found in PATH, video thumbnails will be skipped")
	}
	_, exifErr := exec.LookPath("exiftool")
	hasExifTool := exifErr == nil
	if !hasExifTool {
		slog.Info("exiftool not found in PATH, RAW photo thumbnails will use EXIF preview only")
	}
	return &FilesystemRepository{
		wg:           sync.WaitGroup{},
		jobSvc:       jobService,
		config:       config,
		ffmpegFound:  hasFfmpeg,
		exifToolFound: hasExifTool,
	}
}

func (fs *FilesystemRepository) CreateDirectoryIfNotExists(basePhotoPath string, directoryName string) {
	fullPath := fmt.Sprintf("%s/%s", basePhotoPath, directoryName)
	err := os.Mkdir(fullPath, os.ModePerm)
	if err != nil {
		slog.Warn("Unable to create album folder, check default_new_albums_dir setting", "path", fullPath, "error", err)
	}
}

func (fs *FilesystemRepository) ScanFilesystem(photoChan chan string) {
	photosRoot := fs.config.PhotoDir

	// Initialize semaphore with limit of 4x CPU count for concurrent directory walkers
	maxConcurrentDirs := runtime.NumCPU() * 4
	fs.dirSem = make(chan struct{}, maxConcurrentDirs)

	fs.wg.Add(1)
	fs.walkDir(photosRoot, photoChan)
	fs.wg.Wait()
}

func (fs *FilesystemRepository) walkDir(dir string, photoChan chan string) {
	// Acquire semaphore slot
	fs.dirSem <- struct{}{}
	defer func() {
		<-fs.dirSem // Release semaphore slot
		fs.wg.Done()
	}()

	visit := func(path string, d os.DirEntry, err error) error {
		if d.IsDir() && path != dir {
			slog.Debug("Processing directory", "name", d.Name())
			fs.wg.Add(1)
			go fs.walkDir(path, photoChan)
			return filepath.SkipDir
		}

		if d.Type().IsRegular() && utils.IsMediaFile(d.Name()) {
			photoChan <- path
		}
		return nil
	}

	err := filepath.WalkDir(dir, visit)
	if err != nil {
		slog.Error("Error walking directory", "dir", dir, "error", err)
	}
}

func (fs *FilesystemRepository) GenerateThumbnail(path string, width, height int) []byte {
	if utils.IsVideoFile(path) {
		return fs.generateVideoThumbnail(path, width, height)
	}
	// Check if source file still exists before attempting decode
	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Error("Source file missing, cannot generate thumbnail", "path", path)
		return nil
	}

	img, err := fs.decodeImage(path, width, height)
	if err != nil {
		if strings.Contains(err.Error(), "unsupported feature") {
			slog.Warn("TIFF color model not supported, falling back to RAW preview", "path", path, "error", err)
		} else {
			slog.Warn("Unable to decode image for thumbnail", "path", path, "error", err)
		}
		return fs.generateRawThumbnail(path, width, height)
	}
	thumb := imaging.Thumbnail(img, width, height, imaging.CatmullRom)
	var buffer bytes.Buffer
	if err := webp.Encode(&buffer, thumb, webp.OptionsForPreset(webp.PresetPhoto, 75)); err != nil {
		slog.Error("WebP thumbnail encode failed", "path", path, "error", err)
		return nil
	}
	return buffer.Bytes()
}

// decodeImage opens and decodes an image, using a memory-efficient path for
// large JPEGs (>20MP) via scaled DCT decoding.
func (fs *FilesystemRepository) decodeImage(path string, targetWidth, targetHeight int) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return nil, err
	}

	const largeImageThreshold = 20_000_000 // 20 megapixels
	isLargeJPEG := format == "jpeg" && cfg.Width*cfg.Height > largeImageThreshold

	if isLargeJPEG {
		// Calculate DCT scale: 8=full, 4=1/2, 2=1/4, 1=1/8.
		// Pick the smallest scale where decoded size >= target * 2 to avoid upscaling.
		scale := 8
		for _, s := range []int{1, 2, 4, 8} {
			decodedW := cfg.Width * s / 8
			if decodedW >= targetWidth*2 {
				scale = s
				break
			}
		}
		if _, err := f.Seek(0, 0); err != nil {
			return nil, err
		}
		img, err := jpegscaled.Decode(f, jpegscaled.DecodeOptions{DCTSizeScaled: scale})
		if err == nil {
			slog.Debug("Using scaled JPEG decode", "path", path, "scale", scale, "dims", fmt.Sprintf("%dx%d", cfg.Width, cfg.Height))
			return img, nil
		}
		slog.Warn("Scaled JPEG decode failed, falling back to full decode", "path", path, "error", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}
	img, _, err := image.Decode(f)
	return img, err
}

func (fs *FilesystemRepository) generateVideoThumbnail(path string, width, height int) []byte {
	if !fs.ffmpegFound {
		return nil
	}
	cmd := exec.Command("ffmpeg", "-i", path, "-ss", "00:00:01", "-vframes", "1", "-f", "image2pipe", "-vcodec", "png", "-")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		slog.Error("Unable to generate video thumbnail", "path", path, "error", err)
		return nil
	}
	img, _, err := image.Decode(bytes.NewReader(out.Bytes()))
	if err != nil {
		slog.Error("Unable to decode video frame", "path", path, "error", err)
		return nil
	}
	thumb := imaging.Thumbnail(img, width, height, imaging.CatmullRom)
	var buffer bytes.Buffer
	if err := webp.Encode(&buffer, thumb, webp.OptionsForPreset(webp.PresetPhoto, 75)); err != nil {
		slog.Error("WebP video thumbnail encode failed", "path", path, "error", err)
		return nil
	}
	return buffer.Bytes()
}

func (fs *FilesystemRepository) generateRawThumbnail(path string, width, height int) []byte {
	if !fs.exifToolFound {
		return nil
	}

	// exiftool extracts camera-embedded JPEG previews from RAW/DNG files.
	// Try common tag names across formats (Canon CR2, Nikon NEF, Sony ARW, DNG, etc.).
	tags := []string{"-PreviewImage", "-JpgFromRaw", "-ThumbnailImage"}
	for _, tag := range tags {
		cmd := exec.Command("exiftool", "-b", tag, path)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil || out.Len() == 0 {
			continue
		}

		img, _, err := image.Decode(&out)
		if err != nil {
			continue
		}

		thumb := imaging.Thumbnail(img, width, height, imaging.CatmullRom)
		var buf bytes.Buffer
		if err := webp.Encode(&buf, thumb, webp.OptionsForPreset(webp.PresetPhoto, 75)); err != nil {
			slog.Error("WebP RAW thumbnail encode failed", "path", path, "error", err)
			continue
		}
		return buf.Bytes()
	}

	slog.Debug("exiftool found no preview in RAW file", "path", path)
	return nil
}

func (fs *FilesystemRepository) MoveToTrash(path string) (string, error) {
	trashDir := filepath.Join(fs.config.PhotoDir, ".trash")
	if err := os.MkdirAll(trashDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create trash directory: %w", err)
	}
	filename := filepath.Base(path)
	trashPath := filepath.Join(trashDir, filename)
	// Avoid overwriting existing trashed files
	if _, err := os.Stat(trashPath); err == nil {
		trashPath = filepath.Join(trashDir, fmt.Sprintf("%d_%s", time.Now().UnixMilli(), filename))
	}
	if err := os.Rename(path, trashPath); err != nil {
		return "", fmt.Errorf("failed to move to trash: %w", err)
	}
	return trashPath, nil
}

func (fs *FilesystemRepository) RestoreFromTrash(trashPath, originalPath string) error {
	if err := os.MkdirAll(filepath.Dir(originalPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create original directory: %w", err)
	}
	if err := os.Rename(trashPath, originalPath); err != nil {
		return fmt.Errorf("failed to restore from trash: %w", err)
	}
	return nil
}

func (fs *FilesystemRepository) RenameDirectory(oldPath, newPath string) error {
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename directory: %w", err)
	}
	return nil
}

// PermanentlyDelete removes a file from disk (used when purging trashed
// originals). Returns nil if the file does not exist.
func (fs *FilesystemRepository) PermanentlyDelete(path string) error {
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to permanently delete file: %w", err)
	}
	return nil
}

// PermanentlyDeleteTrashFile removes every trashed copy of an original file
// from the .trash directory. MoveToTrash may have stored the file as either
// <name> or <millis>_<name> when collisions occurred, so all variants are
// removed. Missing files are not an error.
func (fs *FilesystemRepository) PermanentlyDeleteTrashFile(originalPath string) error {
	if originalPath == "" {
		return nil
	}
	trashDir := filepath.Join(fs.config.PhotoDir, ".trash")
	filename := filepath.Base(originalPath)

	entries, err := os.ReadDir(trashDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read trash directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if name == filename || strings.HasSuffix(name, "_"+filename) {
			_ = os.Remove(filepath.Join(trashDir, name))
		}
	}
	return nil
}
