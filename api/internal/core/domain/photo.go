package domain

import (
	"gorm.io/datatypes"
	"time"
)

type Photo struct {
	ID             string         `gorm:"primarykey" json:"id"`
	Name           string         `json:"name"`
	FilesystemPath string         `json:"filesystemPath"`
	SourcePath     string         `json:"sourcePath"`
	AlbumId        string         `json:"albumId" gorm:"index:idx_album_deleted_epoch,priority:1"`
	Metadata       datatypes.JSON `json:"metadata"`
	CreatedAt      time.Time      `json:"createdAt"`
	CreatedEpoch   int64          `json:"createdEpoch" gorm:"index;index:idx_album_deleted_epoch,priority:3;index:idx_deleted_epoch,priority:2"`
	Year           int            `json:"year" gorm:"index"`
	Month          int            `json:"month" gorm:"index"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Thumbnail      []byte         `json:"-"`
	ThumbnailPath  string         `json:"-"`
	ThumbnailUrl   string         `json:"thumbnailUrl" gorm:"-"`
	Favorite       bool           `json:"favorite" gorm:"default:false;index"`
	DeletedAt      *time.Time     `json:"deletedAt" gorm:"index;index:idx_album_deleted_epoch,priority:2;index:idx_deleted_epoch,priority:1"`
	Blurhash       string         `json:"blurhash"`
	DominantColor  string         `json:"dominantColor"`
	FileHash       string `json:"fileHash" gorm:"index"`
	FileModifiedTime int64 `json:"fileModifiedTime" gorm:"index"`
	MediaType      string `json:"mediaType" gorm:"default:'image'"`
	Duration       int            `json:"duration"`
	Width          int            `json:"width"`
	Height         int            `json:"height"`
	Latitude       float64        `json:"latitude" gorm:"index"`
	Longitude      float64        `json:"longitude" gorm:"index"`
	Hidden         bool           `json:"hidden" gorm:"default:false;index"`
	Description    string         `json:"description" gorm:"type:text"`
	QualityScore   int            `json:"qualityScore" gorm:"default:0"`
	BlurScore      float64        `json:"blurScore" gorm:"default:0"`
	IsLowQuality   bool           `json:"isLowQuality" gorm:"default:false;index"`
	// TrashPath records where the original file was moved when the photo was
	// soft-deleted, so RestorePhoto can move it back. Empty for photos that
	// were never deleted (or have no filesystem file).
	TrashPath string `json:"-" gorm:"type:text"`
	// EditParams holds the last-applied non-destructive edit parameters
	// (rotate/crop/adjustments) as JSON. Empty when the photo is unedited.
	EditParams datatypes.JSON `json:"editParams" gorm:"type:jsonb"`
}

type PhotoTag struct {
	PhotoID string `gorm:"primaryKey;index:idx_photo_tag,priority:1"`
	Tag     string `gorm:"primaryKey;index:idx_photo_tag,priority:2;index:idx_tag_photo,priority:1"`
	Source  string `gorm:"default:''"`
}

// EditParams describes a non-destructive edit. All fields are optional;
// only provided fields are applied. Applied to a copy, never the original.
type EditParams struct {
	// Rotate is clockwise rotation in degrees: 0, 90, 180 or 270.
	Rotate *int `json:"rotate,omitempty"`
	// Crop is an optional region in fractional coords (0-1) of the image.
	Crop *CropParams `json:"crop,omitempty"`
	// Brightness / Contrast / Saturation are percentage adjustments
	// (imaging semantics: -100..100).
	Brightness *float64 `json:"brightness,omitempty"`
	Contrast   *float64 `json:"contrast,omitempty"`
	Saturation *float64 `json:"saturation,omitempty"`
	// AutoEnhance applies a gamma correction tuned to brighten shadows
	// (simple auto-levels without a full histogram stretch).
	AutoEnhance *bool `json:"autoEnhance,omitempty"`
}

// CropParams is a fractional crop rectangle (0-1 across each axis).
type CropParams struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// PhotoSearchFilters holds advanced multi-filter search criteria. Zero-value
// fields are ignored, so filters combine freely.
type PhotoSearchFilters struct {
	StartDate string
	EndDate   string
	MediaType string
	// Camera matches EXIF Make/Model (e.g. "Canon" or "iPhone").
	Camera string
	// HasGPS restricts to photos with coordinates when true.
	HasGPS bool
	// Orientation: "landscape", "portrait" or "square".
	Orientation string
	// Tags requires all listed tags (AND).
	Tags []string
	// Favorite restricts to favourited photos when true.
	Favorite bool
	// LowQuality restricts to quality-flagged photos when true.
	LowQuality bool
}

type TimelineEntry struct {
	Year  int   `json:"year"`
	Month int   `json:"month"`
	Count int64 `json:"count"`
}

type PhotoGeoData struct {
	ID          string `json:"id"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Thumbnail   string `json:"thumbnail"`
	DateTaken   string `json:"dateTaken"`
}

type PhotoUpload struct {
	Name          string `json:"name"`
	AlbumName     string `json:"albumName"`
	BinaryContent string `json:"binaryContent"`
}

type PhotoAnalysis struct {
	PhotoID       string         `gorm:"primaryKey" json:"photoId"`
	Model         string         `json:"model"`
	Caption       string         `json:"caption"`
	Tags          datatypes.JSON `json:"tags"`
	Objects       datatypes.JSON `json:"objects"`
	IsNSFW        bool           `json:"isNsfw"`
	IsPortrait    bool           `json:"isPortrait"`
	Status        string         `json:"status" gorm:"default:'pending'"`
	Attempts      int            `json:"attempts" gorm:"default:0"`
	ErrorMessage  string         `json:"errorMessage"`
	Retryable     bool           `json:"retryable" gorm:"default:true"`
	StartedAt     *time.Time     `json:"startedAt"`
	LastAttemptAt *time.Time     `json:"lastAttemptAt"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

func (PhotoAnalysis) TableName() string {
	return "photo_analysis"
}
