-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_local_games" table
CREATE TABLE `new_local_games` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `game_id` integer NULL,
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
  `exe` uniqueIndex NULL,
  `status` integer NULL,
  `status_message` text NULL,
  `started` datetime NULL,
  `done` datetime NULL,
  `launch_exe` text NULL
);
-- copy rows from old table "local_games" to new temporary table "new_local_games"
INSERT INTO `new_local_games` (`id`, `created_at`, `updated_at`, `deleted_at`, `game_id`, `provider_type`, `game_db_id`, `url`, `name`, `short_desc`, `full_desc`, `rating`, `rating_count`, `release_date`, `release_status`, `category`, `platforms`, `genres`, `client`, `download_id`, `state`, `complete`, `left`, `progress`, `download_url`, `download_path`, `incomplete_path`, `indexer_type`, `game_type`, `title`, `image_url`, `file_size`, `created_iso`, `exe`, `status`, `status_message`, `started`, `done`, `launch_exe`) SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `game_id`, `provider_type`, `game_db_id`, `url`, `name`, `short_desc`, `full_desc`, `rating`, `rating_count`, `release_date`, `release_status`, `category`, `platforms`, `genres`, `client`, `download_id`, `state`, `complete`, `left`, `progress`, `download_url`, `download_path`, `incomplete_path`, `indexer_type`, `game_type`, `title`, `image_url`, `file_size`, `created_iso`, `exe`, `status`, `status_message`, `started`, `done`, `launch_exe` FROM `local_games`;
-- drop "local_games" table after copying rows
DROP TABLE `local_games`;
-- rename temporary table "new_local_games" to "local_games"
ALTER TABLE `new_local_games` RENAME TO `local_games`;
-- create index "idx_provider_game" to table: "local_games"
CREATE UNIQUE INDEX `idx_provider_game` ON `local_games` (`provider_type`, `game_db_id`);
-- create index "idx_local_games_deleted_at" to table: "local_games"
CREATE INDEX `idx_local_games_deleted_at` ON `local_games` (`deleted_at`, `deleted_at`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "idx_local_games_deleted_at" to table: "local_games"
DROP INDEX `idx_local_games_deleted_at`;
-- reverse: create index "idx_provider_game" to table: "local_games"
DROP INDEX `idx_provider_game`;
-- reverse: create "new_local_games" table
DROP TABLE `new_local_games`;
