-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_folder_manifests" table
CREATE TABLE `new_folder_manifests` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `game_id` integer NULL,
  `total_size` integer NULL,
  `file_info` text NULL,
  CONSTRAINT `fk_folder_manifests_game` FOREIGN KEY (`game_id`) REFERENCES `games` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- copy rows from old table "folder_manifests" to new temporary table "new_folder_manifests"
INSERT INTO `new_folder_manifests` (`id`, `created_at`, `updated_at`, `deleted_at`, `game_id`, `total_size`, `file_info`) SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `game_id`, `total_size`, `file_info` FROM `folder_manifests`;
-- drop "folder_manifests" table after copying rows
DROP TABLE `folder_manifests`;
-- rename temporary table "new_folder_manifests" to "folder_manifests"
ALTER TABLE `new_folder_manifests` RENAME TO `folder_manifests`;
-- create index "idx_folder_manifests_game_id" to table: "folder_manifests"
CREATE UNIQUE INDEX `idx_folder_manifests_game_id` ON `folder_manifests` (`game_id`);
-- create index "idx_folder_manifests_deleted_at" to table: "folder_manifests"
CREATE INDEX `idx_folder_manifests_deleted_at` ON `folder_manifests` (`deleted_at`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "idx_folder_manifests_deleted_at" to table: "folder_manifests"
DROP INDEX `idx_folder_manifests_deleted_at`;
-- reverse: create index "idx_folder_manifests_game_id" to table: "folder_manifests"
DROP INDEX `idx_folder_manifests_game_id`;
-- reverse: create "new_folder_manifests" table
DROP TABLE `new_folder_manifests`;
