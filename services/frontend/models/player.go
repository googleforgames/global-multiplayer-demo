// Copyright 2023 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

// SingleGameStats provides a structure for updating a player's stats based on a single game outcome
type SingleGameStats struct {
	Player_google_id string `json:"player_google_id" uri:"id"`
	Won              bool   `json:"won"`
	Score            int64  `json:"score"`
	Kills            int64  `json:"kills"`
	Deaths           int64  `json:"deaths"`
}

// PlayerStats provides various statistics for a player
type PlayerStats struct {
	Games_played int64 `json:"games_played"`
	Games_won    int64 `json:"games_won"`
	Total_score  int64 `json:"total_score"`
	Total_kills  int64 `json:"total_kills"`
	Total_deaths int64 `json:"total_deaths"`
}

// Player maps to the fields stored for the backend database
type Player struct {
	Player_google_id string      `json:"player_google_id"`
	Player_name      string      `json:"player_name"`
	Profile_image    string      `json:"profile_image"`
	Region           string      `json:"region"`
	Stats            PlayerStats `json:"stats"`
	Skill_level      int64       `json:"skill_level"`
	Tier             string      `json:"tier"`
}
