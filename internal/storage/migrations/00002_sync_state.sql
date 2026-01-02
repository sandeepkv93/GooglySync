-- +goose Up
ALTER TABLE files ADD COLUMN account_id TEXT NOT NULL DEFAULT '';
DROP INDEX IF EXISTS idx_files_path;
CREATE UNIQUE INDEX IF NOT EXISTS idx_files_account_path ON files(account_id, path);
CREATE INDEX IF NOT EXISTS idx_files_account_drive_id ON files(account_id, drive_id);

CREATE TABLE IF NOT EXISTS accounts (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL,
  display_name TEXT NOT NULL DEFAULT '',
  is_primary INTEGER NOT NULL DEFAULT 0,
  created_at INTEGER NOT NULL DEFAULT 0,
  updated_at INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);

CREATE TABLE IF NOT EXISTS token_refs (
  account_id TEXT PRIMARY KEY,
  key_id TEXT NOT NULL,
  token_type TEXT NOT NULL DEFAULT '',
  scope TEXT NOT NULL DEFAULT '',
  expiry INTEGER NOT NULL DEFAULT 0,
  updated_at INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sync_state (
  account_id TEXT PRIMARY KEY,
  start_page_token TEXT NOT NULL DEFAULT '',
  last_sync_at INTEGER NOT NULL DEFAULT 0,
  last_error TEXT NOT NULL DEFAULT '',
  paused INTEGER NOT NULL DEFAULT 0,
  updated_at INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pending_ops (
  id TEXT PRIMARY KEY,
  account_id TEXT NOT NULL,
  path TEXT NOT NULL,
  drive_id TEXT NOT NULL DEFAULT '',
  op_type TEXT NOT NULL,
  state TEXT NOT NULL DEFAULT 'queued',
  retry_count INTEGER NOT NULL DEFAULT 0,
  last_error TEXT NOT NULL DEFAULT '',
  created_at INTEGER NOT NULL DEFAULT 0,
  updated_at INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_pending_ops_account_state ON pending_ops(account_id, state, created_at);

CREATE TABLE IF NOT EXISTS folders (
  id TEXT PRIMARY KEY,
  account_id TEXT NOT NULL,
  path TEXT NOT NULL,
  drive_id TEXT NOT NULL,
  parent_id TEXT NOT NULL DEFAULT '',
  modified_at INTEGER NOT NULL DEFAULT 0,
  created_at INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY(account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_folders_account_path ON folders(account_id, path);

CREATE TABLE IF NOT EXISTS shared_drives (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  created_at INTEGER NOT NULL DEFAULT 0,
  updated_at INTEGER NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS pending_ops;
DROP TABLE IF EXISTS sync_state;
DROP TABLE IF EXISTS token_refs;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS shared_drives;
DROP TABLE IF EXISTS accounts;

DROP INDEX IF EXISTS idx_files_account_drive_id;
DROP INDEX IF EXISTS idx_files_account_path;

CREATE TABLE IF NOT EXISTS files_old (
  id TEXT PRIMARY KEY,
  path TEXT NOT NULL,
  drive_id TEXT NOT NULL,
  etag TEXT,
  checksum TEXT,
  size INTEGER NOT NULL DEFAULT 0,
  modified_at INTEGER NOT NULL DEFAULT 0,
  created_at INTEGER NOT NULL DEFAULT 0
);
INSERT INTO files_old (id, path, drive_id, etag, checksum, size, modified_at, created_at)
  SELECT id, path, drive_id, etag, checksum, size, modified_at, created_at FROM files;
DROP TABLE files;
ALTER TABLE files_old RENAME TO files;
CREATE UNIQUE INDEX IF NOT EXISTS idx_files_path ON files(path);
