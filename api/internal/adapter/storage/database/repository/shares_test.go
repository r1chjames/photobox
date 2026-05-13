package repository

import (
	"testing"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// TestCreateShare_Success tests successful share creation
func TestCreateShare_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	share := &domain.SharedLink{
		Token:        "token123",
		ResourceType: "photo",
		ResourceId:   "photo123",
		CreatedBy:    "user123",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "shared_links"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewShareRepository(env)
	err := repo.CreateShare(share)

	if err != nil {
		t.Errorf("error was not expected while creating share: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetShareByToken_Success tests successful share retrieval by token
func TestGetShareByToken_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	token := "token123"
	expectedShare := sqlmock.NewRows([]string{"token", "resource_type", "resource_id", "created_by", "view_count", "created_at", "updated_at"}).
		AddRow(token, "photo", "photo123", "user123", 0, time.Now(), time.Now())

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "shared_links" WHERE token = \$1 ORDER BY "shared_links"\."token" LIMIT \$2`, expectedShare)

	repo := NewShareRepository(env)
	share, err := repo.GetShareByToken(token)

	if err != nil {
		t.Errorf("error was not expected while getting share: %s", err)
	}

	if share == nil {
		t.Error("expected share to be returned, got nil")
	}

	if share.Token != token {
		t.Errorf("expected share token %s, got %s", token, share.Token)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetShareByToken_NotFound tests share not found scenario
func TestGetShareByToken_NotFound(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	token := "nonexistent"
	database.ShouldReturnNotFoundErrorForQuery(mock, `SELECT \* FROM "shared_links" WHERE token = \$1 ORDER BY "shared_links"\."token" LIMIT \$2`)

	repo := NewShareRepository(env)
	share, err := repo.GetShareByToken(token)

	if err == nil {
		t.Error("expected error when share not found, got nil")
	}

	if share != nil {
		t.Error("expected nil share when not found, got non-nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListShares_Success tests successful listing of shares
func TestListShares_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedShares := sqlmock.NewRows([]string{"token", "resource_type", "resource_id", "created_by", "view_count", "created_at", "updated_at"}).
		AddRow("token1", "photo", "photo123", "user123", 5, time.Now(), time.Now()).
		AddRow("token2", "album", "album123", "user123", 10, time.Now(), time.Now())

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "shared_links"`, expectedShares)

	repo := NewShareRepository(env)
	shares, err := repo.ListShares()

	if err != nil {
		t.Errorf("error was not expected while listing shares: %s", err)
	}

	expectedCount := 2
	if len(shares) != expectedCount {
		t.Errorf("expected %d shares, got %d", expectedCount, len(shares))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestListShares_Empty tests listing when no shares exist
func TestListShares_Empty(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	expectedShares := sqlmock.NewRows([]string{"token", "resource_type", "resource_id", "created_by", "view_count", "created_at", "updated_at"})

	database.ShouldReturnRowsForQuery(mock, `SELECT \* FROM "shared_links"`, expectedShares)

	repo := NewShareRepository(env)
	shares, err := repo.ListShares()

	if err != nil {
		t.Errorf("error was not expected for empty result: %s", err)
	}

	if shares == nil {
		t.Error("expected empty slice, got nil")
	}

	if len(shares) != 0 {
		t.Errorf("expected 0 shares, got %d", len(shares))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestDeleteShare_Success tests successful share deletion
func TestDeleteShare_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	token := "token123"

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "shared_links" WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewShareRepository(env)
	err := repo.DeleteShare(token)

	if err != nil {
		t.Errorf("error was not expected while deleting share: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestDeleteShare_NotFound tests deleting non-existent share
func TestDeleteShare_NotFound(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	token := "nonexistent"

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "shared_links" WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	repo := NewShareRepository(env)
	err := repo.DeleteShare(token)

	if err != nil {
		t.Errorf("error was not expected for non-existent share: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestIncrementViewCount_Success tests successful view count increment
func TestIncrementViewCount_Success(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	token := "token123"

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "shared_links" SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "shared_links" SET .+ WHERE .+`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewShareRepository(env)
	err := repo.IncrementViewCount(token)

	if err != nil {
		t.Errorf("error was not expected while incrementing view count: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
