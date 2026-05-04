package domain

import (
	"github.com/rwcarlsen/goexif/exif"
)

type PhotoFile struct {
	ID        string `gorm:"primarykey" json:"-"`
	MD5       string `json:"md5"`
	Path      string `json:"path"`
	Directory string `json:"directory"`
	//DateCreated os.time
	//DateModified os.time
	Size      int64     `json:"size"`
	Extension string    `json:"extension"`
	Name      string    `json:"name"`
	Exif      exif.Exif `json:"exif"`
	Mime      string    `json:"mime"`
	Thumbnail []byte    `json:"-" gorm:"-"`
	MediaType string    `json:"mediaType"`
	Duration  int       `json:"duration"`
	Width     int   `json:"width"`
	Height    int   `json:"height"`
	ModifiedTime int64 `json:"modifiedTime" gorm:"-"`
	Latitude     float64 `json:"latitude" gorm:"-"`
	Longitude    float64 `json:"longitude" gorm:"-"`
}
