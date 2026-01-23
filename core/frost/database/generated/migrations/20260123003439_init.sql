-- +goose Up
-- create index "idx_provider_game" to table: "local_games"
CREATE UNIQUE INDEX `idx_provider_game` ON `local_games` (`provider_type`, `game_db_id`);

-- +goose Down
-- reverse: create index "idx_provider_game" to table: "local_games"
DROP INDEX `idx_provider_game`;
