package database

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Env struct {
	Db *gorm.DB
}

func InitDbConnection(appConfig *appconfig.AppConfig) *Env {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: appConfig.DbUrl,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "photobox.",
			SingularTable: false,
		}})

	if err != nil {
		slog.Error("Failed to connect database", "error", err)
		os.Exit(1)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get database instance", "error", err)
		os.Exit(1)
	}

	// Set maximum number of idle connections in the pool
	sqlDB.SetMaxIdleConns(10)

	// Set maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// Set maximum lifetime of a connection (reuse connections for up to 1 hour)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Env{Db: db}
}

func (dbEnv *Env) PerformDbSetup() {
	// Migrate the schema
	err := dbEnv.Db.AutoMigrate(&domain.Album{}, &domain.Photo{}, &domain.Setting{}, &domain.Job{}, &domain.User{})
	if err != nil {
		slog.Error("Failed to perform database migration", "error", err)
		os.Exit(1)
	}
}

// Ping checks the database connection is alive
func (dbEnv *Env) Ping() error {
	sqlDB, err := dbEnv.Db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func HandleError(result *gorm.DB) error {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.ErrDataNotFound
	} else if result.Error != nil {
		return result.Error
	}
	return nil
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
