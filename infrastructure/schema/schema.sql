-- Copyright 2023 Google LLC
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     https://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.
--

CREATE TABLE players (
  player_google_id STRING(MAX) NOT NULL,
  player_name STRING(64) NOT NULL,
  profile_image STRING(MAX) NOT NULL,
  region STRING(10) NOT NULL,
  skill_level int64 NOT NULL,
  tier STRING(1) NOT NULL,
  stats JSON,
) PRIMARY KEY(player_google_id);

CREATE UNIQUE INDEX PlayerName ON players(player_name);

CREATE TABLE game_assets
(
  asset_uuid STRING(36) NOT NULL,
  asset_name STRING(MAX) NOT NULL,
  available_time TIMESTAMP NOT NULL
)PRIMARY KEY (asset_uuid);

CREATE TABLE player_assets
(
  player_asset_uuid STRING(36) NOT NULL,
  player_google_id STRING(MAX) NOT NULL,
  asset_uuid STRING(36) NOT NULL,
  acquire_time TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP()),
  FOREIGN KEY (asset_uuid) REFERENCES game_assets (asset_uuid)
) PRIMARY KEY (player_google_id, player_asset_uuid),
    INTERLEAVE IN PARENT players ON DELETE CASCADE;
