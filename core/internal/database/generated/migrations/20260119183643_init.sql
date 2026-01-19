-- +goose Up
-- create "games" table
CREATE TABLE `games` (
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
  `client` text NULL,
  `download_id` text NULL,
  `state` text NULL,
  `progress` text NULL,
  `download_url` text NULL,
  `download_path` text NULL,
  `incomplete_path` text NULL,
  `indexer_type` text NULL,
  `game_type` text NULL,
  `title` text NULL,
  `image_url` text NULL,
  `file_size` text NULL,
  `created_iso` text NULL
);
-- create index "idx_games_deleted_at" to table: "games"
CREATE INDEX `idx_games_deleted_at` ON `games` (`deleted_at`);

-- +goose Down
-- reverse: create index "idx_games_deleted_at" to table: "games"
DROP INDEX `idx_games_deleted_at`;
-- reverse: create "games" table
DROP TABLE `games`;
