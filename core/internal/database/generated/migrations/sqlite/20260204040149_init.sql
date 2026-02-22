-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_sessions" table
CREATE TABLE `new_sessions` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `user_id` integer NULL,
  `hashed_refresh_token` text NULL,
  `refresh_token_expiry` datetime NULL,
  `hashed_session_token` text NULL,
  `session_token_expiry` datetime NULL,
  `session_type` text NULL,
  CONSTRAINT `fk_sessions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- copy rows from old table "sessions" to new temporary table "new_sessions"
INSERT INTO `new_sessions` (`id`, `created_at`, `updated_at`, `deleted_at`, `user_id`, `hashed_refresh_token`, `refresh_token_expiry`, `hashed_session_token`, `session_token_expiry`, `session_type`) SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `user_id`, `hashed_refresh_token`, `refresh_token_expiry`, `hashed_session_token`, `session_token_expiry`, `session_type` FROM `sessions`;
-- drop "sessions" table after copying rows
DROP TABLE `sessions`;
-- rename temporary table "new_sessions" to "sessions"
ALTER TABLE `new_sessions` RENAME TO `sessions`;
-- create index "idx_sessions_hashed_session_token" to table: "sessions"
CREATE UNIQUE INDEX `idx_sessions_hashed_session_token` ON `sessions` (`hashed_session_token`);
-- create index "idx_sessions_hashed_refresh_token" to table: "sessions"
CREATE UNIQUE INDEX `idx_sessions_hashed_refresh_token` ON `sessions` (`hashed_refresh_token`);
-- create index "idx_sessions_deleted_at" to table: "sessions"
CREATE INDEX `idx_sessions_deleted_at` ON `sessions` (`deleted_at`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "idx_sessions_deleted_at" to table: "sessions"
DROP INDEX `idx_sessions_deleted_at`;
-- reverse: create index "idx_sessions_hashed_refresh_token" to table: "sessions"
DROP INDEX `idx_sessions_hashed_refresh_token`;
-- reverse: create index "idx_sessions_hashed_session_token" to table: "sessions"
DROP INDEX `idx_sessions_hashed_session_token`;
-- reverse: create "new_sessions" table
DROP TABLE `new_sessions`;
