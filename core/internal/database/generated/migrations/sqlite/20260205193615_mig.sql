-- +goose Up
-- add column "email" to table: "users"
ALTER TABLE `users` ADD COLUMN `email` text NULL;
-- create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX `idx_users_email` ON `users` (`email`);

-- +goose Down
-- reverse: create index "idx_users_email" to table: "users"
DROP INDEX `idx_users_email`;
-- reverse: add column "email" to table: "users"
ALTER TABLE `users` DROP COLUMN `email`;
