// Package migrations embeds and applies versioned SQL schema migrations
// (golang-migrate). Phase 0 of issue #74: existing databases created by GORM
// AutoMigrate are adopted via the idempotent 000001_baseline migration, and
// every subsequent schema change ships as a numbered .sql migration here.
package migrations

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq" // postgres driver used by golang-migrate
)

//go:embed *.sql
var migrationsFS embed.FS

// Run applies all pending migrations to the database described by dbUrl, a
// lib/pq-style key=value DSN matching appconfig.AppConfig.DbUrl (e.g.
// "host=... user=... dbname=... sslmode=... search_path=photobox,public").
//
// golang-migrate's postgres driver serializes concurrent replicas with a
// session-level advisory lock and records applied versions in the
// schema_migrations table (created in the first schema of search_path).
func Run(dbUrl string) error {
	migrateURL, err := toMigrateURL(dbUrl)
	if err != nil {
		return err
	}

	src, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("migrations: embed source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, migrateURL)
	if err != nil {
		return fmt.Errorf("migrations: init: %w", err)
	}
	m.Log = &migrateLogger{}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("migrations: apply: %w", err)
	}
	return nil
}

// toMigrateURL converts a lib/pq key=value DSN into the postgres:// URL form
// golang-migrate requires, preserving search_path so migrations and the
// schema_migrations table land in the photobox schema.
func toMigrateURL(dsn string) (string, error) {
	params := map[string]string{}
	for _, kv := range strings.Fields(dsn) {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		}
	}

	host := params["host"]
	if host == "" {
		host = "localhost"
	}
	port := params["port"]
	if port == "" {
		port = "5432"
	}
	user := url.QueryEscape(params["user"])
	password := url.QueryEscape(params["password"])
	dbname := params["dbname"]
	if dbname == "" {
		return "", fmt.Errorf("migrations: dbname missing from DSN")
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
		Path:   "/" + dbname,
	}
	q := url.Values{}
	if sslmode := params["sslmode"]; sslmode != "" {
		q.Set("sslmode", sslmode)
	}
	if searchPath := params["search_path"]; searchPath != "" {
		// migrate parses x- queries out before opening the connection, but the
		// driver keeps other query params in the DSN handed to lib/pq, so
		// search_path survives as a runtime connection option.
		q.Set("search_path", searchPath)
	}
	// Multi-statement .sql files require explicit opt-in.
	q.Set("x-multi-statement", "true")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

type migrateLogger struct{}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	slog.Info(fmt.Sprintf("migrate: "+format, v...))
}

func (l *migrateLogger) Verbose() bool { return false }
