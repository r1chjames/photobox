package database

import (
	"errors"
	"fmt"
	. "gitlab.com/r1chjames/photobox/api/internal/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Env struct {
	db *gorm.DB
}

func InitDbConnection(appConfig AppConfig) *Env {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s?charset=utf8&parseTime=True&loc=Local", appConfig.DbUrl),
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}
	return &Env{db: db}
}

func (dbEnv *Env) PerformDbSetup(appConfig AppConfig) {
	// Migrate the schema
	err := dbEnv.db.AutoMigrate(&Album{}, &Photo{}, &Setting{}, &Job{})
	if err != nil {
		log.Fatal("failed to perform database migration")
	}

	dbEnv.createBaseSettings(appConfig.ResetSettings)
	dbEnv.createBaseJobs()
}

func checkNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func Paginate(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		var offset int
		if page == 1 {
			offset = 0
		} else {
			offset = (page - 1) * limit
		}

		return db.Offset(offset).Limit(limit)
	}
}
