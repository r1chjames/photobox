package database

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/google/uuid"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
	"strings"
)

func SavePhotoRecordsToDatabase(appConfig AppConfig, photoRecords []PhotoFile) {
	defer jobCompleted(appConfig, "Photo_index")
	for _, photo := range photoRecords {
		albumId := GetAlbumByName(appConfig, photo.Directory).Id
		if albumId == "" {
			createAlbum(appConfig, photo.Directory)
		}
		log.Printf("Adding photo: %s to album: %s", photo.Name, photo.Directory)
		photoHash := b64.StdEncoding.EncodeToString([]byte(photo.Path))
		photoInsertQuery := fmt.Sprintf("INSERT IGNORE INTO photos VALUES ('%s','%s','%s','%s','%s','%s')",
			photoHash,
			escapeValue(photo.Name),
			escapeValue(photo.Path),
			albumId,
			"",
			"{}")
		executeDbInsert(appConfig, photoInsertQuery)
	}
}

func escapeValue(field string) string {
	return strings.ReplaceAll(field, "'", "\\'")
}

func createAlbum(appConfig AppConfig, name string) {
	log.Printf("Creating new album: %s", name)
	photoInsertQuery := fmt.Sprintf("INSERT IGNORE INTO albums VALUES ('%s','%s','%s','%s','%s')",
		uuid.New(),
		escapeValue(name),
		"",
		"",
		"{}")
	executeDbInsert(appConfig, photoInsertQuery)
}