-- +goose Up
-- add column "started" to table: "local_games"
ALTER TABLE `local_games` ADD COLUMN `started` datetime NULL;

-- +goose Down
-- reverse: add column "started" to table: "local_games"
ALTER TABLE `local_games` DROP COLUMN `started`;
