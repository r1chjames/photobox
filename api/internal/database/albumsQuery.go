package database

import (
	"database/sql"
	"fmt"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
)

func transformRowsIntoAlbumsArray(results *sql.Rows) interface{} {
	defer results.Close()
	var records []AlbumRecord
	for results.Next() {
		var record AlbumRecord
		err := results.Scan(
			&record.Id,
			&record.Name,
			&record.Description,
			&record.Tags,
			&record.Metadata)
		if err != nil {
			log.Fatal(fmt.Sprintf("An error parsing Album record: %s", err))
		}
		records = append(records, record)
	}
	return records
}

func GetAlbumById(appConfig AppConfig, albumId string) AlbumRecord {
	allAlbumFields := "ID, NAME, DESCRIPTION, TAGS, METADATA"
	query := fmt.Sprintf("SELECT %s FROM albums WHERE id = '%s'", allAlbumFields, albumId)

	albumResults := executeTransformDbQuery(appConfig, query, transformRowsIntoAlbumsArray).([]AlbumRecord)
	if len(albumResults) > 0 {
		return albumResults[0]
	}
	return AlbumRecord{}
}

func GetAlbumByName(appConfig AppConfig, albumName string) AlbumRecord {
	allAlbumFields := "ID, NAME, DESCRIPTION, TAGS, METADATA"
	query := fmt.Sprintf("SELECT %s FROM albums WHERE name = '%s'", allAlbumFields, albumName)

	albumResults := executeTransformDbQuery(appConfig, query, transformRowsIntoAlbumsArray).([]AlbumRecord)
	if len(albumResults) > 0 {
		return albumResults[0]
	}
	return AlbumRecord{}
}

func GetAllAlbums(appConfig AppConfig) []AlbumRecord {
	allAlbumFields := "ID, NAME, DESCRIPTION, TAGS, METADATA"
	query := fmt.Sprintf("SELECT %s FROM albums", allAlbumFields)
	albumResults := executeTransformDbQuery(appConfig, query, transformRowsIntoAlbumsArray).([]AlbumRecord)
	if len(albumResults) > 0 {
		return albumResults
	}
	return []AlbumRecord{}
}

func GetAlbumCount(appConfig AppConfig) int {
	query := "SELECT COUNT(*) FROM albums"
	transform := func(results *sql.Rows) interface{} {
		defer results.Close()
		var count int

		for results.Next() {
			if err := results.Scan(&count); err != nil {
				log.Fatal(err)
			}
		}
		return count
	}
	count := executeTransformDbQuery(appConfig, query, transform).(int)
	return count
}
