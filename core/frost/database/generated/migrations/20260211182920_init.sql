-- +goose Up
-- add column "exe" to table: "local_games"
ALTER TABLE `local_games` ADD COLUMN `exe` uniqueIndex NULL;

-- +goose Down
-- reverse: add column "exe" to table: "local_games"
ALTER TABLE `local_games` DROP COLUMN `exe`;
