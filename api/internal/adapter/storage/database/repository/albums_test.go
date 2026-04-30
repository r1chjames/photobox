package repository

import (
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestShouldGetAlbumByID(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "albums" WHERE "albums"\."id" = \$1 ORDER BY "albums"\."id" LIMIT \$2`, sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("albumId1", "albumName1", "albumDesc1"))

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

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "albums" WHERE name = \$1 ORDER BY "albums"\."id" LIMIT \$2`, sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("albumId1", "albumName1", "albumDesc1"))

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

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "albums" ORDER BY created_epoch ASC LIMIT \$1`, sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(
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

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "albums" WHERE created_epoch > \$1 ORDER BY created_epoch ASC LIMIT \$2`, sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(
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

// TestUpdateAlbum_Success tests successful album update
func TestUpdateAlbum_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	album := &domain.Album{
		ID:          "album123",
		Name:        "Updated Album",
		Description: "Updated Description",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "albums" SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewAlbumRepository(env)
	err := repo.UpdateAlbum(album)

	if err != nil {
		t.Errorf("error was not expected while updating album: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestDeleteAlbum_Success tests successful album deletion
func TestDeleteAlbum_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	albumID := "album123"

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "albums" WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewAlbumRepository(env)
	err := repo.DeleteAlbum(albumID)

	if err != nil {
		t.Errorf("error was not expected while deleting album: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestReassignPhotosToAlbum_Success tests successful photo reassignment
func TestReassignPhotosToAlbum_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	fromAlbumID := "album1"
	toAlbumID := "album2"

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "photos" SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()

	repo := NewAlbumRepository(env)
	err := repo.ReassignPhotosToAlbum(fromAlbumID, toAlbumID)

	if err != nil {
		t.Errorf("error was not expected while reassigning photos: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestSearchAlbums_Success tests successful album search
func TestSearchAlbums_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	query := "vacation"
	expectedAlbums := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow("album1", "Vacation 2024", "Summer trip").
		AddRow("album2", "Vacation 2023", "Winter trip")

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "albums" WHERE to_tsvector\('english', coalesce\(name, ''\)\) @@ plainto_tsquery\('english', \$1\) ORDER BY created_epoch DESC LIMIT \$2`, expectedAlbums)

	repo := NewAlbumRepository(env)
	albums, err := repo.SearchAlbums(query, 10)

	if err != nil {
		t.Errorf("error was not expected while searching albums: %s", err)
	}

	expectedCount := 2
	if len(albums) != expectedCount {
		t.Errorf("expected %d albums, got %d", expectedCount, len(albums))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestSearchAlbums_Empty tests search with no results
func TestSearchAlbums_Empty(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	query := "nonexistent"
	expectedAlbums := sqlmock.NewRows([]string{"id", "name", "description"})

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "albums" WHERE to_tsvector\('english', coalesce\(name, ''\)\) @@ plainto_tsquery\('english', \$1\) ORDER BY created_epoch DESC LIMIT \$2`, expectedAlbums)

	repo := NewAlbumRepository(env)
	albums, err := repo.SearchAlbums(query, 10)

	if err == nil {
		t.Error("expected ErrDataNotFound when search returns no results, got nil error")
	}

	if albums != nil {
		t.Error("expected nil albums when search is empty, got non-nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
