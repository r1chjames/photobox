package testutil

import (
	"fmt"
	"testing"
	"time"

	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDBConfig returns a test database configuration
func TestDBConfig() appconfig.AppConfig {
	timezone, _ := time.LoadLocation("UTC")
	return appconfig.AppConfig{
		DbUrl:    "host=localhost user=photobox_test password=photobox_test dbname=photobox_test port=5433 sslmode=disable search_path=photobox,public",
		Timezone: timezone,
		PhotoDir: "/tmp/photobox-test",
	}
}

// SetupTestDB creates a test database connection and runs migrations
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	config := TestDBConfig()

	// Connect to database
	db, err := gorm.Open(postgres.Open(config.DbUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Ensure the photobox schema exists (mirrors production PerformDbSetup, so
	// tables land in the schema that TeardownTestDB truncates).
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS photobox").Error; err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&domain.Album{},
		&domain.Photo{},
		&domain.Setting{},
		&domain.Job{},
		&domain.User{},
	)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Clean all tables
	tables := []string{"photobox.photos", "photobox.albums", "photobox.users", "photobox.jobs", "photobox.settings"}
	for _, table := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
	}

	// Close connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// CreateTestEnv creates a database.Env for testing
func CreateTestEnv(t *testing.T) *database.Env {
	t.Helper()
	db := SetupTestDB(t)
	return &database.Env{Db: db}
}

// CleanupTestEnv cleans up the test environment
func CleanupTestEnv(t *testing.T, env *database.Env) {
	t.Helper()
	TeardownTestDB(t, env.Db)
}

// SeedTestData inserts common test data
func SeedTestData(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Create test albums
	albums := []domain.Album{
		{ID: "album1", Name: "Test Album 1", Description: "First test album"},
		{ID: "album2", Name: "Test Album 2", Description: "Second test album"},
	}
	for _, album := range albums {
		db.Create(&album)
	}

	// Create test user
	user := domain.User{
		ID:       "user1",
		Username: "testuser",
		Password: "$2a$10$test", // bcrypt hash
		Role:     domain.ADMINISTRATOR,
	}
	db.Create(&user)

	// Create test jobs
	jobs := []domain.Job{
		{Name: "photo_index", Status: "NOT_RUNNING"},
	}
	for _, job := range jobs {
		db.Create(&job)
	}
}
