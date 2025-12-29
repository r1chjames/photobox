package repository

import (
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestShouldGetAlbumByID(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "albums" WHERE "albums"\."id" = \$1 ORDER BY "albums"\."id" LIMIT \$2`, sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("albumId1", "albumName1", "albumDesc1"))

	repo := NewAlbumRepository(env)
	if _, err := repo.GetAlbumById("albumId1"); err != nil {
		t.Errorf("error was not expected while getting albums: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetAlbumByName(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "albums" WHERE name = \$1 ORDER BY "albums"\."id" LIMIT \$2`, sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("albumId1", "albumName1", "albumDesc1"))

	repo := NewAlbumRepository(env)
	if _, err := repo.GetAlbumByName("albumName1"); err != nil {
		t.Errorf("error was not expected while getting albums: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldListAllAlbumsWithNoOffset(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "albums" LIMIT \$1`, sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(
		"albumId1", "albumName1", "albumDesc1").AddRow(
		"albumId2", "albumName2", "albumDesc2").AddRow(
		"albumId3", "albumName3", "albumDesc3").AddRow(
		"albumId4", "albumName4", "albumDesc4").AddRow(
		"albumId5", "albumName5", "albumDesc5"))

	repo := NewAlbumRepository(env)
	albums, err := repo.ListAllAlbums("", 5)
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

func TestShouldListAllAlbumsWithFromId(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	fromId := "dGVzdA==" // base64 encoded

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "albums" WHERE created_epoch > \$1 LIMIT \$2`, sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(
		"albumId1", "albumName1", "albumDesc1").AddRow(
		"albumId2", "albumName2", "albumDesc2"))

	repo := NewAlbumRepository(env)
	albums, err := repo.ListAllAlbums(fromId, 2)
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
