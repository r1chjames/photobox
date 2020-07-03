package types

import (
	"github.com/rwcarlsen/goexif/exif"
	"time"
)

type PhotoRecord struct {
	Id string `json:"id"`
	Name string `json:"name"`
	FilesystemPath string `json:"filesystemPath"`
	AlbumId string `json:"albumId"`
	Tags string `json:"tags"`
	Metadata string `json:"metadata"`
}

type photoMetadataRecord struct {
	Exif exif.Exif `json:"exif,omitempty"`
}

type AlbumRecord struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Tags string `json:"tags"`
	Metadata string `json:"metadata"`
}

type albumMetadataRecord struct {
	Created	time.Duration `json:"created"`
	LastUpdated time.Duration `json:"lastUpdated"`
}

type ServerTask struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Interval string `json:"interval"`
	LastRun string `json:"lastRun"`
}

type SettingRecord struct {
	Setting string `json:"setting"`
	Value string `json:"value"`
}

type JobsRecord struct {
	Job	string `json:"job"`
	Running bool `json:"running"`
	LastRun string `json:"lastRun"`
}