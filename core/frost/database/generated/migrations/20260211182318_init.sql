-- +goose Up
-- rename a column from "installer_path" to "download_exe"
ALTER TABLE `local_games` RENAME COLUMN `installer_path` TO `download_exe`;
-- rename a column from "exe_path" to "launch_exe"
ALTER TABLE `local_games` RENAME COLUMN `exe_path` TO `launch_exe`;

-- +goose Down
-- reverse: rename a column from "exe_path" to "launch_exe"
ALTER TABLE `local_games` RENAME COLUMN `launch_exe` TO `exe_path`;
-- reverse: rename a column from "installer_path" to "download_exe"
ALTER TABLE `local_games` RENAME COLUMN `download_exe` TO `installer_path`;
