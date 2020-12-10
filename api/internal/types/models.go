package types

import (
	"gorm.io/datatypes"
	"time"
)

type Album struct {
	ID string `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Tags string `json:"tags"`
	Metadata datatypes.JSON `json:"metadata"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Photo struct {
	ID string `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
	FilesystemPath string `json:"filesystemPath"`
	AlbumId string `json:"albumId"`
	Tags string `json:"tags"`
	Metadata datatypes.JSON `json:"metadata"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Setting struct {
	Key string `gorm:"primarykey" json:"key"`
	Value string `json:"value"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Job struct {
	Name string `gorm:"primarykey" json:"name"`
	Status string `json:"status"`
	LastRun time.Time `json:"lastRun"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}