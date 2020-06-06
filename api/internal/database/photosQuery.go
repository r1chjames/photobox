package database

import (
	"database/sql"
	"errors"
	"fmt"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"log"
)

func transformRowsIntoPhotosArray(results *sql.Rows) interface{} {
	defer results.Close()
	var records []PhotoRecord
	for results.Next() {
		var record PhotoRecord
		err := results.Scan(
			&record.Id,
			&record.Name,
			&record.FilesystemPath,
			&record.AlbumId,
			&record.Tags,
			&record.Metadata)
		if err != nil {
			log.Fatal(fmt.Sprintf("An error parsing photo record: %s", err))
		}
		records = append(records, record)
	}
	return records
}

func GetPhotoInfoById(appConfig AppConfig, photoId string) (PhotoRecord, error) {
	allPhotoFields := "ID, NAME, FILESYSTEMPATH, ALBUMID, TAGS, METADATA"
	query := fmt.Sprintf("SELECT %s FROM photos WHERE id = '%s'", allPhotoFields, photoId)

	photoResults := executeTransformDbQuery(appConfig, query, transformRowsIntoPhotosArray).([]PhotoRecord)
	if len(photoResults) > 0 {
		return photoResults[0], nil
	}
	return PhotoRecord{}, errors.New("no photo found")
}

func GetAllPhotos(appConfig AppConfig) ([]PhotoRecord, error) {
	allPhotoFields := "ID, NAME, FILESYSTEMPATH, ALBUMID, TAGS, METADATA"
	query := fmt.Sprintf("SELECT %s FROM photos", allPhotoFields)

	photoResults := executeTransformDbQuery(appConfig, query, transformRowsIntoPhotosArray).([]PhotoRecord)
	if len(photoResults) > 0 {
		return photoResults, nil
	}
	return []PhotoRecord{}, errors.New("no photos found")
}

func GetAllPhotosInfoInAlbum(appConfig AppConfig, albumId string) ([]PhotoRecord, error) {
	allPhotoFields := "ID, NAME, FILESYSTEMPATH, ALBUMID, TAGS, METADATA"
	query := fmt.Sprintf("SELECT %s FROM photos WHERE albumId = '%s'", allPhotoFields, albumId)

	photoResults := executeTransformDbQuery(appConfig, query, transformRowsIntoPhotosArray).([]PhotoRecord)
	if len(photoResults) > 0 {
		return photoResults, nil
	}
	return []PhotoRecord{}, errors.New("no photos found")
}

func GetPhotosInAlbumCount(appConfig AppConfig, albumId string) int {
	query := fmt.Sprintf("SELECT COUNT(*) FROM photos WHERE albumId = '%s'", albumId)
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
