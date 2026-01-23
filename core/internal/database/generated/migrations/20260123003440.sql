-- +goose Up
-- create index "idx_provider_game" to table: "games"
CREATE UNIQUE INDEX `idx_provider_game` ON `games` (`provider_type`, `game_db_id`);

-- +goose Down
-- reverse: create index "idx_provider_game" to table: "games"
DROP INDEX `idx_provider_game`;
