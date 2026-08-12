package appconfig

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// testToken is exactly 32 bytes, matching the Paseto chacha20poly1305 key size.
const testToken = "0123456789abcdef0123456789abcdef"

// testAdminPassword is a strong, non-denylisted password.
const testAdminPassword = "Str0ngAdminPass!"

// clearEnv unsets every configuration variable the loader reads.
func clearEnv() {
	envVars := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE",
		"PHOTO_DIR", "API_BASE_PATH", "RESET_SETTINGS", "DEBUG_MODE",
		"TIMEZONE", "TOKEN_DURATION", "DEFAULT_ADMIN_USERNAME",
		"DEFAULT_ADMIN_PASSWORD", "CACHE_HOST", "CACHE_PORT", "CACHE_PASSWORD",
		"CACHE_ENABLED", "CACHE_DB", "AI_ENABLED", "OLLAMA_HOST", "OLLAMA_MODEL",
		"THUMBNAIL_STORAGE", "THUMBNAIL_DIR",
		"S3_ENDPOINT", "S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_BUCKET", "S3_USE_SSL",
		"CORS_ALLOWED_ORIGINS", "TOKEN",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}
}

func setRequired(t *testing.T) {
	t.Helper()
	os.Setenv("TOKEN", testToken)
	os.Setenv("DEFAULT_ADMIN_PASSWORD", testAdminPassword)
	t.Cleanup(func() {
		os.Unsetenv("TOKEN")
		os.Unsetenv("DEFAULT_ADMIN_PASSWORD")
	})
}

// TestNew_WithDefaultValues tests config creation with default values
func TestNew_WithDefaultValues(t *testing.T) {
	clearEnv()
	setRequired(t)

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
	assert.Equal(t, testAdminPassword, config.AdminPassword)
	assert.Equal(t, time.Hour, config.TokenDuration)
	assert.Equal(t, testToken, config.Token)
}

// TestNew_WithCustomValues tests config creation with custom environment values
func TestNew_WithCustomValues(t *testing.T) {
	clearEnv()
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
	os.Setenv("TOKEN", testToken)
	os.Setenv("TOKEN_DURATION", "2h")
	os.Setenv("DEFAULT_ADMIN_USERNAME", "customadmin")
	os.Setenv("DEFAULT_ADMIN_PASSWORD", "customadminpass1")
	t.Cleanup(clearEnv)

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
	assert.Equal(t, testToken, config.Token)
	assert.Equal(t, 2*time.Hour, config.TokenDuration)
	assert.Equal(t, "customadmin", config.AdminUsername)
	assert.Equal(t, "customadminpass1", config.AdminPassword)
}

// TestNew_WithPartialCustomValues tests config with some custom and some default values
func TestNew_WithPartialCustomValues(t *testing.T) {
	clearEnv()
	os.Setenv("TOKEN", testToken)
	os.Setenv("DEFAULT_ADMIN_PASSWORD", testAdminPassword)
	os.Setenv("DB_HOST", "partial-host")
	os.Setenv("PHOTO_DIR", "/partial/photos")
	os.Setenv("DEBUG_MODE", "true")
	t.Cleanup(clearEnv)

	config := New()

	assert.NotNil(t, config)
	assert.Equal(t, "/partial/photos", config.PhotoDir)
	assert.Contains(t, config.DbUrl, "host=partial-host")
	assert.Contains(t, config.DbUrl, "port=5432")   // Default
	assert.Contains(t, config.DbUrl, "user=photobox") // Default
	assert.True(t, config.DebugMode)
	assert.False(t, config.ResetSettings) // Default
}

// TestNew_DbUrlFormat tests the database URL format
func TestNew_DbUrlFormat(t *testing.T) {
	clearEnv()
	os.Setenv("TOKEN", testToken)
	os.Setenv("DEFAULT_ADMIN_PASSWORD", testAdminPassword)
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	t.Cleanup(clearEnv)

	config := New()

	expectedDbUrl := "host=testhost user=testuser password=testpass dbname=testdb port=5433 sslmode=require search_path=photobox,public"
	assert.Equal(t, expectedDbUrl, config.DbUrl)
}

// TestNew_InvalidBooleanValues tests config with invalid boolean values (defaults to false)
func TestNew_InvalidBooleanValues(t *testing.T) {
	clearEnv()
	os.Setenv("TOKEN", testToken)
	os.Setenv("DEFAULT_ADMIN_PASSWORD", testAdminPassword)
	os.Setenv("RESET_SETTINGS", "invalid")
	os.Setenv("DEBUG_MODE", "not-a-bool")
	t.Cleanup(clearEnv)

	config := New()

	assert.False(t, config.ResetSettings)
	assert.False(t, config.DebugMode)
}

// TestNew_InvalidTimezone tests config with invalid timezone (results in nil)
func TestNew_InvalidTimezone(t *testing.T) {
	clearEnv()
	os.Setenv("TOKEN", testToken)
	os.Setenv("DEFAULT_ADMIN_PASSWORD", testAdminPassword)
	os.Setenv("TIMEZONE", "Invalid/Timezone")
	t.Cleanup(clearEnv)

	config := New()

	assert.Nil(t, config.Timezone)
}

