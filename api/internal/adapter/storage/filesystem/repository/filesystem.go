package repository

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/service"
	"gitlab.com/r1chjames/photobox/api/internal/core/utils"
)

type FilesystemRepository struct {
	wg     sync.WaitGroup
	jobSvc *service.JobService
	config appconfig.AppConfig
	dirSem chan struct{} // Semaphore to limit concurrent directory walking
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

		if d.Type().IsRegular() && utils.IsImageFile(d.Name()) {
			photoChan <- path
		}
		return nil
	}

	err := filepath.WalkDir(dir, visit)
	if err != nil {
		slog.Error("Error walking directory", "dir", dir, "error", err)
	}
}

func (fs *FilesystemRepository) GenerateThumbnail(path string) []byte {
	extension := utils.GetFileExtension(path)
	img, err := imaging.Open(path)
	if err != nil {
		slog.Error("Unable to open file for thumbnail", "path", path, "error", err)
		return nil
	}
	thumb := imaging.Thumbnail(img, 600, 600, imaging.CatmullRom)
	var buffer bytes.Buffer
	writer := io.MultiWriter(&buffer)
	format, _ := imaging.FormatFromExtension(extension)
	_ = imaging.Encode(writer, thumb, format)
	return buffer.Bytes()
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
