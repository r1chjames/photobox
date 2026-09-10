-- Down: reverse of 000001_baseline. The deployment is append-only
-- (migrate down is never run in production). This exists for local dev
-- reset. Names are unqualified like the up migration — the session
-- search_path determines the schema. Drops tables in dependency-safe
-- reverse order.

DROP TABLE IF EXISTS face_clusters;
DROP TABLE IF EXISTS face_persons;
DROP TABLE IF EXISTS face_detections;
DROP TABLE IF EXISTS photo_tags;
DROP TABLE IF EXISTS photo_analysis;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS shared_links;
DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS photos;
DROP TABLE IF EXISTS albums;
