package repository

import (
	"testing"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// TestGetPhotoById_Success tests successful photo retrieval by ID
func TestGetPhotoById_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photoID := "photo123"
	expectedPhoto := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "tags", "metadata", "thumbnail", "created_epoch"}).
		AddRow(photoID, "test.jpg", "/path/to/test.jpg", "album1", "", []byte("{}"), []byte("thumbnail-data"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`, expectedPhoto)

	repo := NewPhotoRepository(env)
	photo, err := repo.GetPhotoById(photoID, true)

	if err != nil {
		t.Errorf("error was not expected while getting photo: %s", err)
	}

	if photo == nil {
		t.Error("expected photo to be returned, got nil")
	}

	if photo.ID != photoID {
		t.Errorf("expected photo ID %s, got %s", photoID, photo.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetPhotoById_WithoutThumbnail tests photo retrieval without thumbnail
func TestGetPhotoById_WithoutThumbnail(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photoID := "photo456"
	expectedPhoto := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "tags", "metadata", "created_at", "created_epoch", "updated_at"}).
		AddRow(photoID, "test2.jpg", "/path/to/test2.jpg", "", "album2", "", []byte("{}"), time.Now(), time.Now().UnixMilli(), time.Now())

	// When includeThumbnail is false, GORM selects specific columns (omits thumbnail)
	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`, expectedPhoto)

	repo := NewPhotoRepository(env)
	photo, err := repo.GetPhotoById(photoID, false)

	if err != nil {
		t.Errorf("error was not expected while getting photo: %s", err)
	}

	if photo == nil {
		t.Error("expected photo to be returned, got nil")
	}

	if photo.ID != photoID {
		t.Errorf("expected photo ID %s, got %s", photoID, photo.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetPhotoById_NotFound tests photo not found scenario
func TestGetPhotoById_NotFound(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photoID := "nonexistent"
	database.ShouldReturnNotFoundErrorForQuery(mock, `SELECT \* FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`)

	repo := NewPhotoRepository(env)
	photo, err := repo.GetPhotoById(photoID, true)

	if err == nil {
		t.Error("expected error when photo not found, got nil")
	}

	if photo != nil {
		t.Error("expected nil photo when not found, got non-nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListAllPhotos_Success tests successful listing of photos
func TestListAllPhotos_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "tags", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo1", "test1.jpg", "/path/to/test1.jpg", "album1", "", []byte("{}"), []byte("thumb1"), time.Now().UnixMilli()).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", "album1", "", []byte("{}"), []byte("thumb2"), time.Now().UnixMilli()).
		AddRow("photo3", "test3.jpg", "/path/to/test3.jpg", "album2", "", []byte("{}"), []byte("thumb3"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "photos" WHERE deleted_at IS NULL LIMIT \$1`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotos("", 10, true)

	if err != nil {
		t.Errorf("error was not expected while listing photos: %s", err)
	}

	expectedCount := 3
	if len(photos) != expectedCount {
		t.Errorf("expected %d photos, got %d", expectedCount, len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListAllPhotos_WithFromId tests listing photos with pagination from a specific ID
func TestListAllPhotos_WithFromId(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	fromID := "photo1"
	fromEpoch := time.Now().UnixMilli()

	// First query to get the fromPhoto - uses Omit("thumbnail") so specific columns are selected
	fromPhotoRow := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "tags", "metadata", "created_at", "created_epoch", "updated_at"}).
		AddRow(fromID, "test1.jpg", "/path/to/test1.jpg", "", "album1", "", []byte("{}"), time.Now(), fromEpoch, time.Now())
	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`, fromPhotoRow)

	// Second query to get photos after fromEpoch - also uses Omit("thumbnail")
	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "tags", "metadata", "created_at", "created_epoch", "updated_at"}).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", "", "album1", "", []byte("{}"), time.Now(), fromEpoch+1000, time.Now()).
		AddRow("photo3", "test3.jpg", "/path/to/test3.jpg", "", "album2", "", []byte("{}"), time.Now(), fromEpoch+2000, time.Now())

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NULL AND created_epoch > \$1 LIMIT \$2`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotos(fromID, 5, false)

	if err != nil {
		t.Errorf("error was not expected while listing photos: %s", err)
	}

	expectedCount := 2
	if len(photos) != expectedCount {
		t.Errorf("expected %d photos, got %d", expectedCount, len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListAllPhotos_Empty tests listing when no photos exist
func TestListAllPhotos_Empty(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "tags", "metadata", "thumbnail", "created_epoch"})

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "photos" WHERE deleted_at IS NULL LIMIT \$1`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotos("", 10, true)

	if err == nil {
		t.Error("expected ErrDataNotFound when no photos exist, got nil error")
	}

	if photos != nil {
		t.Error("expected nil photos when list is empty, got non-nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListAllPhotosInAlbum_Success tests successful listing of photos in an album
func TestListAllPhotosInAlbum_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	albumID := "album123"
	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "tags", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo1", "test1.jpg", "/path/to/test1.jpg", albumID, "", []byte("{}"), []byte("thumb1"), time.Now().UnixMilli()).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", albumID, "", []byte("{}"), []byte("thumb2"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "photos" WHERE deleted_at IS NULL AND album_id = \$1 LIMIT \$2`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotosInAlbum(albumID, "", 10, true)

	if err != nil {
		t.Errorf("error was not expected while listing photos in album: %s", err)
	}

	expectedCount := 2
	if len(photos) != expectedCount {
		t.Errorf("expected %d photos, got %d", expectedCount, len(photos))
	}

	// Verify all photos belong to the correct album
	for _, photo := range photos {
		if photo.AlbumId != albumID {
			t.Errorf("expected photo to belong to album %s, got %s", albumID, photo.AlbumId)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListAllPhotosInAlbum_WithFromId tests pagination in album listing
func TestListAllPhotosInAlbum_WithFromId(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	albumID := "album123"
	fromID := "cGhvdG8x" // base64 encoded epoch

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "tags", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", albumID, "", []byte("{}"), []byte("thumb2"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "photos" WHERE deleted_at IS NULL AND created_epoch > \$1 AND album_id = \$2 LIMIT \$3`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotosInAlbum(albumID, fromID, 5, true)

	if err != nil {
		t.Errorf("error was not expected while listing photos in album: %s", err)
	}

	if len(photos) != 1 {
		t.Errorf("expected 1 photo, got %d", len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListAllPhotosInAlbum_Empty tests listing when album has no photos
func TestListAllPhotosInAlbum_Empty(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	albumID := "emptyAlbum"
	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "tags", "metadata", "thumbnail", "created_epoch"})

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "photos" WHERE deleted_at IS NULL AND album_id = \$1 LIMIT \$2`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotosInAlbum(albumID, "", 10, true)

	if err == nil {
		t.Error("expected ErrDataNotFound when album is empty, got nil error")
	}

	if photos != nil {
		t.Error("expected nil photos when album is empty, got non-nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetPhotosInAlbumCount_Success tests successful photo count retrieval
func TestGetPhotosInAlbumCount_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	albumID := "album123"
	expectedCount := int64(42)

	countRow := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)
	database.ShouldReturnRowsForQuery(mock, `SELECT count\(\*\) FROM "photos" WHERE album_id = \$1 AND deleted_at IS NULL`, countRow)

	repo := NewPhotoRepository(env)
	count, err := repo.GetPhotosInAlbumCount(albumID)

	if err != nil {
		t.Errorf("error was not expected while counting photos: %s", err)
	}

	if count != expectedCount {
		t.Errorf("expected count %d, got %d", expectedCount, count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetPhotosInAlbumCount_Zero tests count when album is empty
func TestGetPhotosInAlbumCount_Zero(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	albumID := "emptyAlbum"
	expectedCount := int64(0)

	countRow := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)
	database.ShouldReturnRowsForQuery(mock, `SELECT count\(\*\) FROM "photos" WHERE album_id = \$1 AND deleted_at IS NULL`, countRow)

	repo := NewPhotoRepository(env)
	count, err := repo.GetPhotosInAlbumCount(albumID)

	if err != nil {
		t.Errorf("error was not expected while counting photos: %s", err)
	}

	if count != expectedCount {
		t.Errorf("expected count %d, got %d", expectedCount, count)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreatePhotoInfo_Success tests successful photo creation
func TestCreatePhotoInfo_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photo := domain.Photo{
		ID:             "photo123",
		Name:           "test.jpg",
		FilesystemPath: "/path/to/test.jpg",
		AlbumId:        "album1",
		Tags:           "",
		Metadata:       []byte("{}"),
		Thumbnail:      []byte("thumbnail-data"),
		CreatedEpoch:   time.Now().UnixMilli(),
	}

	// Expect INSERT with ON CONFLICT DO UPDATE
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "photos"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewPhotoRepository(env)
	err := repo.CreatePhotoInfo(photo)

	if err != nil {
		t.Errorf("error was not expected while creating photo: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreatePhotoInfo_Upsert tests photo upsert on conflict
func TestCreatePhotoInfo_Upsert(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photo := domain.Photo{
		ID:             "existingPhoto",
		Name:           "updated.jpg",
		FilesystemPath: "/path/to/updated.jpg",
		AlbumId:        "album1",
		Tags:           "updated",
		Metadata:       []byte("{}"),
		Thumbnail:      []byte("new-thumbnail"),
		CreatedEpoch:   time.Now().UnixMilli(),
	}

	// Expect INSERT with ON CONFLICT DO UPDATE (upsert behavior)
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "photos"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewPhotoRepository(env)
	err := repo.CreatePhotoInfo(photo)

	if err != nil {
		t.Errorf("error was not expected while upserting photo: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}