-- +goose Up
-- add column "name" to table: "games"
ALTER TABLE `games` ADD COLUMN `name` text NULL;
-- add column "metadata_id" to table: "games"
ALTER TABLE `games` ADD COLUMN `metadata_id` text NULL;
-- add column "game_type" to table: "games"
ALTER TABLE `games` ADD COLUMN `game_type` integer NULL;

-- +goose Down
-- reverse: add column "game_type" to table: "games"
ALTER TABLE `games` DROP COLUMN `game_type`;
-- reverse: add column "metadata_id" to table: "games"
ALTER TABLE `games` DROP COLUMN `metadata_id`;
-- reverse: add column "name" to table: "games"
ALTER TABLE `games` DROP COLUMN `name`;
