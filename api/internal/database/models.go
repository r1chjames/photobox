package database

import (
	"gorm.io/datatypes"
	"time"
)

type Album struct {
	ID          string `gorm:"primarykey"`
	Name        string
	Description string
	Tags        string
	Metadata    datatypes.JSON
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Thumnail    string
}

type Photo struct {
	ID             string `gorm:"primarykey"`
	Name           string
	FilesystemPath string
	SourcePath     string
	AlbumId        string `gorm:"index"`
	Tags           string
	Metadata       datatypes.JSON
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Thumbnail      []byte
}

type Setting struct {
	Key          string `gorm:"primarykey"`
	Value        string
	FriendlyName string
	Category     string
	Type         string
	Options      string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Settings struct {
	Settings []Setting `binding:"required"`
}

type Job struct {
	Name      string `gorm:"primarykey"`
	Status    string
	LastRun   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRole string

const (
	ADMINISTRATOR UserRole = "administrator"
	VIEWER        UserRole = "viewer"
	CONTRIBUTOR   UserRole = "contributor"
)

type User struct {
	ID       string `gorm:"primarykey"`
	Username string `gorm:"index"`
	Role     UserRole
	Email    string
	Approved bool
	Password string
}
