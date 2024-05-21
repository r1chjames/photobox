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
		dbEnv.SavePhotoRecordToDatabase(photo)
	}
}

func (dbEnv *Env) SavePhotoRecordToDatabase(photo PhotoFile) {
	result, err := dbEnv.GetAlbumByName(photo.Directory)
	albumId := result.ID
	if checkNotFoundError(err) {
		albumId, _ = dbEnv.CreateAlbum(photo.Directory)
		if err != nil {
			log.Printf("unable to insert album record, %s", err)
		}
	}
	log.Printf("Adding photo: %s to album: %s", photo.Name, photo.Directory)
	photoHash := b64.StdEncoding.EncodeToString([]byte(photo.Path))
	photoMetadata, _ := json.Marshal(&photo)
	photoInfo := Photo{
		Name:           utils.EscapeInvalidCharacters(photo.Name),
		FilesystemPath: utils.EscapeInvalidCharacters(photo.Path),
		AlbumId:        albumId,
		Tags:           "",
		Metadata:       photoMetadata,
		Thumbnail:      photo.Thumbnail,
	}
	photo.ID = photoHash
	err = dbEnv.CreatePhotoInfo(photoInfo)
	if err != nil {
		log.Printf("unable to insert photo record, %s", err)
	}
}
