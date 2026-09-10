-- 000002_workspaces down: reverse of the tenancy schema migration. Only ever
-- used for local dev reset (append-only in production).

DROP INDEX IF EXISTS idx_photos_thumb_cap;
DROP INDEX IF EXISTS idx_shared_links_ws;
DROP INDEX IF EXISTS idx_albums_ws_created;
DROP INDEX IF EXISTS idx_photos_ws_album_deleted_epoch;
DROP INDEX IF EXISTS idx_photos_ws_deleted_epoch;
DROP INDEX IF EXISTS idx_workspace_members_user;
DROP INDEX IF EXISTS idx_workspaces_owner;

-- Restore the original global-unique album name index.
DROP INDEX IF EXISTS idx_albums_ws_name;
CREATE UNIQUE INDEX IF NOT EXISTS idx_albums_name ON albums (name);

ALTER TABLE photos DROP COLUMN IF EXISTS thumb_cap;
ALTER TABLE photos DROP COLUMN IF EXISTS workspace_id;
ALTER TABLE albums DROP COLUMN IF EXISTS workspace_id;
ALTER TABLE shared_links DROP COLUMN IF EXISTS workspace_id;

DROP TABLE IF EXISTS workspace_members;
DROP TABLE IF EXISTS workspaces;
DROP INDEX IF EXISTS workspaces_slug;
