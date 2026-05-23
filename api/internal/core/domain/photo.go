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
}

type PhotoTag struct {
	PhotoID string `gorm:"primaryKey;index:idx_photo_tag,priority:1"`
	Tag     string `gorm:"primaryKey;index:idx_photo_tag,priority:2;index:idx_tag_photo,priority:1"`
	Source  string `gorm:"default:''"`
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
