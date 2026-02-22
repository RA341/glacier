-- +goose Up
-- create "assets" table
CREATE TABLE "public"."assets" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "game_id" bigint NULL,
  "type" bigint NULL,
  "remote_url" text NULL,
  "local_path" text NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_assets_deleted_at" to table: "assets"
CREATE INDEX "idx_assets_deleted_at" ON "public"."assets" ("deleted_at");
-- create index "idx_assets_game_id" to table: "assets"
CREATE INDEX "idx_assets_game_id" ON "public"."assets" ("game_id");
-- create index "idx_assets_type" to table: "assets"
CREATE INDEX "idx_assets_type" ON "public"."assets" ("type");
-- create "games" table
CREATE TABLE "public"."games" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "provider_type" text NULL,
  "game_db_id" text NULL,
  "url" text NULL,
  "name" text NULL,
  "short_desc" text NULL,
  "full_desc" text NULL,
  "rating" text NULL,
  "rating_count" bigint NULL,
  "release_date" timestamptz NULL,
  "release_status" text NULL,
  "category" text NULL,
  "platforms" text NULL,
  "genres" text NULL,
  "client" text NULL,
  "download_id" text NULL,
  "state" text NULL,
  "complete" bigint NULL,
  "left" bigint NULL,
  "progress" text NULL,
  "download_url" text NULL,
  "download_path" text NULL,
  "incomplete_path" text NULL,
  "indexer_type" text NULL,
  "game_type" text NULL,
  "title" text NULL,
  "image_url" text NULL,
  "file_size" text NULL,
  "created_iso" text NULL,
  "exe" text NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_games_deleted_at" to table: "games"
CREATE INDEX "idx_games_deleted_at" ON "public"."games" ("deleted_at");
-- create index "idx_games_exe" to table: "games"
CREATE UNIQUE INDEX "idx_games_exe" ON "public"."games" ("exe");
-- create index "idx_provider_game" to table: "games"
CREATE UNIQUE INDEX "idx_provider_game" ON "public"."games" ("provider_type", "game_db_id");
-- create "service_configs" table
CREATE TABLE "public"."service_configs" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "service_type" text NULL,
  "name" text NULL,
  "enabled" boolean NULL,
  "flavour" text NULL,
  "config" text NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_service_config" to table: "service_configs"
CREATE UNIQUE INDEX "idx_service_config" ON "public"."service_configs" ("service_type", "name");
-- create index "idx_service_configs_deleted_at" to table: "service_configs"
CREATE INDEX "idx_service_configs_deleted_at" ON "public"."service_configs" ("deleted_at");
-- create "folder_manifests" table
CREATE TABLE "public"."folder_manifests" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "game_id" bigint NULL,
  "total_size" bigint NULL,
  "file_info" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_games_manifest" FOREIGN KEY ("game_id") REFERENCES "public"."games" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- create index "idx_folder_manifests_deleted_at" to table: "folder_manifests"
CREATE INDEX "idx_folder_manifests_deleted_at" ON "public"."folder_manifests" ("deleted_at");
-- create index "idx_folder_manifests_game_id" to table: "folder_manifests"
CREATE UNIQUE INDEX "idx_folder_manifests_game_id" ON "public"."folder_manifests" ("game_id");
-- create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "username" text NULL,
  "email" text NULL,
  "encrypted_password" text NULL,
  "role" text NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" ("username");
-- create "sessions" table
CREATE TABLE "public"."sessions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NULL,
  "hashed_refresh_token" text NULL,
  "refresh_token_expiry" timestamptz NULL,
  "hashed_session_token" text NULL,
  "session_token_expiry" timestamptz NULL,
  "session_type" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_sessions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- create index "idx_sessions_deleted_at" to table: "sessions"
CREATE INDEX "idx_sessions_deleted_at" ON "public"."sessions" ("deleted_at");
-- create index "idx_sessions_hashed_refresh_token" to table: "sessions"
CREATE UNIQUE INDEX "idx_sessions_hashed_refresh_token" ON "public"."sessions" ("hashed_refresh_token");
-- create index "idx_sessions_hashed_session_token" to table: "sessions"
CREATE UNIQUE INDEX "idx_sessions_hashed_session_token" ON "public"."sessions" ("hashed_session_token");

-- +goose Down
-- reverse: create index "idx_sessions_hashed_session_token" to table: "sessions"
DROP INDEX "public"."idx_sessions_hashed_session_token";
-- reverse: create index "idx_sessions_hashed_refresh_token" to table: "sessions"
DROP INDEX "public"."idx_sessions_hashed_refresh_token";
-- reverse: create index "idx_sessions_deleted_at" to table: "sessions"
DROP INDEX "public"."idx_sessions_deleted_at";
-- reverse: create "sessions" table
DROP TABLE "public"."sessions";
-- reverse: create index "idx_users_username" to table: "users"
DROP INDEX "public"."idx_users_username";
-- reverse: create index "idx_users_email" to table: "users"
DROP INDEX "public"."idx_users_email";
-- reverse: create index "idx_users_deleted_at" to table: "users"
DROP INDEX "public"."idx_users_deleted_at";
-- reverse: create "users" table
DROP TABLE "public"."users";
-- reverse: create index "idx_folder_manifests_game_id" to table: "folder_manifests"
DROP INDEX "public"."idx_folder_manifests_game_id";
-- reverse: create index "idx_folder_manifests_deleted_at" to table: "folder_manifests"
DROP INDEX "public"."idx_folder_manifests_deleted_at";
-- reverse: create "folder_manifests" table
DROP TABLE "public"."folder_manifests";
-- reverse: create index "idx_service_configs_deleted_at" to table: "service_configs"
DROP INDEX "public"."idx_service_configs_deleted_at";
-- reverse: create index "idx_service_config" to table: "service_configs"
DROP INDEX "public"."idx_service_config";
-- reverse: create "service_configs" table
DROP TABLE "public"."service_configs";
-- reverse: create index "idx_provider_game" to table: "games"
DROP INDEX "public"."idx_provider_game";
-- reverse: create index "idx_games_exe" to table: "games"
DROP INDEX "public"."idx_games_exe";
-- reverse: create index "idx_games_deleted_at" to table: "games"
DROP INDEX "public"."idx_games_deleted_at";
-- reverse: create "games" table
DROP TABLE "public"."games";
-- reverse: create index "idx_assets_type" to table: "assets"
DROP INDEX "public"."idx_assets_type";
-- reverse: create index "idx_assets_game_id" to table: "assets"
DROP INDEX "public"."idx_assets_game_id";
-- reverse: create index "idx_assets_deleted_at" to table: "assets"
DROP INDEX "public"."idx_assets_deleted_at";
-- reverse: create "assets" table
DROP TABLE "public"."assets";
