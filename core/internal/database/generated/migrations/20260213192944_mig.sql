-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_games" table
CREATE TABLE `new_games` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `provider_type` text NULL,
  `game_db_id` text NULL,
  `url` text NULL,
  `name` text NULL,
  `short_desc` text NULL,
  `full_desc` text NULL,
  `rating` text NULL,
  `rating_count` integer NULL,
  `release_date` datetime NULL,
  `release_status` text NULL,
  `category` text NULL,
  `platforms` text NULL,
  `genres` text NULL,
  `client` text NULL,
  `download_id` text NULL,
  `state` text NULL,
  `complete` integer NULL,
  `left` integer NULL,
  `progress` text NULL,
  `download_url` text NULL,
  `download_path` text NULL,
  `incomplete_path` text NULL,
  `indexer_type` text NULL,
  `game_type` text NULL,
  `title` text NULL,
  `image_url` text NULL,
  `file_size` text NULL,
  `created_iso` text NULL,
  `exe` uniqueIndex NULL
);
-- copy rows from old table "games" to new temporary table "new_games"
INSERT INTO `new_games` (`id`, `created_at`, `updated_at`, `deleted_at`, `provider_type`, `game_db_id`, `url`, `name`, `short_desc`, `full_desc`, `rating`, `rating_count`, `release_date`, `release_status`, `category`, `platforms`, `genres`, `client`, `download_id`, `state`, `complete`, `left`, `progress`, `download_url`, `download_path`, `incomplete_path`, `indexer_type`, `game_type`, `title`, `image_url`, `file_size`, `created_iso`, `exe`) SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `provider_type`, `game_db_id`, `url`, `name`, `short_desc`, `full_desc`, `rating`, `rating_count`, `release_date`, `release_status`, `category`, `platforms`, `genres`, `client`, `download_id`, `state`, `complete`, `left`, `progress`, `download_url`, `download_path`, `incomplete_path`, `indexer_type`, `game_type`, `title`, `image_url`, `file_size`, `created_iso`, `exe` FROM `games`;
-- drop "games" table after copying rows
DROP TABLE `games`;
-- rename temporary table "new_games" to "games"
ALTER TABLE `new_games` RENAME TO `games`;
-- create index "idx_provider_game" to table: "games"
CREATE UNIQUE INDEX `idx_provider_game` ON `games` (`provider_type`, `game_db_id`);
-- create index "idx_games_deleted_at" to table: "games"
CREATE INDEX `idx_games_deleted_at` ON `games` (`deleted_at`);
-- create "assets" table
CREATE TABLE `assets` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `game_id` integer NULL,
  `type` integer NULL,
  `remote_url` text NULL,
  `local_path` text NULL
);
-- create index "idx_assets_type" to table: "assets"
CREATE INDEX `idx_assets_type` ON `assets` (`type`);
-- create index "idx_assets_game_id" to table: "assets"
CREATE INDEX `idx_assets_game_id` ON `assets` (`game_id`);
-- create index "idx_assets_deleted_at" to table: "assets"
CREATE INDEX `idx_assets_deleted_at` ON `assets` (`deleted_at`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "idx_assets_deleted_at" to table: "assets"
DROP INDEX `idx_assets_deleted_at`;
-- reverse: create index "idx_assets_game_id" to table: "assets"
DROP INDEX `idx_assets_game_id`;
-- reverse: create index "idx_assets_type" to table: "assets"
DROP INDEX `idx_assets_type`;
-- reverse: create "assets" table
DROP TABLE `assets`;
-- reverse: create index "idx_games_deleted_at" to table: "games"
DROP INDEX `idx_games_deleted_at`;
-- reverse: create index "idx_provider_game" to table: "games"
DROP INDEX `idx_provider_game`;
-- reverse: create "new_games" table
DROP TABLE `new_games`;
