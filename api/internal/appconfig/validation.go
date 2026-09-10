package appconfig

import (
	"fmt"
	"log/slog"
	"strings"
)

// weakAdminPasswords are rejected at startup to prevent an insecure default
// admin credential from being silently deployed.
var weakAdminPasswords = map[string]struct{}{
	"password":     {},
	"password123":  {},
	"admin":        {},
	"admin1234":    {},
	"changeme":     {},
	"change-me":    {},
	"photobox":     {},
	"letmein":      {},
	"12345678":     {},
	"qwerty":       {},
	"qwerty123":    {},
}

// minAdminPasswordLength is the minimum acceptable admin password length.
const minAdminPasswordLength = 8

// ConfigValidator accumulates configuration errors and warnings so they can be
// reported together in a single, aggregated report at startup rather than
// surfacing one-at-a-time or as panics.
type ConfigValidator struct {
	errors   []string
	warnings []string
}

// Required records an error when value is empty.
func (v *ConfigValidator) Required(key, value string) {
	if strings.TrimSpace(value) == "" {
		v.Errorf("missing required environment variable %s", key)
	}
}

// Errorf records a formatting error.
func (v *ConfigValidator) Errorf(format string, args ...any) {
	v.errors = append(v.errors, fmt.Sprintf(format, args...))
}

// Warnf records a non-fatal warning.
func (v *ConfigValidator) Warnf(format string, args ...any) {
	v.warnings = append(v.warnings, fmt.Sprintf(format, args...))
}

// Validate returns an aggregated error containing every recorded error, or nil
// when there are none.
func (v *ConfigValidator) Validate() error {
	if len(v.errors) == 0 {
		return nil
	}
	return fmt.Errorf("configuration errors (%d):\n  - %s", len(v.errors), strings.Join(v.errors, "\n  - "))
}

// validate checks the effective configuration and returns an aggregated error
// (or nil) plus any warnings. It is called by New before the application starts.
func (c *AppConfig) validate() (error, []string) {
	v := &ConfigValidator{}

	// TOKEN is the Paseto symmetric key and must be exactly 32 bytes.
	if c.Token == "" {
		v.Errorf("missing required environment variable TOKEN")
	} else if len(c.Token) != 32 {
		v.Errorf("TOKEN must be exactly 32 bytes (got %d)", len(c.Token))
	}

	// Admin password must be non-empty, long enough, and not obviously weak.
	if c.AdminPassword == "" {
		v.Errorf("missing required environment variable DEFAULT_ADMIN_PASSWORD")
	} else {
		if len(c.AdminPassword) < minAdminPasswordLength {
			v.Errorf("DEFAULT_ADMIN_PASSWORD must be at least %d characters", minAdminPasswordLength)
		}
		if _, weak := weakAdminPasswords[strings.ToLower(c.AdminPassword)]; weak {
			v.Errorf("DEFAULT_ADMIN_PASSWORD is a known-weak value and is rejected")
		}
	}

	// DB password is required (it only has a dev default in code).
	if c.DbPassword == "" {
		v.Errorf("missing required environment variable DB_PASSWORD")
	} else if c.DbPassword == "photobox" {
		v.Warnf("DB_PASSWORD is the insecure dev default; set a real password in production")
	}

	// Optional features must have their dependencies present.
	if c.CacheEnabled && strings.TrimSpace(c.CacheHost) == "" {
		v.Errorf("CACHE_ENABLED is true but CACHE_HOST is empty")
	}

	if c.ThumbnailStorage == "s3" || c.OriginalsStorage == "s3" {
		v.Required("S3_ENDPOINT", c.S3Endpoint)
		v.Required("S3_ACCESS_KEY", c.S3AccessKey)
		v.Required("S3_SECRET_KEY", c.S3SecretKey)
		v.Required("S3_BUCKET", c.S3Bucket)
	}

	if c.OriginalsStorage != "filesystem" && c.OriginalsStorage != "s3" {
		v.Errorf("ORIGINALS_STORAGE must be one of filesystem, s3 (got %q)", c.OriginalsStorage)
	}

	if c.AIEnabled && strings.TrimSpace(c.OllamaHost) == "" {
		v.Errorf("AI_ENABLED is true but OLLAMA_HOST is empty")
	}

	if c.FaceEngineEnabled && strings.TrimSpace(c.FaceEngineURL) == "" {
		v.Errorf("FACE_ENGINE_ENABLED is true but FACE_ENGINE_URL is empty")
	}

	if len(c.CorsAllowedOrigins) == 0 {
		v.Warnf("CORS_ALLOWED_ORIGINS is empty; cross-origin requests will be rejected")
	}

	return v.Validate(), v.warnings
}

// logEffectiveConfig prints the effective configuration at startup, redacting
// any secret values.
func (c *AppConfig) logEffectiveConfig() {
	slog.Info("effective configuration",
		"photo_dir", c.PhotoDir,
		"api_base_path", c.ApiBasePath,
		"db_host", redactURL(c.DbUrl),
		"db_ssl_mode", "see DbUrl",
		"timezone", c.Timezone.String(),
		"token_duration", c.TokenDuration.String(),
		"admin_username", c.AdminUsername,
		"admin_password", redact(c.AdminPassword),
		"cache_enabled", c.CacheEnabled,
		"ai_enabled", c.AIEnabled,
		"face_engine_enabled", c.FaceEngineEnabled,
		"thumbnail_storage", c.ThumbnailStorage,
		"photo_index_workers", c.PhotoIndexWorkers,
	)
}

// redact masks a secret value, keeping only a short non-sensitive hint.
func redact(v string) string {
	if v == "" {
		return "<empty>"
	}
	if len(v) <= 3 {
		return "***"
	}
	return v[:2] + "***"
}

// redactURL masks the password component of a connection string.
func redactURL(u string) string {
	if idx := strings.Index(u, "password="); idx != -1 {
		rest := u[idx+len("password="):]
		end := strings.IndexAny(rest, " ")
		if end == -1 {
			end = len(rest)
		}
		return u[:idx+len("password=")] + "***" + rest[end:]
	}
	return u
}
