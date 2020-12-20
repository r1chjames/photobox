package database

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func MockDB(t *testing.T) (*Env, sqlmock.Sqlmock, *sql.DB) {
	dbConn, mock, err := sqlmock.New()
	mockedDB, err := gorm.Open(mysql.New(mysql.Config{DSN: "test_db", Conn: dbConn, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return &Env{Db: mockedDB}, mock, dbConn
}

func ShouldReturnRowsForQuery(mock sqlmock.Sqlmock, query string, rows *sqlmock.Rows) {
	mock.ExpectQuery(query).WillReturnRows(rows)
}

func ShouldReturnRowsForUpdate(mock sqlmock.Sqlmock, statement string, rowsAffected int64) {
	mock.ExpectExec(statement).WillReturnResult(sqlmock.NewResult(1, rowsAffected))
}

func ShouldReturnNotFoundErrorForQuery(mock sqlmock.Sqlmock, query string) {
	mock.ExpectQuery(query).WillReturnError(gorm.ErrRecordNotFound)
}