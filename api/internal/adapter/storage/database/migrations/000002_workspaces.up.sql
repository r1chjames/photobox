-- 000002_workspaces: tenancy schema (issue #74, Phase 1).
--
-- Adds the isolation dimension: workspaces, workspace_members, and
-- workspace_id/thumb_cap columns on content tables, then backfills all
-- pre-existing content and users into the single deterministic default
-- workspace so today's shared-library semantics are preserved exactly (D6).
--
-- Naming: unqualified, matching 000001_baseline — the session search_path
-- (from the DSN) selects the schema. Idempotency: this migration applies
-- exactly once via golang-migrate statements are written defensively.

-- Workspaces + memberships
CREATE TABLE IF NOT EXISTS workspaces (
	id text PRIMARY KEY,
	name text NOT NULL DEFAULT '',
	slug text NOT NULL DEFAULT '',
	owner_user_id text,
	storage_used_bytes bigint NOT NULL DEFAULT 0,
	storage_limit_bytes bigint NOT NULL DEFAULT 0,
	settings jsonb,
	created_at timestamptz,
	updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS workspace_members (
	workspace_id text NOT NULL,
	user_id text NOT NULL,
	role text NOT NULL DEFAULT 'viewer',
	created_at timestamptz,
	updated_at timestamptz,
	PRIMARY KEY (workspace_id, user_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS workspaces_slug ON workspaces (slug);
CREATE INDEX IF NOT EXISTS idx_workspaces_owner ON workspaces (owner_user_id);
CREATE INDEX IF NOT EXISTS idx_workspace_members_user ON workspace_members (user_id);

-- The deterministic default workspace adopting pre-tenancy data. Its UUID is
-- fixed (domain.DefaultWorkspaceID) so the backfill below is repeatable and
-- future code can refer to it.
INSERT INTO workspaces (id, name, slug, owner_user_id, created_at, updated_at)
SELECT '00000000-0000-0000-0000-000000000001', 'Default', 'default', MIN(id), NOW(), NOW()
FROM users
ON CONFLICT (id) DO NOTHING;

-- Every existing user becomes a member (owner) of the default workspace.
INSERT INTO workspace_members (workspace_id, user_id, role, created_at, updated_at)
SELECT '00000000-0000-0000-0000-000000000001', id, 'owner', NOW(), NOW()
FROM users
ON CONFLICT (workspace_id, user_id) DO NOTHING;

-- Content tables: nullable workspace_id + backfill. Columns are dropped from
-- the NOT NULL set (nullable through Phase 2 the 000003 migration enforces
-- NOT NULL once every write path carries a workspace).

-- Photos: workspace_id + thumb_cap (random 128-bit capability immutable).
ALTER TABLE photos ADD COLUMN IF NOT EXISTS workspace_id text;
ALTER TABLE photos ADD COLUMN IF NOT EXISTS thumb_cap text;

UPDATE photos SET workspace_id = '00000000-0000-0000-0000-000000000001'
WHERE workspace_id IS NULL;

UPDATE photos SET thumb_cap = gen_random_uuid()::text
WHERE thumb_cap IS NULL OR thumb_cap = '';

-- Albums: workspace_id drop the global unique index on name in favor of a
-- per-workspace unique index (two tenants may both have "Holidays").
ALTER TABLE albums ADD COLUMN IF NOT EXISTS workspace_id text;

UPDATE albums SET workspace_id = '00000000-0000-0000-0000-000000000001'
WHERE workspace_id IS NULL;

DROP INDEX IF EXISTS idx_albums_name;
CREATE UNIQUE INDEX IF NOT EXISTS idx_albums_ws_name ON albums (workspace_id, name);

-- Shared links: workspace_id.
ALTER TABLE shared_links ADD COLUMN IF NOT EXISTS workspace_id text;

UPDATE shared_links SET workspace_id = '00000000-0000-0000-0000-000000000001'
WHERE workspace_id IS NULL;

-- Composite workspace-aware indexes on the hot query paths.
CREATE INDEX IF NOT EXISTS idx_photos_ws_deleted_epoch ON photos (workspace_id, deleted_at, created_epoch);
CREATE INDEX IF NOT EXISTS idx_photos_ws_album_deleted_epoch ON photos (workspace_id, album_id, deleted_at, created_epoch);
CREATE INDEX IF NOT EXISTS idx_albums_ws_created ON albums (workspace_id, created_epoch);
CREATE INDEX IF NOT EXISTS idx_shared_links_ws ON shared_links (workspace_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_photos_thumb_cap ON photos (thumb_cap) WHERE thumb_cap IS NOT NULL AND thumb_cap <> '';
