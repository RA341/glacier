-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_service_configs" table
CREATE TABLE `new_service_configs` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `service_type` integer NULL,
  `name` text NULL,
  `enabled` numeric NULL,
  `flavour` text NULL,
  `config` text NULL
);
-- copy rows from old table "service_configs" to new temporary table "new_service_configs"
INSERT INTO `new_service_configs` (`id`, `created_at`, `updated_at`, `deleted_at`, `service_type`, `name`, `enabled`, `flavour`, `config`) SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `service_type`, `name`, `enabled`, `flavour`, `config` FROM `service_configs`;
-- drop "service_configs" table after copying rows
DROP TABLE `service_configs`;
-- rename temporary table "new_service_configs" to "service_configs"
ALTER TABLE `new_service_configs` RENAME TO `service_configs`;
-- create index "idx_service_config" to table: "service_configs"
CREATE UNIQUE INDEX `idx_service_config` ON `service_configs` (`service_type`, `name`);
-- create index "idx_service_configs_deleted_at" to table: "service_configs"
CREATE INDEX `idx_service_configs_deleted_at` ON `service_configs` (`deleted_at`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "idx_service_configs_deleted_at" to table: "service_configs"
DROP INDEX `idx_service_configs_deleted_at`;
-- reverse: create index "idx_service_config" to table: "service_configs"
DROP INDEX `idx_service_config`;
-- reverse: create "new_service_configs" table
DROP TABLE `new_service_configs`;