// TestNew_InvalidTokenDuration tests config with invalid duration (falls back to zero)
func TestNew_InvalidTokenDuration(t *testing.T) {
	clearEnv()
	os.Setenv("TOKEN", testToken)
	os.Setenv("DEFAULT_ADMIN_PASSWORD", testAdminPassword)
	os.Setenv("TOKEN_DURATION", "invalid-duration")
	t.Cleanup(clearEnv)

	config := New()

	assert.Equal(t, time.Duration(0), config.TokenDuration)
}

// --- Validation tests ---

func loadConfig(t *testing.T) *AppConfig {
	t.Helper()
	clearEnv()
	os.Setenv("TOKEN", testToken)
	os.Setenv("DEFAULT_ADMIN_PASSWORD", testAdminPassword)
	os.Setenv("DB_PASSWORD", "testpass")
	t.Cleanup(clearEnv)
	return load()
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := loadConfig(t)

	err, warnings := cfg.validate()
	assert.NoError(t, err)
	assert.Empty(t, warnings)
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := loadConfig(t)
	cfg.Token = ""

	err, _ := cfg.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TOKEN")
}

func TestValidate_ShortToken(t *testing.T) {
	cfg := loadConfig(t)
	cfg.Token = "short"

	err, _ := cfg.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")
}

func TestValidate_MissingAdminPassword(t *testing.T) {
	cfg := loadConfig(t)
	cfg.AdminPassword = ""

	err, _ := cfg.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DEFAULT_ADMIN_PASSWORD")
}

func TestValidate_ShortAdminPassword(t *testing.T) {
	cfg := loadConfig(t)
	cfg.AdminPassword = "short"

	err, _ := cfg.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least")
}

func TestValidate_WeakAdminPassword(t *testing.T) {
	for _, weak := range []string{"password", "admin1234", "admin", "changeme"} {
		cfg := loadConfig(t)
		cfg.AdminPassword = weak

		err, _ := cfg.validate()
		assert.Error(t, err, "expected %q to be rejected", weak)
		assert.Contains(t, err.Error(), "weak")
	}
}

func TestValidate_DefaultDBPasswordWarns(t *testing.T) {
	cfg := loadConfig(t)
	cfg.DbPassword = "photobox"

	err, warnings := cfg.validate()
	assert.NoError(t, err)
	assert.NotEmpty(t, warnings)
	assert.Contains(t, strings.Join(warnings, " "), "DB_PASSWORD")
}

func TestValidate_CacheEnabledWithoutHost(t *testing.T) {
	cfg := loadConfig(t)
	cfg.CacheEnabled = true
	cfg.CacheHost = ""

	err, _ := cfg.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CACHE_HOST")
}

func TestValidate_S3MissingCredentials(t *testing.T) {
	cfg := loadConfig(t)
	cfg.ThumbnailStorage = "s3"
	cfg.S3Endpoint = "http://minio:9000"
	cfg.S3AccessKey = ""
	cfg.S3SecretKey = ""
	cfg.S3Bucket = ""

	err, _ := cfg.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "S3_ACCESS_KEY")
	assert.Contains(t, err.Error(), "S3_SECRET_KEY")
	assert.Contains(t, err.Error(), "S3_BUCKET")
}

func TestValidate_S3Complete(t *testing.T) {
	cfg := loadConfig(t)
	cfg.ThumbnailStorage = "s3"
	cfg.S3Endpoint = "http://minio:9000"
	cfg.S3AccessKey = "key"
	cfg.S3SecretKey = "secret"
	cfg.S3Bucket = "bucket"

	err, _ := cfg.validate()
	assert.NoError(t, err)
}

func TestValidate_AIEnabledWithoutHost(t *testing.T) {
	cfg := loadConfig(t)
	cfg.AIEnabled = true
	cfg.OllamaHost = ""

	err, _ := cfg.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OLLAMA_HOST")
}

func TestValidate_AggregatesErrors(t *testing.T) {
	cfg := loadConfig(t)
	cfg.Token = ""
	cfg.AdminPassword = "weak"
	cfg.CacheEnabled = true
	cfg.CacheHost = ""

	err, _ := cfg.validate()
	assert.Error(t, err)
	// All three problems should appear in the single report.
	assert.Contains(t, err.Error(), "TOKEN")
	assert.Contains(t, err.Error(), "DEFAULT_ADMIN_PASSWORD")
	assert.Contains(t, err.Error(), "CACHE_HOST")
}

func TestRedact(t *testing.T) {
	assert.Equal(t, "<empty>", redact(""))
	assert.Equal(t, "***", redact("ab"))
	assert.Equal(t, "Se***", redact("Secret"))
}

func TestRedactURL(t *testing.T) {
	u := "host=db user=u password=hunter2 dbname=photobox"
	assert.Equal(t, "host=db user=u password=*** dbname=photobox", redactURL(u))

	u2 := "host=db password=hunter2"
	assert.Equal(t, "host=db password=***", redactURL(u2))

	// No password present -> unchanged
	assert.Equal(t, "host=db user=u", redactURL("host=db user=u"))
}
