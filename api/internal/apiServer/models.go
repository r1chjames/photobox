package apiServer

import (
	"gorm.io/datatypes"
	"time"
)

type AlbumResponse struct {
	ID          string         `gorm:"primarykey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Tags        string         `json:"tags"`
	Metadata    datatypes.JSON `json:"metadata"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Thumnail    string         `json:"thumbnail"`
}

type PhotoResponse struct {
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

type SettingResponse struct {
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

type SettingsResponse struct {
	Settings []SettingResponse `binding:"required"`
}

type JobResponse struct {
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

type UserRole string

const (
	ADMINISTRATOR UserRole = "administrator"
	VIEWER        UserRole = "viewer"
	CONTRIBUTOR   UserRole = "contributor"
)

type UserResponse struct {
	ID       string   `gorm:"primarykey" json:"id"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
	Email    string   `json:"email"`
	Approved bool     `json:"-"`
}
