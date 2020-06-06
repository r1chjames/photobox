package types

import "github.com/rwcarlsen/goexif/exif"

type PhotoFile struct {
	MD5 string
	Path string
	Directory string
	//DateCreated os.time
	//DateModified os.time
	Size int64
	Extension string
	Name string
	Exif exif.Exif
	Mime string
}
