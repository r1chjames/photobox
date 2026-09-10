package service

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	ws "gitlab.com/r1chjames/photobox/api/internal/components/websocket"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gitlab.com/r1chjames/photobox/api/internal/core/port"
)

const (
	defaultTakeoutAlbum    = "Google Takeout"
	importBatchSize        = 200
	progressBroadcastEvery = 20
)

// takeoutMediaExtensions are the file types recognised inside a Google Takeout export.
var takeoutMediaExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".heic": true, ".heif": true,
	".gif": true, ".webp": true,
	".mp4": true, ".mov": true, ".m4v": true, ".3gp": true, ".avi": true,
}

// takeoutVideoExtensions is the subset of media extensions that are video.
var takeoutVideoExtensions = map[string]bool{
	".mp4": true, ".mov": true, ".m4v": true, ".3gp": true, ".avi": true,
}

// TakeoutImporter ingests Google Takeout exports (a directory or a .zip archive)
// into the library. It is safe for concurrent reads of progress; at most one
// import may run at a time.
type TakeoutImporter struct {
	photoService  port.PhotoService
	filesystemSvc port.FilesystemService
	photoRepo     port.PhotoRepository
	config        *appconfig.AppConfig
	wsHub         *ws.Hub

	mu       sync.Mutex
	progress domain.TakeoutImportProgress
	running  bool
}

// NewTakeoutImporter constructs a TakeoutImporter. wsHub may be nil (progress
// is still tracked and available via GetProgress; no events are broadcast).
func NewTakeoutImporter(photoSvc port.PhotoService, filesystemSvc port.FilesystemService, photoRepo port.PhotoRepository, config *appconfig.AppConfig, wsHub *ws.Hub) *TakeoutImporter {
	return &TakeoutImporter{
		photoService:  photoSvc,
		filesystemSvc: filesystemSvc,
		photoRepo:     photoRepo,
		config:        config,
		wsHub:         wsHub,
	}
}

// takeoutItem is a single media file discovered during a scan.
type takeoutItem struct {
	srcPath string
	album   string
	size    int64
	isVideo bool
	sidecar *domain.TakeoutSidecarJSON
}

// GetProgress returns a snapshot of the current (or last) import progress.
func (ti *TakeoutImporter) GetProgress() domain.TakeoutImportProgress {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	return ti.progress
}

// ScanTakeout inspects an export without modifying the library and reports
// what an import would do: totals, size, album layout, and duplicate counts.
func (ti *TakeoutImporter) ScanTakeout(path string) (*domain.TakeoutScanResult, error) {
	rootDir, cleanup, err := ti.openTakeout(path)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	items, warnings, err := ti.scanTakeout(rootDir, true)
	if err != nil {
		return nil, err
	}

	existing := ti.existingFileHashes()
	result := &domain.TakeoutScanResult{Warnings: warnings}
	albumSet := make(map[string]bool)
	hashCount := make(map[string]int)

	for _, item := range items {
		if item.isVideo {
			result.TotalVideos++
		} else {
			result.TotalPhotos++
		}
		result.TotalSizeBytes += item.size
		albumSet[item.album] = true

		h := computeFileHash(item.srcPath)
		if h == "" {
			continue
		}
		hashCount[h]++
		if existing[h] {
			result.Duplicates++
		}
	}
	// Count intra-export duplicates (same content appearing more than once).
	for _, n := range hashCount {
		if n > 1 {
			result.Duplicates += n - 1
		}
	}

	for album := range albumSet {
		result.Albums = append(result.Albums, album)
	}
	sort.Strings(result.Albums)

	return result, nil
}

