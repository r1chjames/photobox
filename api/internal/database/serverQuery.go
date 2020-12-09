package database

import (
	b64 "encoding/base64"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/datatypes"
	"log"
	"strings"
)

func SavePhotoRecordsToDatabase(photoRecords []PhotoFile) {
	for _, photo := range photoRecords {
		result, err := GetAlbumByName(photo.Directory)
		albumId := result.ID
		if checkNotFoundError(err) {
			CreateAlbum(photo.Directory)
		}
		log.Printf("Adding photo: %s to album: %s", photo.Name, photo.Directory)
		photoHash := b64.StdEncoding.EncodeToString([]byte(photo.Path))
		photo := Photo{
			Name: escapeValue(photo.Name),
			FilesystemPath: escapeValue(photo.Path),
			AlbumId: albumId,
			Tags: "",
			Metadata: datatypes.JSON(`{}`),
		}
		photo.ID = photoHash
		err = CreatePhoto(photo)
		if err != nil {
			log.Print("unable to insert photo record")
		}
	}
}

func escapeValue(field string) string {
	return strings.ReplaceAll(field, "'", "\\'")
}
