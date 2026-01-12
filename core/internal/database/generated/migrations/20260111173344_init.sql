-- +goose Up
-- create "games" table
CREATE TABLE `games` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `downloaded_path` text NULL
);
-- create index "idx_games_deleted_at" to table: "games"
CREATE INDEX `idx_games_deleted_at` ON `games` (`deleted_at`);

-- +goose Down
-- reverse: create index "idx_games_deleted_at" to table: "games"
DROP INDEX `idx_games_deleted_at`;
-- reverse: create "games" table
DROP TABLE `games`;
