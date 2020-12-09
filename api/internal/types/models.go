package types

import (
	"gorm.io/datatypes"
	"time"
)

type common struct {
	ID        string `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Album struct {
	common
	Name string `json:"name"`
	Description string `json:"description"`
	Tags string `json:"tags"`
	Metadata datatypes.JSON `json:"metadata"`
}

type Photo struct {
	common
	Name string `json:"name"`
	FilesystemPath string `json:"filesystemPath"`
	AlbumId string `json:"albumId"`
	Tags string `json:"tags"`
	Metadata datatypes.JSON `json:"metadata"`
}

type Setting struct {
	common
	Key string `json:"key"`
	Value string `json:"value"`
}

type Job struct {
	common
	Name string `json:"name"`
	Status string `json:"status"`
	LastRun time.Time `json:"lastRun"`
}