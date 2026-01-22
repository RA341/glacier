-- +goose Up
-- create "local_games" table
CREATE TABLE `local_games` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `provider_type` text NULL,
  `game_db_id` text NULL,
  `name` text NULL,
  `summary` text NULL,
  `description` text NULL,
  `url` text NULL,
  `thumbnail_url` text NULL,
  `videos` text NULL,
  `platforms` text NULL,
  `genres` text NULL,
  `rating` text NULL,
  `rating_count` integer NULL,
  `release_date` datetime NULL,
  `release_status` text NULL,
  `category` text NULL,
  `download_path` text NULL,
  `installer_path` text NULL,
  `exe_path` text NULL,
  `status` integer NULL,
  `status_message` text NULL
);
-- create index "idx_local_games_deleted_at" to table: "local_games"
CREATE INDEX `idx_local_games_deleted_at` ON `local_games` (`deleted_at`);

-- +goose Down
-- reverse: create index "idx_local_games_deleted_at" to table: "local_games"
DROP INDEX `idx_local_games_deleted_at`;
-- reverse: create "local_games" table
DROP TABLE `local_games`;
