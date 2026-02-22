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
  `short_desc` text NULL,
  `full_desc` text NULL,
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
  `created_iso` text NULL
);
-- create index "idx_provider_game" to table: "games"
CREATE UNIQUE INDEX `idx_provider_game` ON `games` (`provider_type`, `game_db_id`);
-- create index "idx_games_deleted_at" to table: "games"
CREATE INDEX `idx_games_deleted_at` ON `games` (`deleted_at`);
-- create "folder_manifests" table
CREATE TABLE `folder_manifests` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `game_id` integer NULL,
  `total_size` integer NULL,
  `available_exe_paths` text NULL,
  `file_info` text NULL,
  CONSTRAINT `fk_folder_manifests_game` FOREIGN KEY (`game_id`) REFERENCES `games` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- create index "idx_folder_manifests_game_id" to table: "folder_manifests"
CREATE INDEX `idx_folder_manifests_game_id` ON `folder_manifests` (`game_id`);
-- create index "idx_folder_manifests_deleted_at" to table: "folder_manifests"
CREATE INDEX `idx_folder_manifests_deleted_at` ON `folder_manifests` (`deleted_at`);
-- create "service_configs" table
CREATE TABLE `service_configs` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `service_type` text NULL,
  `name` text NULL,
  `enabled` numeric NULL,
  `flavour` text NULL,
  `config` text NULL
);
-- create index "idx_service_config" to table: "service_configs"
CREATE UNIQUE INDEX `idx_service_config` ON `service_configs` (`service_type`, `name`);
-- create index "idx_service_configs_deleted_at" to table: "service_configs"
CREATE INDEX `idx_service_configs_deleted_at` ON `service_configs` (`deleted_at`);

-- +goose Down
-- reverse: create index "idx_service_configs_deleted_at" to table: "service_configs"
DROP INDEX `idx_service_configs_deleted_at`;
-- reverse: create index "idx_service_config" to table: "service_configs"
DROP INDEX `idx_service_config`;
-- reverse: create "service_configs" table
DROP TABLE `service_configs`;
-- reverse: create index "idx_folder_manifests_deleted_at" to table: "folder_manifests"
DROP INDEX `idx_folder_manifests_deleted_at`;
-- reverse: create index "idx_folder_manifests_game_id" to table: "folder_manifests"
DROP INDEX `idx_folder_manifests_game_id`;
-- reverse: create "folder_manifests" table
DROP TABLE `folder_manifests`;
-- reverse: create index "idx_games_deleted_at" to table: "games"
DROP INDEX `idx_games_deleted_at`;
-- reverse: create index "idx_provider_game" to table: "games"
DROP INDEX `idx_provider_game`;
-- reverse: create "games" table
DROP TABLE `games`;
