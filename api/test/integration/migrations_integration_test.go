//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/r1chjames/photobox/api/internal/adapter/storage/database/migrations"
	"gitlab.com/r1chjames/photobox/api/test/testutil"
)

// TestMigrations_Run_Integration verifies the golang-migrate runner applies
// the embedded baseline and is idempotent on a fresh database (issue #74
// Phase 0). Runs in its own schema (photobox_migrations_test) so it never
// disturbs the shared test schema.
func TestMigrations_Run_Integration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.TeardownTestDB(t, db)

	// Create the isolated schema for this test's migration run.
	err := db.Exec("CREATE SCHEMA IF NOT EXISTS photobox_migrations_test").Error
	require.NoError(t, err)

	dsn := "host=localhost user=photobox_test password=photobox_test dbname=photobox_test port=5433 sslmode=disable search_path=photobox_migrations_test,public"

	// Apply migrations (baseline creates the app tables).
	require.NoError(t, migrations.Run(dsn))

	// Re-running is a clean no-op (not an error, not a dirty state).
	require.NoError(t, migrations.Run(dsn))

	// The baseline created the app tables in the isolated schema.
	var n int
	err = db.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'photobox_migrations_test' AND table_name = 'photos'`).Scan(&n).Error
	require.NoError(t, err)
	assert.Equal(t, 1, n, "photos table should exist after baseline migration")

	// Exactly one applied migration version is recorded.
	err = db.Raw(`SELECT COUNT(*) FROM photobox_migrations_test.schema_migrations`).Scan(&n).Error
	require.NoError(t, err)
	assert.Equal(t, 1, n, "exactly one applied migration version")
}
