package types

import (
	"github.com/rwcarlsen/goexif/exif"
)

type PhotoFile struct {
	ID        string `gorm:"primarykey" json:"id"`
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
	Thumbnail []byte    `json:"thumbnail"`
}
