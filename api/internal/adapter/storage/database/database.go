package database

import (
	"errors"
	. "gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

type Env struct {
	Db *gorm.DB
}

func InitDbConnection(appConfig *AppConfig) *Env {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: appConfig.DbUrl,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "photobox.",
			SingularTable: false,
		}})

	if err != nil {
		log.Fatalf("failed to connect database, %s", err)
	}
	return &Env{Db: db}
}

func (dbEnv *Env) PerformDbSetup(appConfig *AppConfig) {
	// Migrate the schema
	err := dbEnv.Db.AutoMigrate(&Album{}, &Photo{}, &Setting{}, &Job{}, &User{})
	if err != nil {
		log.Fatalf("failed to perform database migration, %s", err)
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