// ImportTakeout runs an import. If req.Mode is "scan" it delegates to ScanTakeout
// and returns nil; callers should use the returned scan via a separate call. For
// mode "import" it copies media into the library and saves metadata in batches.
func (ti *TakeoutImporter) ImportTakeout(ctx context.Context, req domain.TakeoutImportRequest) error {
	if strings.EqualFold(req.Mode, "scan") {
		// Scan is handled by ScanTakeout; nothing to import.
		return nil
	}

	ti.mu.Lock()
	if ti.running {
		ti.mu.Unlock()
		return fmt.Errorf("an import is already in progress")
	}
	start := time.Now()
	ti.running = true
	ti.progress = domain.TakeoutImportProgress{Status: "running", Phase: "copying", StartedAt: start}
	ti.mu.Unlock()

	defer func() {
		ti.mu.Lock()
		ti.running = false
		ti.mu.Unlock()
	}()

	rootDir, cleanup, err := ti.openTakeout(req.Path)
	if err != nil {
		ti.finalize("error", 0, 0, 0, []string{err.Error()}, start)
		return err
	}
	defer cleanup()

	items, _, err := ti.scanTakeout(rootDir, req.PreserveAlbums)
	if err != nil {
		ti.finalize("error", 0, 0, 0, []string{err.Error()}, start)
		return err
	}

	var existing map[string]bool
	if req.Deduplicate {
		existing = ti.existingFileHashes()
	}
	seen := make(map[string]bool)

	errors := make([]string, 0)
	batch := make([]domain.PhotoFile, 0, importBatchSize)
	imported, skipped := 0, 0
	fatal := false

	report := func(phase string, current int) {
		ti.mu.Lock()
		ti.progress = domain.TakeoutImportProgress{
			Status: "running", Phase: phase, Current: current, Total: len(items),
			Imported: imported, Skipped: skipped, Errors: errors, StartedAt: start,
		}
		ti.mu.Unlock()
		if ti.wsHub != nil {
			// Takeout import targets the shared default workspace until
			// per-workspace import lands; scope the event so it reaches only
			// that workspace's clients (issue #74).
			ti.wsHub.BroadcastWorkspaceEvent(domain.DefaultWorkspaceID, ws.Event{
				Type: ws.EventImportProgress,
				Payload: ws.ImportProgressPayload{
					Status: "running", Phase: phase, Current: current, Total: len(items),
					Imported: imported, Skipped: skipped,
				},
			})
		}
	}

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		if ti.photoService != nil {
			if err := ti.photoService.SavePhotos(batch); err != nil {
				return err
			}
		}
		batch = batch[:0]
		return nil
	}

	for i, item := range items {
		if err := ctx.Err(); err != nil {
			fatal = true
			errors = append(errors, fmt.Sprintf("import cancelled: %v", err))
			break
		}
		current := i + 1

		if req.Deduplicate {
			h := computeFileHash(item.srcPath)
			if h != "" {
				if existing[h] || seen[h] {
					skipped++
					continue
				}
				seen[h] = true
			}
		}

		albumDir := filepath.Join(ti.config.PhotoDir, item.album)
		if err := os.MkdirAll(albumDir, 0o755); err != nil {
			errors = append(errors, fmt.Sprintf("create dir %s: %v", albumDir, err))
			continue
		}
		destPath := filepath.Join(albumDir, filepath.Base(item.srcPath))
		if err := copyFile(item.srcPath, destPath); err != nil {
			errors = append(errors, fmt.Sprintf("copy %s: %v", item.srcPath, err))
			continue
		}

		meta := ti.filesystemSvc.GetPhotoMetadata(destPath)
		meta.Directory = item.album
		if sc := item.sidecar; sc != nil {
			if t := sc.PhotoTakenTime; t != nil && t.Seconds > 0 {
				meta.ModifiedTime = t.Seconds
			}
			if g := sc.GeoData; g != nil {
				meta.Latitude = float64(g.LatitudeE7) / 1e7
				meta.Longitude = float64(g.LongitudeE7) / 1e7
			}
		}

		batch = append(batch, meta)
		imported++

		if len(batch) >= importBatchSize {
			if err := flush(); err != nil {
				fatal = true
				errors = append(errors, fmt.Sprintf("save batch: %v", err))
				break
			}
		}

		if current%progressBroadcastEvery == 0 || current == len(items) {
			report("copying", current)
		}
	}

	if !fatal {
		if err := flush(); err != nil {
			fatal = true
			errors = append(errors, fmt.Sprintf("save final batch: %v", err))
		}
	}

	status := "complete"
	if fatal {
		status = "error"
		slog.Error("takeout import failed", "status", status, "errors", errors)
	}
	ti.finalize(status, len(items), imported, skipped, errors, start)

	if fatal {
		return fmt.Errorf("import completed with errors: %s", strings.Join(errors, "; "))
	}
	return nil
}

// finalize records terminal progress and broadcasts the completion event.
func (ti *TakeoutImporter) finalize(status string, total, imported, skipped int, errors []string, start time.Time) {
	ti.mu.Lock()
	ti.progress = domain.TakeoutImportProgress{
		Status: status, Phase: "done", Current: total, Total: total,
		Imported: imported, Skipped: skipped, Errors: errors, StartedAt: start,
	}
	ti.mu.Unlock()

	if ti.wsHub != nil {
		// Scoped: the completion payload may carry error strings derived from
		// filesystem paths (issue #74).
		ti.wsHub.BroadcastWorkspaceEvent(domain.DefaultWorkspaceID, ws.Event{
			Type: ws.EventImportComplete,
			Payload: ws.ImportCompletePayload{
				Imported: imported, Skipped: skipped, Errors: errors,
			},
		})
	}
}

