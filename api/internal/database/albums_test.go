package database

import (
	"fmt"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)


func TestShouldGetAlbumByID(t *testing.T) {
	env, mock, dbConn := MockDB(t)
	defer dbConn.Close()

	ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `albums` WHERE `albums`.`id` = (.+) ORDER BY `albums`.`id` LIMIT 1", sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("albumId1", "albumName1", "albumDesc1"))

	if _, err := env.GetAlbumById("albumId1"); err != nil {
		t.Errorf("error was not expected while getting albums: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetAlbumByName(t *testing.T) {
	env, mock, dbConn := MockDB(t)
	defer dbConn.Close()

	ShouldReturnRowsForQuery(mock, "SELECT (.+) FROM `albums` WHERE name = (.+) ORDER BY `albums`.`id` LIMIT 1", sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("albumId1", "albumName1", "albumDesc1"))

	if _, err := env.GetAlbumByName("albumName1"); err != nil {
		t.Errorf("error was not expected while getting albums: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetAllAlbumsWithNoOffset(t *testing.T) {
	env, mock, dbConn := MockDB(t)
	defer dbConn.Close()

	page := 1 //page 1 means offset = 0 which will be omitted from query
	limit := 5
	ShouldReturnRowsForQuery(mock, fmt.Sprintf("SELECT (.+) FROM `albums` LIMIT %d", limit), sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(
		"albumId1", "albumName1", "albumDesc1").AddRow(
		"albumId2", "albumName2", "albumDesc2").AddRow(
		"albumId3", "albumName3", "albumDesc3").AddRow(
		"albumId4", "albumName4", "albumDesc4").AddRow(
		"albumId5", "albumName5", "albumDesc5"))

	albums, err := env.GetAllAlbums(page, limit)
	if err != nil {
		t.Errorf("error was not expected while getting all albums: %s", err)
	}


	expectedAlbumCount := 5
	if len(albums) != expectedAlbumCount {
		t.Errorf("incorrect number of albums returned, expected: %d, got: %d", expectedAlbumCount, len(albums))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetAllAlbumsWithWithOffset(t *testing.T) {
	env, mock, dbConn := MockDB(t)
	defer dbConn.Close()

	page := 2
	offset := 2 //page 2 means offset = 2 which will be omitted from query
	limit := 2

	ShouldReturnRowsForQuery(mock, fmt.Sprintf("SELECT (.+) FROM `albums` LIMIT %d OFFSET %d", limit, offset), sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(
		"albumId1", "albumName1", "albumDesc1").AddRow(
		"albumId2", "albumName2", "albumDesc2"))

	albums, err := env.GetAllAlbums(page, limit)
	if err != nil {
		t.Errorf("error was not expected while getting all albums: %s", err)
	}


	expectedAlbumCount := 2
	if len(albums) != expectedAlbumCount {
		t.Errorf("incorrect number of albums returned, expected: %d, got: %d", expectedAlbumCount, len(albums))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}