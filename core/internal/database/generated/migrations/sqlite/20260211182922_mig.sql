-- +goose Up
-- add column "exe" to table: "games"
ALTER TABLE `games` ADD COLUMN `exe` uniqueIndex NULL;

-- +goose Down
-- reverse: add column "exe" to table: "games"
ALTER TABLE `games` DROP COLUMN `exe`;
