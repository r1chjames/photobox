package appconfig

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNew_WithDefaultValues tests config creation with default values
func TestNew_WithDefaultValues(t *testing.T) {
	// Clean environment
	envVars := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"PHOTO_DIR", "API_BASE_PATH", "RESET_SETTINGS", "DEBUG_MODE",
		"TIMEZONE", "TOKEN_DURATION", "DEFAULT_ADMIN_USERNAME",
		"DEFAULT_ADMIN_PASSWORD", "CACHE_HOST", "CACHE_PORT", "CACHE_PASSWORD",
		"CACHE_ENABLED", "CACHE_DB",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}

	// TOKEN is required, set a default
	os.Setenv("TOKEN", "test-token")
	os.Setenv("DEFAULT_ADMIN_PASSWORD", "password")
	defer os.Unsetenv("TOKEN")
	defer os.Unsetenv("DEFAULT_ADMIN_PASSWORD")

	config := New()

	assert.NotNil(t, config)
	assert.Equal(t, "/photos", config.PhotoDir)
	assert.Equal(t, "/api", config.ApiBasePath)
	assert.Contains(t, config.DbUrl, "host=localhost")
	assert.Contains(t, config.DbUrl, "port=5432")
	assert.Contains(t, config.DbUrl, "user=photobox")
	assert.Contains(t, config.DbUrl, "dbname=photobox")
	assert.False(t, config.ResetSettings)
	assert.False(t, config.DebugMode)
	assert.NotNil(t, config.Timezone)
	assert.Equal(t, "admin", config.AdminUsername)
	assert.Equal(t, "password", config.AdminPassword)
	assert.Equal(t, time.Hour, config.TokenDuration)
	assert.Equal(t, "test-token", config.Token)
}

// TestNew_WithCustomValues tests config creation with custom environment values
func TestNew_WithCustomValues(t *testing.T) {
	// Set custom environment variables
	os.Setenv("DB_HOST", "custom-host")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("PHOTO_DIR", "/custom/photos")
	os.Setenv("API_BASE_PATH", "/custom/api")
	os.Setenv("RESET_SETTINGS", "true")
	os.Setenv("DEBUG_MODE", "true")
	os.Setenv("TIMEZONE", "America/New_York")
	os.Setenv("TOKEN", "custom-token")
	os.Setenv("TOKEN_DURATION", "2h")
	os.Setenv("DEFAULT_ADMIN_USERNAME", "customadmin")
	os.Setenv("DEFAULT_ADMIN_PASSWORD", "custompass")

	defer func() {
		// Cleanup
	envVars := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"PHOTO_DIR", "API_BASE_PATH", "RESET_SETTINGS", "DEBUG_MODE",
		"TIMEZONE", "TOKEN_DURATION", "DEFAULT_ADMIN_USERNAME",
		"DEFAULT_ADMIN_PASSWORD", "CACHE_HOST", "CACHE_PORT", "CACHE_PASSWORD",
		"CACHE_ENABLED", "CACHE_DB",
		"THUMBNAIL_STORAGE", "THUMBNAIL_DIR",
		"S3_ENDPOINT", "S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_BUCKET", "S3_USE_SSL",
	}
		for _, v := range envVars {
			os.Unsetenv(v)
		}
	}()

	config := New()

	assert.NotNil(t, config)
	assert.Equal(t, "/custom/photos", config.PhotoDir)
	assert.Equal(t, "/custom/api", config.ApiBasePath)
	assert.Contains(t, config.DbUrl, "host=custom-host")
	assert.Contains(t, config.DbUrl, "port=3306")
	assert.Contains(t, config.DbUrl, "user=testuser")
	assert.Contains(t, config.DbUrl, "password=testpass")
	assert.Contains(t, config.DbUrl, "dbname=testdb")
	assert.True(t, config.ResetSettings)
	assert.True(t, config.DebugMode)
	assert.Equal(t, "America/New_York", config.Timezone.String())
	assert.Equal(t, "custom-token", config.Token)
	assert.Equal(t, 2*time.Hour, config.TokenDuration)
	assert.Equal(t, "customadmin", config.AdminUsername)
	assert.Equal(t, "custompass", config.AdminPassword)
}

