-- +goose Up
-- create index "idx_sessions_hashed_session_token" to table: "sessions"
CREATE UNIQUE INDEX `idx_sessions_hashed_session_token` ON `sessions` (`hashed_session_token`);
-- create index "idx_sessions_hashed_refresh_token" to table: "sessions"
CREATE UNIQUE INDEX `idx_sessions_hashed_refresh_token` ON `sessions` (`hashed_refresh_token`);

-- +goose Down
-- reverse: create index "idx_sessions_hashed_refresh_token" to table: "sessions"
DROP INDEX `idx_sessions_hashed_refresh_token`;
-- reverse: create index "idx_sessions_hashed_session_token" to table: "sessions"
DROP INDEX `idx_sessions_hashed_session_token`;
