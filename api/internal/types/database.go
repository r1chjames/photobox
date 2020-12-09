package types

import (
	"github.com/rwcarlsen/goexif/exif"
	"time"
)

type PhotoOld struct {
	Id string `json:"id"`
	Name string `json:"name"`
	FilesystemPath string `json:"filesystemPath"`
	AlbumId string `json:"albumId"`
	Tags string `json:"tags"`
	Metadata string `json:"metadata"`
}

type photoMetadata struct {
	Exif exif.Exif `json:"exif,omitempty"`
}

type AlbumOld struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Tags string `json:"tags"`
	Metadata string `json:"metadata"`
}

type albumMetadata struct {
	Created	time.Duration `json:"created"`
	LastUpdated time.Duration `json:"lastUpdated"`
}

type ServerTask struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Interval string `json:"interval"`
	LastRun string `json:"lastRun"`
}

type SettingOld struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type JobOld struct {
	Job	string `json:"job"`
	Running bool `json:"running"`
	LastRun string `json:"lastRun"`
}