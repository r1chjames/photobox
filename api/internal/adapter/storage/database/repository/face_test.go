package repository

import (
	"regexp"
	"testing"

	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var faceIDPattern = `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`

func TestInsertDetectionsForPhoto_Replaces(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	// delete-then-insert in a transaction (BEGIN, DELETE, INSERT, COMMIT)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "face_detections" WHERE photo_id = $1`)).
		WithArgs("photo1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`INSERT INTO "face_detections"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewFaceRepository(env)
	dets := []domain.FaceDetection{{
		ID: "face1", PhotoID: "photo1",
		BoxX: 10, BoxY: 20, BoxW: 30, BoxH: 40,
		Score: 0.9, Embedding: []byte("encrypted"), Status: "detected",
	}}
	if err := repo.InsertDetectionsForPhoto("photo1", dets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDeleteAllFaceData(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "face_clusters" WHERE 1 = 1`)).
		WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "face_persons" WHERE 1 = 1`)).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "face_detections" WHERE 1 = 1`)).
		WillReturnResult(sqlmock.NewResult(0, 9))
	mock.ExpectCommit()

	repo := NewFaceRepository(env)
	total, err := repo.DeleteAllFaceData()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 14 {
		t.Errorf("expected 14 deleted rows, got %d", total)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestPersons_WithPhotoCount(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "cover_face_id", "created_at", "updated_at", "photo_count"}).
		AddRow("p1", "Alice", "f1", nil, nil, 3).
		AddRow("p2", "Bob", "f9", nil, nil, 1)
	mock.ExpectQuery(`SELECT face_persons\.\*, COUNT\(DISTINCT fd\.photo_id\) AS photo_count FROM "face_persons" LEFT JOIN face_detections fd ON fd\.person_id = face_persons\.id GROUP BY "face_persons"\."id" ORDER BY face_persons\.name ASC`).
		WillReturnRows(rows)

	repo := NewFaceRepository(env)
	persons, err := repo.Persons()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(persons) != 2 {
		t.Fatalf("expected 2 persons, got %d", len(persons))
	}
	if persons[0].Name != "Alice" || persons[0].PhotoCount != 3 {
		t.Errorf("unexpected person[0]: %+v", persons[0])
	}
	if persons[1].Name != "Bob" || persons[1].PhotoCount != 1 {
		t.Errorf("unexpected person[1]: %+v", persons[1])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestPersonByID_NotFound(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	database.ShouldReturnNotFoundErrorForQuery(mock, `SELECT .+ FROM "face_persons" WHERE "face_persons"\."id" = \$1 ORDER BY "face_persons"\."id" LIMIT \$2`)

	repo := NewFaceRepository(env)
	if _, err := repo.PersonByID("missing"); err == nil {
		t.Error("expected error for missing person")
	}
}

func TestSetDetectionPerson_Assign(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "face_detections" SET .+ WHERE id = \$\d+`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewFaceRepository(env)
	pid := "p1"
	if err := repo.SetDetectionPerson("f1", &pid); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestUnassignedDetections(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"id", "photo_id", "person_id", "box_x", "box_y", "box_w", "box_h", "score", "status"}).
		AddRow("f1", "photo1", nil, 1, 2, 3, 4, 0.9, "detected")
	mock.ExpectQuery(`SELECT .+ FROM "face_detections" WHERE person_id IS NULL ORDER BY created_at ASC LIMIT \$1`).
		WillReturnRows(rows)

	repo := NewFaceRepository(env)
	dets, err := repo.UnassignedDetections(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dets) != 1 || dets[0].ID != "f1" {
		t.Errorf("unexpected detections: %+v", dets)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestPersonPhotoIDs(t *testing.T) {
	env, mock, dbConn := database.MockDB(t)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"photo_id"}).AddRow("photo1").AddRow("photo2")
	mock.ExpectQuery(`SELECT DISTINCT "photo_id" FROM "face_detections" WHERE person_id = \$1`).
		WillReturnRows(rows)

	repo := NewFaceRepository(env)
	ids, err := repo.PersonPhotoIDs("p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 || ids[0] != "photo1" || ids[1] != "photo2" {
		t.Errorf("unexpected ids: %v", ids)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
