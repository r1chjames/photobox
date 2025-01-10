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
	AlbumId        string         `json:"albumId" gorm:"index"`
	Tags           string         `json:"tags"`
	Metadata       datatypes.JSON `json:"metadata"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Thumbnail      []byte         `json:"thumbnail"`
}

type PhotoUpload struct {
	Name          string `json:"name"`
	AlbumName     string `json:"albumName"`
	BinaryContent string `json:"binaryContent"`
}
