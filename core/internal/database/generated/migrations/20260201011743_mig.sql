-- +goose Up
-- create "users" table
CREATE TABLE `users` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `username` text NULL,
  `encrypted_password` text NULL,
  `role` text NULL
);
-- create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX `idx_users_username` ON `users` (`username`);
-- create index "idx_users_deleted_at" to table: "users"
CREATE INDEX `idx_users_deleted_at` ON `users` (`deleted_at`);
-- create "sessions" table
CREATE TABLE `sessions` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `user_id` integer NULL,
  `hashed_refresh_token` text NULL,
  `refresh_token_expiry` datetime NULL,
  `hashed_session_token` text NULL,
  `session_token_expiry` datetime NULL,
  `session_type` integer NULL,
  CONSTRAINT `fk_sessions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- create index "idx_sessions_deleted_at" to table: "sessions"
CREATE INDEX `idx_sessions_deleted_at` ON `sessions` (`deleted_at`);

-- +goose Down
-- reverse: create index "idx_sessions_deleted_at" to table: "sessions"
DROP INDEX `idx_sessions_deleted_at`;
-- reverse: create "sessions" table
DROP TABLE `sessions`;
-- reverse: create index "idx_users_deleted_at" to table: "users"
DROP INDEX `idx_users_deleted_at`;
-- reverse: create index "idx_users_username" to table: "users"
DROP INDEX `idx_users_username`;
-- reverse: create "users" table
DROP TABLE `users`;