// openTakeout resolves the export path to a readable root directory. A directory
// is used as-is; a .zip is extracted to a temp dir (returned cleanup removes it).
func (ti *TakeoutImporter) openTakeout(path string) (string, func(), error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", nil, fmt.Errorf("takeout path not accessible: %w", err)
	}

	noCleanup := func() {}
	if info.IsDir() {
		return path, noCleanup, nil
	}
	if !strings.EqualFold(filepath.Ext(path), ".zip") {
		return "", nil, fmt.Errorf("takeout path must be a directory or .zip file: %s", path)
	}

	tmpDir, err := os.MkdirTemp("", "photobox-takeout-")
	if err != nil {
		return "", nil, fmt.Errorf("create temp dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }

	zr, err := zip.OpenReader(path)
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("open zip: %w", err)
	}
	defer zr.Close()

	for _, f := range zr.File {
		dest := filepath.Join(tmpDir, f.Name)
		// Guard against zip-slip path traversal.
		if dest != tmpDir && !strings.HasPrefix(dest, filepath.Clean(tmpDir)+string(os.PathSeparator)) {
			continue
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				cleanup()
				return "", nil, fmt.Errorf("create dir %s: %w", dest, err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("create dir %s: %w", filepath.Dir(dest), err)
		}
		src, err := f.Open()
		if err != nil {
			cleanup()
			return "", nil, fmt.Errorf("open zip entry %s: %w", f.Name, err)
		}
		out, err := os.Create(dest)
		if err != nil {
			src.Close()
			cleanup()
			return "", nil, fmt.Errorf("create %s: %w", dest, err)
		}
		_, copyErr := io.Copy(out, src)
		src.Close()
		out.Close()
		if copyErr != nil {
			cleanup()
			return "", nil, fmt.Errorf("extract %s: %w", f.Name, copyErr)
		}
	}

	return tmpDir, cleanup, nil
}

// scanTakeout walks rootDir and returns every recognised media file with its
// resolved album and (if present) parsed sidecar JSON. It also returns any
// warnings encountered during the walk.
func (ti *TakeoutImporter) scanTakeout(rootDir string, preserveAlbums bool) ([]takeoutItem, []string, error) {
	var items []takeoutItem
	warnings := make([]string, 0)
	unparsedSidecars := 0

	err := filepath.WalkDir(rootDir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if !takeoutMediaExtensions[ext] {
			return nil
		}

		item := takeoutItem{
			srcPath: p,
			size:    0,
			isVideo: takeoutVideoExtensions[ext],
		}
		if fi, statErr := d.Info(); statErr == nil {
			item.size = fi.Size()
		}

		sidecarPath := filepath.Join(filepath.Dir(p), strings.TrimSuffix(d.Name(), ext)+".json")
		if data, readErr := os.ReadFile(sidecarPath); readErr == nil {
			var sc domain.TakeoutSidecarJSON
			if jsonErr := json.Unmarshal(data, &sc); jsonErr != nil {
				unparsedSidecars++
			} else {
				item.sidecar = &sc
				if sc.Video {
					item.isVideo = true
				}
			}
		}

		dir := filepath.Dir(p)
		rel, relErr := filepath.Rel(rootDir, dir)
		album := defaultTakeoutAlbum
		if preserveAlbums && relErr == nil && rel != "." && !strings.HasPrefix(rel, "..") {
			album = filepath.Base(dir)
		}
		item.album = album

		items = append(items, item)
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("scan takeout: %w", err)
	}
	if unparsedSidecars > 0 {
		warnings = append(warnings, fmt.Sprintf("%d sidecar JSON file(s) could not be parsed and were ignored", unparsedSidecars))
	}
	return items, warnings, nil
}

// existingFileHashes returns the set of file hashes already present in the library.
func (ti *TakeoutImporter) existingFileHashes() map[string]bool {
	set := make(map[string]bool)
	cache, err := ti.photoRepo.GetPhotoIndexCache()
	if err != nil {
		return set
	}
	for _, entry := range cache {
		if entry.FileHash != "" {
			set[entry.FileHash] = true
		}
	}
	return set
}

// copyFile copies src to dest, fully closing both handles.
func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
