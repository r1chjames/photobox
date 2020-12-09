package database

import (
	"errors"
	"fmt"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var dbConn *gorm.DB

func InitDbConnection(appConfig AppConfig) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		//DriverName: "mysql_driver",
		DSN:        fmt.Sprintf("%s?charset=utf8&parseTime=True&loc=Local", appConfig.DbUrl),
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	dbConn = db
}

func PerformDbSetup() {
	// Migrate the schema
	err := dbConn.AutoMigrate(&Album{}, &Photo{}, &Setting{}, &Job{})
	if err != nil {
		log.Fatal("failed to perform database migration")
	}

	dbConn.Create(&Job{Name: "Photo_index", Status: "NOT_RUNNING"})
}

func checkNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}