package database

import (
	"database/sql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)


func mockDB(t *testing.T) (*Env, sqlmock.Sqlmock, *sql.DB) {
	dbConn, mock, err := sqlmock.New()
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: "test_db", Conn: dbConn, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return &Env{db: db}, mock, dbConn
}

func TestShouldGetAlbumByID(t *testing.T) {
	env, mock, dbConn := mockDB(t)
	defer dbConn.Close()

	mock.ExpectQuery("SELECT (.+) FROM `albums` WHERE `albums`.`id` = (.+) ORDER BY `albums`.`id` LIMIT 1").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("a", "b", "c"))

	if _, err := env.GetAlbumById("1"); err != nil {
		t.Errorf("error was not expected while getting albums: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldGetAlbumByName(t *testing.T) {
	env, mock, dbConn := mockDB(t)
	defer dbConn.Close()

	mock.ExpectQuery("SELECT (.+) FROM `albums` WHERE `albums`.`name` = (.+) ORDER BY `albums`.`id` LIMIT 1").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("a", "b", "c"))

	if _, err := env.GetAlbumByName("Photos"); err != nil {
		t.Errorf("error was not expected while getting albums: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}