// TestNew_WithPartialCustomValues tests config with some custom and some default values
func TestNew_WithPartialCustomValues(t *testing.T) {
	// Clean environment first
	envVars := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"PHOTO_DIR", "API_BASE_PATH", "RESET_SETTINGS", "DEBUG_MODE",
		"TIMEZONE", "TOKEN_DURATION", "DEFAULT_ADMIN_USERNAME",
		"DEFAULT_ADMIN_PASSWORD", "CACHE_HOST", "CACHE_PORT", "CACHE_PASSWORD",
		"CACHE_ENABLED", "CACHE_DB",
		"THUMBNAIL_STORAGE", "THUMBNAIL_DIR",
		"S3_ENDPOINT", "S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_BUCKET", "S3_USE_SSL",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}

	// Set only some variables (TOKEN is required)
	os.Setenv("TOKEN", "test-token")
	os.Setenv("DEFAULT_ADMIN_PASSWORD", "password")
	os.Setenv("DB_HOST", "partial-host")
	os.Setenv("PHOTO_DIR", "/partial/photos")
	os.Setenv("DEBUG_MODE", "true")

	defer func() {
		os.Unsetenv("TOKEN")
		os.Unsetenv("DEFAULT_ADMIN_PASSWORD")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("PHOTO_DIR")
		os.Unsetenv("DEBUG_MODE")
	}()

	config := New()

	assert.NotNil(t, config)
	assert.Equal(t, "/partial/photos", config.PhotoDir)
	assert.Contains(t, config.DbUrl, "host=partial-host")
	assert.Contains(t, config.DbUrl, "port=5432") // Default
	assert.Contains(t, config.DbUrl, "user=photobox") // Default
	assert.True(t, config.DebugMode)
	assert.False(t, config.ResetSettings) // Default
}

// TestNew_DbUrlFormat tests the database URL format
func TestNew_DbUrlFormat(t *testing.T) {
	os.Setenv("TOKEN", "test-token")
	os.Setenv("DEFAULT_ADMIN_PASSWORD", "password")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")

	defer func() {
		os.Unsetenv("TOKEN")
		os.Unsetenv("DEFAULT_ADMIN_PASSWORD")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
	}()

	config := New()

	expectedDbUrl := "host=testhost user=testuser password=testpass dbname=testdb port=5433 sslmode=require"
	assert.Equal(t, expectedDbUrl, config.DbUrl)
}

// TestNew_InvalidBooleanValues tests config with invalid boolean values (defaults to false)
func TestNew_InvalidBooleanValues(t *testing.T) {
	os.Setenv("TOKEN", "test-token")
	os.Setenv("DEFAULT_ADMIN_PASSWORD", "password")
	os.Setenv("RESET_SETTINGS", "invalid")
	os.Setenv("DEBUG_MODE", "not-a-bool")

	defer func() {
		os.Unsetenv("TOKEN")
		os.Unsetenv("DEFAULT_ADMIN_PASSWORD")
		os.Unsetenv("RESET_SETTINGS")
		os.Unsetenv("DEBUG_MODE")
	}()

	config := New()

	// Invalid booleans should parse to false
	assert.False(t, config.ResetSettings)
	assert.False(t, config.DebugMode)
}

// TestNew_InvalidTimezone tests config with invalid timezone (results in nil)
func TestNew_InvalidTimezone(t *testing.T) {
	os.Setenv("TOKEN", "test-token")
	os.Setenv("DEFAULT_ADMIN_PASSWORD", "password")
	os.Setenv("TIMEZONE", "Invalid/Timezone")

	defer func() {
		os.Unsetenv("TOKEN")
		os.Unsetenv("DEFAULT_ADMIN_PASSWORD")
		os.Unsetenv("TIMEZONE")
	}()

	config := New()

	// Invalid timezone results in nil (error is ignored in appconfig.New)
	assert.Nil(t, config.Timezone)
}

// TestNew_InvalidTokenDuration tests config with invalid duration (falls back to zero)
func TestNew_InvalidTokenDuration(t *testing.T) {
	os.Setenv("TOKEN", "test-token")
	os.Setenv("DEFAULT_ADMIN_PASSWORD", "password")
	os.Setenv("TOKEN_DURATION", "invalid-duration")

	defer func() {
		os.Unsetenv("TOKEN")
		os.Unsetenv("DEFAULT_ADMIN_PASSWORD")
		os.Unsetenv("TOKEN_DURATION")
	}()

	config := New()

	// Invalid duration should parse to zero
	assert.Equal(t, time.Duration(0), config.TokenDuration)
}

// NOTE: TestNew_EmptyToken skipped - TOKEN is required and panics when empty
// The application requires TOKEN to be set. Empty token causes panic in utils.GetEnv
// func TestNew_EmptyToken(t *testing.T) {
// 	t.Skip("TOKEN is required, empty value causes panic")
// }
