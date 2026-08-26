package domain

import "time"

// TakeoutSidecarJSON matches the Google Takeout sidecar JSON schema.
// Each photo/video file has a corresponding .json file with metadata.
type TakeoutSidecarJSON struct {
	AlbumKey        string   `json:"albumKey,omitempty"`
	Description     string   `json:"description,omitempty"`
	GeoData         *GeoData `json:"geoData,omitempty"`
	LivePhotoURL    string   `json:"livePhotoUrl,omitempty"`
	Panoramic       bool     `json:"panoramic,omitempty"`
	People          []string `json:"people,omitempty"`
	PhotoTakenTime  *TimeRef `json:"photoTakenTime,omitempty"`
	Screenshot      bool     `json:"screenshot,omitempty"`
	Starred         bool     `json:"starred,omitempty"`
	Video           bool     `json:"video,omitempty"`
}

// GeoData holds latitude/longitude as E7 (integer × 1e-7 degrees).
type GeoData struct {
	LatitudeE7   int64 `json:"latitudeE7"`
	LongitudeE7  int64 `json:"longitudeE7"`
	Timezone     string `json:"timezone,omitempty"`
	LocationName string `json:"locationName,omitempty"`
}

// TimeRef is a Google protobuf Timestamp serialized as JSON.
type TimeRef struct {
	Seconds int64 `json:"seconds"`
	Nanos   int32 `json:"nanos,omitempty"`
}

// TakeoutImportRequest is the POST /api/import/takeout request body.
type TakeoutImportRequest struct {
	Path            string `json:"path" binding:"required"`
	Mode            string `json:"mode" binding:"required,oneof=scan import"`
	PreserveAlbums  bool   `json:"preserve_albums"`
	Deduplicate     bool   `json:"deduplicate"`
}

// TakeoutScanResult is the response for mode=scan.
type TakeoutScanResult struct {
	TotalPhotos    int      `json:"total_photos"`
	TotalVideos    int      `json:"total_videos"`
	TotalSizeBytes int64    `json:"total_size_bytes"`
	Duplicates     int      `json:"duplicates"`
	Albums         []string `json:"albums"`
	Warnings       []string `json:"warnings,omitempty"`
}

// TakeoutImportProgress tracks the in-flight import state.
type TakeoutImportProgress struct {
	Status    string   `json:"status"` // "running" | "complete" | "error"
	Phase     string   `json:"phase"`  // "copying" | "indexing" | "done"
	Current   int      `json:"current"`
	Total     int      `json:"total"`
	Imported  int      `json:"imported"`
	Skipped   int      `json:"skipped"`
	Errors    []string `json:"errors,omitempty"`
	StartedAt time.Time `json:"started_at"`
}
