-- +goose Up
CREATE TABLE IF NOT EXISTS files (
  id TEXT PRIMARY KEY,
  path TEXT NOT NULL,
  drive_id TEXT NOT NULL,
  etag TEXT,
  checksum TEXT,
  size INTEGER NOT NULL DEFAULT 0,
  modified_at INTEGER NOT NULL DEFAULT 0,
  created_at INTEGER NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_files_path ON files(path);

-- +goose Down
DROP INDEX IF EXISTS idx_files_path;
DROP TABLE IF EXISTS files;
