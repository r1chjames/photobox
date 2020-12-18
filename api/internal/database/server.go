package database

import (
	b64 "encoding/base64"
	"encoding/json"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gitlab.com/r1chjames/photobox/api/internal/utils"
	"log"
)

func (dbEnv *Env) SavePhotoRecordsToDatabase(photoRecords []PhotoFile) {
	for _, photo := range photoRecords {
		result, err := dbEnv.GetAlbumByName(photo.Directory)
		albumId := result.ID
		if checkNotFoundError(err) {
			albumId, _ = dbEnv.CreateAlbum(photo.Directory)
			if err != nil {
				log.Print("unable to insert album record")
			}
		}
		log.Printf("Adding photo: %s to album: %s", photo.Name, photo.Directory)
		photoHash := b64.StdEncoding.EncodeToString([]byte(photo.Path))
		photoMetadata, _ := json.Marshal(&photo)
		photo := Photo{
			Name:           utils.EscapeInvalidCharacters(photo.Name),
			FilesystemPath: utils.EscapeInvalidCharacters(photo.Path),
			AlbumId:        albumId,
			Tags:           "",
			Metadata:       photoMetadata,
		}
		photo.ID = photoHash
		err = dbEnv.CreatePhoto(photo)
		if err != nil {
			log.Print("unable to insert photo record")
		}
	}
}
