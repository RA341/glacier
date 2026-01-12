-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_games" table
CREATE TABLE `new_games` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `name` text NULL,
  `game_type` integer NULL,
  `downloaded_path` text NULL
);
-- copy rows from old table "games" to new temporary table "new_games"
INSERT INTO `new_games` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `game_type`, `downloaded_path`) SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `name`, `game_type`, `downloaded_path` FROM `games`;
-- drop "games" table after copying rows
DROP TABLE `games`;
-- rename temporary table "new_games" to "games"
ALTER TABLE `new_games` RENAME TO `games`;
-- create index "idx_games_deleted_at" to table: "games"
CREATE INDEX `idx_games_deleted_at` ON `games` (`deleted_at`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "idx_games_deleted_at" to table: "games"
DROP INDEX `idx_games_deleted_at`;
-- reverse: create "new_games" table
DROP TABLE `new_games`;
