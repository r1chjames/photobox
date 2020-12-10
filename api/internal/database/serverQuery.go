package database

import (
	b64 "encoding/base64"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gitlab.com/r1chjames/photobox/api/internal/utils"
	"gorm.io/datatypes"
	"log"
)

func SavePhotoRecordsToDatabase(photoRecords []PhotoFile) {
	for _, photo := range photoRecords {
		result, err := GetAlbumByName(photo.Directory)
		albumId := result.ID
		if checkNotFoundError(err) {
			albumId, _ = CreateAlbum(photo.Directory)
			if err != nil {
				log.Print("unable to insert album record")
			}
		}
		log.Printf("Adding photo: %s to album: %s", photo.Name, photo.Directory)
		photoHash := b64.StdEncoding.EncodeToString([]byte(photo.Path))
		photo := Photo{
			Name: utils.EscapeInvalidCharacters(photo.Name),
			FilesystemPath: utils.EscapeInvalidCharacters(photo.Path),
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