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
	expectedPhoto := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"}).
		AddRow(photoID, "test.jpg", "/path/to/test.jpg", "album1", []byte("{}"), []byte("thumbnail-data"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`, expectedPhoto)

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
	expectedPhoto := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "metadata", "created_at", "created_epoch", "updated_at"}).
		AddRow(photoID, "test2.jpg", "/path/to/test2.jpg", "", "album2", []byte("{}"), time.Now(), time.Now().UnixMilli(), time.Now())

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
	database.ShouldReturnNotFoundErrorForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`)

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

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo1", "test1.jpg", "/path/to/test1.jpg", "album1", []byte("{}"), []byte("thumb1"), time.Now().UnixMilli()).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", "album1", []byte("{}"), []byte("thumb2"), time.Now().UnixMilli()).
		AddRow("photo3", "test3.jpg", "/path/to/test3.jpg", "album2", []byte("{}"), []byte("thumb3"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NULL AND hidden = \$1 ORDER BY created_epoch ASC LIMIT \$2`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotos("", 10, true, "", "", "")

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
	fromPhotoRow := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "metadata", "created_at", "created_epoch", "updated_at"}).
		AddRow(fromID, "test1.jpg", "/path/to/test1.jpg", "", "album1", []byte("{}"), time.Now(), fromEpoch, time.Now())
	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`, fromPhotoRow)

	// Second query to get photos after fromEpoch - also uses Omit("thumbnail")
	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "metadata", "created_at", "created_epoch", "updated_at"}).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", "", "album1", []byte("{}"), time.Now(), fromEpoch+1000, time.Now()).
		AddRow("photo3", "test3.jpg", "/path/to/test3.jpg", "", "album2", []byte("{}"), time.Now(), fromEpoch+2000, time.Now())

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NULL AND hidden = \$1 AND created_epoch > \$2 ORDER BY created_epoch ASC LIMIT \$3`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotos(fromID, 5, false, "", "", "")

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

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"})

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NULL AND hidden = \$1 ORDER BY created_epoch ASC LIMIT \$2`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotos("", 10, true, "", "", "")

	if err != nil {
		t.Errorf("error was not expected when no photos exist: %s", err)
	}

	if len(photos) != 0 {
		t.Errorf("expected 0 photos when list is empty, got %d", len(photos))
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
	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo1", "test1.jpg", "/path/to/test1.jpg", albumID, []byte("{}"), []byte("thumb1"), time.Now().UnixMilli()).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", albumID, []byte("{}"), []byte("thumb2"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE \(deleted_at IS NULL AND album_id = \$1\) AND hidden = \$2 ORDER BY created_epoch ASC LIMIT \$3`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotosInAlbum(albumID, "", 10, true, "", "", "")

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

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", albumID, []byte("{}"), []byte("thumb2"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE \(deleted_at IS NULL AND album_id = \$1\) AND hidden = \$2 AND created_epoch > \$3 ORDER BY created_epoch ASC LIMIT \$4`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotosInAlbum(albumID, fromID, 5, true, "", "", "")

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
	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"})

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE \(deleted_at IS NULL AND album_id = \$1\) AND hidden = \$2 ORDER BY created_epoch ASC LIMIT \$3`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListAllPhotosInAlbum(albumID, "", 10, true, "", "", "")

	if err != nil {
		t.Errorf("error was not expected when album is empty: %s", err)
	}

	if len(photos) != 0 {
		t.Errorf("expected 0 photos when album is empty, got %d", len(photos))
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

// TestSoftDeletePhoto_Success tests successful soft deletion of a photo
func TestSoftDeletePhoto_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photoID := "photo123"
	now := time.Now()

	expectedPhoto := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "metadata", "created_at", "created_epoch", "updated_at", "favorite", "deleted_at", "blurhash", "dominant_color"}).
		AddRow(photoID, "test.jpg", "/path/to/test.jpg", "", "album1", []byte("{}"), now, now.UnixMilli(), now, false, nil, "", "")

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`, expectedPhoto)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "photos" SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewPhotoRepository(env)
	photo, err := repo.SoftDeletePhoto(photoID)

	if err != nil {
		t.Errorf("error was not expected while soft deleting photo: %s", err)
	}

	if photo == nil {
		t.Error("expected photo to be returned, got nil")
	}

	if photo.ID != photoID {
		t.Errorf("expected photo ID %s, got %s", photoID, photo.ID)
	}

	if photo.DeletedAt == nil {
		t.Error("expected photo to have DeletedAt set")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestSoftDeletePhoto_NotFound tests soft deletion when photo doesn't exist
func TestSoftDeletePhoto_NotFound(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photoID := "nonexistent"
	database.ShouldReturnNotFoundErrorForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`)

	repo := NewPhotoRepository(env)
	photo, err := repo.SoftDeletePhoto(photoID)

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

// TestRestorePhoto_Success tests successful restoration of a soft-deleted photo
func TestRestorePhoto_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photoID := "photo123"
	now := time.Now()
	deletedAt := now.Add(-time.Hour)

	expectedPhoto := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "metadata", "created_at", "created_epoch", "updated_at", "favorite", "deleted_at", "blurhash", "dominant_color"}).
		AddRow(photoID, "test.jpg", "/path/to/test.jpg", "", "album1", []byte("{}"), now, now.UnixMilli(), now, false, deletedAt, "", "")

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`, expectedPhoto)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "photos" SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewPhotoRepository(env)
	photo, err := repo.RestorePhoto(photoID)

	if err != nil {
		t.Errorf("error was not expected while restoring photo: %s", err)
	}

	if photo == nil {
		t.Error("expected photo to be returned, got nil")
	}

	if photo.ID != photoID {
		t.Errorf("expected photo ID %s, got %s", photoID, photo.ID)
	}

	if photo.DeletedAt != nil {
		t.Error("expected photo DeletedAt to be nil after restore")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestRestorePhoto_NotFound tests restoration when photo doesn't exist
func TestRestorePhoto_NotFound(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photoID := "nonexistent"
	database.ShouldReturnNotFoundErrorForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`)

	repo := NewPhotoRepository(env)
	photo, err := repo.RestorePhoto(photoID)

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

// TestListTrashPhotos_Success tests successful listing of trash photos
func TestListTrashPhotos_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo1", "test1.jpg", "/path/to/test1.jpg", "album1", []byte("{}"), []byte("thumb1"), time.Now().UnixMilli()).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", "album1", []byte("{}"), []byte("thumb2"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NOT NULL ORDER BY created_epoch ASC LIMIT \$1`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListTrashPhotos("", 10, true)

	if err != nil {
		t.Errorf("error was not expected while listing trash photos: %s", err)
	}

	expectedCount := 2
	if len(photos) != expectedCount {
		t.Errorf("expected %d photos, got %d", expectedCount, len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListTrashPhotos_WithFromId tests listing trash photos with pagination
func TestListTrashPhotos_WithFromId(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	fromID := "photo1"
	fromEpoch := time.Now().UnixMilli()

	fromPhotoRow := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "metadata", "created_at", "created_epoch", "updated_at", "favorite", "deleted_at", "blurhash", "dominant_color"}).
		AddRow(fromID, "test1.jpg", "/path/to/test1.jpg", "", "album1", []byte("{}"), time.Now(), fromEpoch, time.Now(), false, nil, "", "")

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE "photos"\."id" = \$1 ORDER BY "photos"\."id" LIMIT \$2`, fromPhotoRow)

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", "album1", []byte("{}"), []byte("thumb2"), fromEpoch+1000)

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NOT NULL AND created_epoch > \$1 ORDER BY created_epoch ASC LIMIT \$2`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListTrashPhotos(fromID, 5, true)

	if err != nil {
		t.Errorf("error was not expected while listing trash photos: %s", err)
	}

	if len(photos) != 1 {
		t.Errorf("expected 1 photo, got %d", len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListTrashPhotos_Empty tests listing when trash is empty
func TestListTrashPhotos_Empty(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"})

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NOT NULL ORDER BY created_epoch ASC LIMIT \$1`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListTrashPhotos("", 10, true)

	if err != nil {
		t.Errorf("error was not expected when trash is empty: %s", err)
	}

	if len(photos) != 0 {
		t.Errorf("expected 0 photos when trash is empty, got %d", len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestEmptyTrash_Success tests successful emptying of trash
func TestEmptyTrash_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "photos" WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()

	repo := NewPhotoRepository(env)
	err := repo.EmptyTrash()

	if err != nil {
		t.Errorf("error was not expected while emptying trash: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdatePhoto_Success tests successful photo update
func TestUpdatePhoto_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	photo := domain.Photo{
		ID:             "photo123",
		Name:           "updated.jpg",
		FilesystemPath: "/path/to/updated.jpg",
		AlbumId:        "album1",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "photos" SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewPhotoRepository(env)
	err := repo.UpdatePhoto(photo)

	if err != nil {
		t.Errorf("error was not expected while updating photo: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListFavoritePhotos_Success tests successful listing of favorite photos
func TestListFavoritePhotos_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo1", "test1.jpg", "/path/to/test1.jpg", "album1", []byte("{}"), []byte("thumb1"), time.Now().UnixMilli()).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", "album2", []byte("{}"), []byte("thumb2"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE \(favorite = \$1 AND deleted_at IS NULL\) AND hidden = \$2 ORDER BY created_epoch ASC LIMIT \$3`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListFavoritePhotos("", 10, true, "", "")

	if err != nil {
		t.Errorf("error was not expected while listing favorite photos: %s", err)
	}

	expectedCount := 2
	if len(photos) != expectedCount {
		t.Errorf("expected %d photos, got %d", expectedCount, len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListFavoritePhotos_Empty tests listing when no favorite photos exist
func TestListFavoritePhotos_Empty(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"})

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE \(favorite = \$1 AND deleted_at IS NULL\) AND hidden = \$2 ORDER BY created_epoch ASC LIMIT \$3`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.ListFavoritePhotos("", 10, true, "", "")

	if err != nil {
		t.Errorf("error was not expected when no favorite photos exist: %s", err)
	}

	if len(photos) != 0 {
		t.Errorf("expected 0 photos when no favorite photos exist, got %d", len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestSearchPhotos_Success tests successful photo search
func TestSearchPhotos_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	query := "vacation"
	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"}).
		AddRow("photo1", "vacation1.jpg", "/path/to/vacation1.jpg", "album1", []byte("{}"), []byte("thumb1"), time.Now().UnixMilli()).
		AddRow("photo2", "vacation2.jpg", "/path/to/vacation2.jpg", "album1", []byte("{}"), []byte("thumb2"), time.Now().UnixMilli())

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NULL AND hidden = \$1 AND to_tsvector\('english', coalesce\(name, ''\) \|\| ' ' \|\| coalesce\(tags, ''\)\) @@ plainto_tsquery\('english', \$2\) ORDER BY created_epoch DESC LIMIT \$3`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.SearchPhotos(query, 10)

	if err != nil {
		t.Errorf("error was not expected while searching photos: %s", err)
	}

	expectedCount := 2
	if len(photos) != expectedCount {
		t.Errorf("expected %d photos, got %d", expectedCount, len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestSearchPhotos_Empty tests search with no results
func TestSearchPhotos_Empty(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	query := "nonexistent"
	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "album_id", "metadata", "thumbnail", "created_epoch"})

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NULL AND hidden = \$1 AND to_tsvector\('english', coalesce\(name, ''\) \|\| ' ' \|\| coalesce\(tags, ''\)\) @@ plainto_tsquery\('english', \$2\) ORDER BY created_epoch DESC LIMIT \$3`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.SearchPhotos(query, 10)

	if err != nil {
		t.Errorf("error was not expected when search returns no results: %s", err)
	}

	if len(photos) != 0 {
		t.Errorf("expected 0 photos when search is empty, got %d", len(photos))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetTimeline_Success tests successful timeline retrieval
func TestGetTimeline_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedEntries := sqlmock.NewRows([]string{"year", "month", "count"}).
		AddRow(2024, 1, 10).
		AddRow(2024, 2, 15).
		AddRow(2023, 12, 5)

	database.ShouldReturnRowsForQuery(mock, `SELECT "year","month",COUNT\(\*\) as count FROM "photos" WHERE deleted_at IS NULL AND year > 0 AND month > 0 GROUP BY year, month ORDER BY year DESC, month DESC`, expectedEntries)

	repo := NewPhotoRepository(env)
	entries, err := repo.GetTimeline()

	if err != nil {
		t.Errorf("error was not expected while getting timeline: %s", err)
	}

	expectedCount := 3
	if len(entries) != expectedCount {
		t.Errorf("expected %d timeline entries, got %d", expectedCount, len(entries))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetPhotosWithGeodata_Success tests successful retrieval of photos with geodata
func TestGetPhotosWithGeodata_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedPhotos := sqlmock.NewRows([]string{"id", "name", "filesystem_path", "source_path", "album_id", "metadata", "created_at", "created_epoch", "updated_at", "favorite", "deleted_at", "blurhash", "dominant_color", "latitude", "longitude"}).
		AddRow("photo1", "test1.jpg", "/path/to/test1.jpg", "source1", "album1", []byte("{}"), time.Now(), time.Now().UnixMilli(), time.Now(), false, nil, "", "", 51.5, -0.1).
		AddRow("photo2", "test2.jpg", "/path/to/test2.jpg", "source2", "album1", []byte("{}"), time.Now(), time.Now().UnixMilli(), time.Now(), false, nil, "", "", 40.7, -74.0)

	database.ShouldReturnRowsForQuery(mock, `SELECT .+ FROM "photos" WHERE deleted_at IS NULL AND hidden = \$1 AND \(latitude IS NOT NULL AND longitude IS NOT NULL\) AND \(latitude BETWEEN \$2 AND \$3\) AND \(longitude BETWEEN \$4 AND \$5\) LIMIT \$6`, expectedPhotos)

	repo := NewPhotoRepository(env)
	photos, err := repo.GetPhotosWithGeodata(52.0, 50.0, 1.0, -1.0)

	if err != nil {
		t.Errorf("error was not expected while getting photos with geodata: %s", err)
	}

	expectedCount := 2
	if len(photos) != expectedCount {
		t.Errorf("expected %d photos, got %d", expectedCount, len(photos))
	}

	// Verify GPS data is populated
	if photos[0].Lat != 51.5 || photos[0].Lng != -0.1 {
		t.Errorf("expected photo1 lat=51.5, lng=-0.1, got lat=%f, lng=%f", photos[0].Lat, photos[0].Lng)
	}
	if photos[1].Lat != 40.7 || photos[1].Lng != -74.0 {
		t.Errorf("expected photo2 lat=40.7, lng=-74.0, got lat=%f, lng=%f", photos[1].Lat, photos[1].Lng)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}