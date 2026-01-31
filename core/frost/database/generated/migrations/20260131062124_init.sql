-- +goose Up
-- add column "done" to table: "local_games"
ALTER TABLE `local_games` ADD COLUMN `done` datetime NULL;

-- +goose Down
-- reverse: add column "done" to table: "local_games"
ALTER TABLE `local_games` DROP COLUMN `done`;
