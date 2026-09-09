-- Baseline: existing AutoMigrate-produced schema (issue #74, Phase 0).
--
-- This migration is a no-op on databases created fresh by AutoMigrate and
-- version-marks existing databases that predate golang-migrate. Every
-- statement is idempotent so it is safe to run against both an empty
-- database and a populated AutoMigrate database. Subsequent migrations
-- (000002_*) apply the workspace/tenancy changes on top.
--
-- Names are intentionally unqualified: the runner session's search_path
-- (from the connection DSN) places them in the photobox schema in
-- production and in an isolated schema under test.

-- Albums
CREATE TABLE IF NOT EXISTS albums (
	id text NOT NULL,
	name text,
	description text,
	tags text,
	metadata jsonb,
	created_at timestamptz,
	created_epoch bigint,
	updated_at timestamptz,
	thumbnail text,
	cover_photo_id text,
	PRIMARY KEY (id)
);

-- Photos
CREATE TABLE IF NOT EXISTS photos (
	id text NOT NULL,
	name text,
	filesystem_path text,
	source_path text,
	album_id text,
	metadata jsonb,
	created_at timestamptz,
	created_epoch bigint,
	year bigint,
	month bigint,
	updated_at timestamptz,
	thumbnail bytea,
	thumbnail_path text,
	favorite boolean DEFAULT false,
	deleted_at timestamptz,
	blurhash text,
	dominant_color text,
	file_hash text,
	file_modified_time bigint,
	media_type text DEFAULT 'image',
	duration bigint,
	width bigint,
	height bigint,
	latitude numeric,
	longitude numeric,
	hidden boolean DEFAULT false,
	live_photo_path text,
	description text,
	quality_score bigint DEFAULT 0,
	blur_score numeric DEFAULT 0,
	is_low_quality boolean DEFAULT false,
	trash_path text,
	edit_params jsonb,
	PRIMARY KEY (id)
);

-- Users
CREATE TABLE IF NOT EXISTS users (
	id text NOT NULL,
	username text,
	role text,
	password text,
	email text,
	approved boolean,
	created_at timestamptz,
	updated_at timestamptz,
	PRIMARY KEY (id)
);

-- Jobs
CREATE TABLE IF NOT EXISTS jobs (
	name text NOT NULL,
	status text,
	last_run timestamptz,
	created_at timestamptz,
	updated_at timestamptz,
	PRIMARY KEY (name)
);

-- Settings
CREATE TABLE IF NOT EXISTS settings (
	key text NOT NULL,
	value text,
	friendly_name text,
	category text,
	type text,
	options text,
	description text,
	created_at timestamptz,
	updated_at timestamptz,
	PRIMARY KEY (key)
);

-- Shared links
CREATE TABLE IF NOT EXISTS shared_links (
	token varchar(64) NOT NULL,
	resource_type varchar(20) NOT NULL,
	resource_id varchar(255) NOT NULL,
	created_by varchar(36),
	expiry timestamptz,
	password_hash varchar(255),
	view_count bigint DEFAULT 0,
	created_at timestamptz,
	updated_at timestamptz,
	PRIMARY KEY (token)
);

-- API keys
CREATE TABLE IF NOT EXISTS api_keys (
	id text NOT NULL,
	name varchar(100),
	key_hash varchar(255),
	scope varchar(20),
	created_by varchar(36),
	created_at timestamptz,
	last_used timestamptz,
	revoked boolean DEFAULT false,
	PRIMARY KEY (id)
);

-- Photo analysis
CREATE TABLE IF NOT EXISTS photo_analysis (
	photo_id text NOT NULL,
	model text,
	caption text,
	tags jsonb,
	objects jsonb,
	is_nsfw boolean,
	is_portrait boolean,
	status text DEFAULT 'pending',
	attempts bigint DEFAULT 0,
	error_message text,
	retryable boolean DEFAULT true,
	started_at timestamptz,
	last_attempt_at timestamptz,
	created_at timestamptz,
	updated_at timestamptz,
	PRIMARY KEY (photo_id)
);

-- Photo tags
CREATE TABLE IF NOT EXISTS photo_tags (
	photo_id text NOT NULL,
	tag text NOT NULL,
	source text DEFAULT '',
	PRIMARY KEY (photo_id, tag)
);

-- Face recognition (issue #80)
CREATE TABLE IF NOT EXISTS face_detections (
	id text NOT NULL,
	photo_id text,
	person_id text,
	box_x numeric,
	box_y numeric,
	box_w numeric,
	box_h numeric,
	score numeric,
	embedding bytea,
	status text DEFAULT 'detected',
	quality numeric DEFAULT 0,
	created_at timestamptz,
	updated_at timestamptz,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS face_persons (
	id text NOT NULL,
	name text,
	cover_face_id text,
	created_at timestamptz,
	updated_at timestamptz,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS face_clusters (
	id text NOT NULL,
	person_id text,
	centroid bytea,
	face_count bigint DEFAULT 0,
	needs_review boolean DEFAULT false,
	updated_at timestamptz,
	PRIMARY KEY (id)
);

-- Indexes (idempotent — match AutoMigrate plus explicit setup output)
CREATE UNIQUE INDEX IF NOT EXISTS albums_pkey ON albums (id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_albums_name ON albums (name);
CREATE INDEX IF NOT EXISTS idx_albums_created_epoch ON albums (created_epoch);
CREATE INDEX IF NOT EXISTS idx_album_search ON albums USING gin (to_tsvector('english'::regconfig, coalesce(name, ''::text)));

CREATE UNIQUE INDEX IF NOT EXISTS photos_pkey ON photos (id);
CREATE INDEX IF NOT EXISTS idx_photos_created_epoch ON photos (created_epoch);
CREATE INDEX IF NOT EXISTS idx_photos_deleted_at ON photos (deleted_at);
CREATE INDEX IF NOT EXISTS idx_photos_favorite ON photos (favorite);
CREATE INDEX IF NOT EXISTS idx_photos_file_hash ON photos (file_hash);
CREATE INDEX IF NOT EXISTS idx_photos_file_modified_time ON photos (file_modified_time);
CREATE INDEX IF NOT EXISTS idx_photos_hidden ON photos (hidden);
CREATE INDEX IF NOT EXISTS idx_photos_is_low_quality ON photos (is_low_quality);
CREATE INDEX IF NOT EXISTS idx_photos_latitude ON photos (latitude);
CREATE INDEX IF NOT EXISTS idx_photos_longitude ON photos (longitude);
CREATE INDEX IF NOT EXISTS idx_photos_live_photo_path ON photos (live_photo_path);
CREATE INDEX IF NOT EXISTS idx_photos_month ON photos (month);
CREATE INDEX IF NOT EXISTS idx_photos_year ON photos (year);
CREATE INDEX IF NOT EXISTS idx_deleted_epoch ON photos (deleted_at, created_epoch);
CREATE INDEX IF NOT EXISTS idx_album_deleted_epoch ON photos (album_id, deleted_at, created_epoch);
CREATE INDEX IF NOT EXISTS idx_photo_search ON photos USING gin (to_tsvector('english'::regconfig, coalesce(name, ''::text)));
CREATE INDEX IF NOT EXISTS idx_photos_lat_lng ON photos (latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS users_pkey ON users (id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username);

CREATE UNIQUE INDEX IF NOT EXISTS jobs_pkey ON jobs (name);

CREATE UNIQUE INDEX IF NOT EXISTS settings_pkey ON settings (key);

CREATE UNIQUE INDEX IF NOT EXISTS shared_links_pkey ON shared_links (token);

CREATE UNIQUE INDEX IF NOT EXISTS api_keys_pkey ON api_keys (id);

CREATE UNIQUE INDEX IF NOT EXISTS photo_analysis_pkey ON photo_analysis (photo_id);

CREATE INDEX IF NOT EXISTS idx_photo_tag ON photo_tags (photo_id, tag);
CREATE INDEX IF NOT EXISTS idx_tag_photo ON photo_tags (tag);

CREATE UNIQUE INDEX IF NOT EXISTS face_detections_pkey ON face_detections (id);
CREATE INDEX IF NOT EXISTS idx_face_detections_photo ON face_detections (photo_id);
CREATE INDEX IF NOT EXISTS idx_face_detections_person ON face_detections (person_id);
CREATE INDEX IF NOT EXISTS idx_face_detections_status ON face_detections (status);

CREATE UNIQUE INDEX IF NOT EXISTS face_persons_pkey ON face_persons (id);
CREATE INDEX IF NOT EXISTS idx_face_persons_name ON face_persons (name);

CREATE UNIQUE INDEX IF NOT EXISTS face_clusters_pkey ON face_clusters (id);
CREATE INDEX IF NOT EXISTS idx_face_clusters_person ON face_clusters (person_id);
