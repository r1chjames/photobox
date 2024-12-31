package types

import (
	"gorm.io/datatypes"
	"time"
)

type Album struct {
	ID          string         `gorm:"primarykey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Tags        string         `json:"tags"`
	Metadata    datatypes.JSON `json:"metadata"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Thumnail    string         `json:"thumbnail"`
}

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

type Setting struct {
	Key          string    `gorm:"primarykey" json:"key"`
	Value        string    `json:"value"`
	FriendlyName string    `json:"friendlyName"`
	Category     string    `json:"category"`
	Type         string    `json:"type"`
	Options      string    `json:"options"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type Settings struct {
	Settings []Setting `binding:"required"`
}

type Job struct {
	Name      string    `gorm:"primarykey" json:"name"`
	Status    string    `json:"status"`
	LastRun   time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type PhotoUpload struct {
	Name          string `json:"name"`
	AlbumName     string `json:"albumName"`
	BinaryContent string `json:"binaryContent"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
