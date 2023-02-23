// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package models interacts with the backend database to handle the stateful
// data for the profile service.
package models

import (
	"context"
	"encoding/json"
	"fmt"

	spanner "cloud.google.com/go/spanner"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

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
	Player_google_id string           `json:"player_google_id" validate:"required" uri:"id"`
	Player_name      string           `json:"player_name"`
	Profile_image    string           `json:"profile_image"`
	Region           string           `json:"region"`
	Stats            spanner.NullJSON `json:"stats"`
	Skill_level      int64            `json:"skill_level"`
	Tier             string           `json:"tier"`
}

// Validate that the player has the required information based on the type's validation rules.
func (p *Player) Validate() error {
	validate = validator.New()
	err := validate.Struct(p)
	if err != nil {
		return err
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	return nil
}

// GetPlayerByGoogleId returns a Player based on a provided google_id. In the event of an error
// retrieving the player, an empty Player is returned with the error.
func GetPlayerByGoogleId(ctx context.Context, client spanner.Client, google_id string) (Player, error) {
	// Retrieve most columns of the player. Does not retrieve stats, as this is not necessary for most player requests.
	row, err := client.Single().ReadRow(ctx, "players",
		spanner.Key{google_id}, []string{"player_google_id", "player_name", "profile_image", "region", "skill_level", "tier"})
	if err != nil {
		return Player{}, err
	}

	player := Player{}
	err = row.ToStruct(&player)

	if err != nil {
		return Player{}, err
	}
	return player, nil
}

// AddPlayer provides functionality to insert a player into the backend.
// Provide with the required fields from the API call. This is then inserted, along with empty stats, into
// the Spanner database.
func (p *Player) AddPlayer(ctx context.Context, client spanner.Client) error {
	// Validate based on struct validation rules
	err := p.Validate()
	if err != nil {
		return err
	}

	// Initialize player stats
	emptyStats := spanner.NullJSON{Value: PlayerStats{
		Games_played: 0,
		Games_won:    0,
		Total_score:  0,
		Total_kills:  0,
		Total_deaths: 0,
	}, Valid: true}

	// insert into spanner.
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT players (player_google_id, player_name, profile_image, region, skill_level, tier, stats) VALUES
					(@playerGoogleId, @playerName, @profileImage, @region, @skill, @tier, @pStats)
			`,
			Params: map[string]interface{}{
				"playerGoogleId": p.Player_google_id,
				"playerName":     p.Player_name,
				"profileImage":   p.Profile_image,
				"region":         p.Region,
				"skill":          0,   // Initial skill rating
				"tier":           "U", // Initial tier is U=Unknown
				"pStats":         emptyStats,
			},
		}

		_, err := txn.Update(ctx, stmt)
		return err
	})

	// TODO: Handle 'AlreadyExists' errors
	if err != nil {
		return err
	}

	// return empty error on success
	return nil
}

// UpdatePlayer updates a game's player profile with provided information
func (p *Player) UpdatePlayer(ctx context.Context, client spanner.Client) error {
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// Update player
		cols := []string{"player_google_id", "player_name", "profile_image", "region"}

		err := txn.BufferWrite([]*spanner.Mutation{
			spanner.Update("players", cols, []interface{}{p.Player_google_id, p.Player_name, p.Profile_image, p.Region}),
		})

		if err != nil {
			return fmt.Errorf("could not buffer write: %s", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// GetPlayerStats returns a Player's stats based on a provided google_id. In the event of an error
// retrieving the player, an empty Player is returned with the error.
func GetPlayerStats(ctx context.Context, client spanner.Client, google_id string) (Player, error) {
	// Retrieve columns related to player stats.
	row, err := client.Single().ReadRow(ctx, "players",
		spanner.Key{google_id}, []string{"player_google_id", "stats", "skill_level", "tier"})
	if err != nil {
		return Player{}, err
	}

	player := Player{}
	err = row.ToStruct(&player)

	if err != nil {
		return Player{}, err
	}
	return player, nil
}

// UpdateStats updates a player's stats with statistics of a game's outcome
func UpdateStats(ctx context.Context, client spanner.Client, gStats SingleGameStats) (Player, error) {
	player := Player{}

	// Transaction to request and update player's stats with the new stats
	// We must read the stats first before updating, because there is no function
	// to update values in place. This is done in a transaction to avoide concurrent updates
	// on a single player losing stats.
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// Retrieve columns related to player stats.
		row, err := txn.ReadRow(ctx, "players",
			spanner.Key{gStats.Player_google_id}, []string{"player_google_id", "stats", "skill_level", "tier"})

		if err != nil {
			return err
		}

		err = row.ToStruct(&player)

		if err != nil {
			return err
		}

		// Modify stat totals
		var pStats PlayerStats
		if err := json.Unmarshal([]byte(player.Stats.String()), &pStats); err != nil {
			return fmt.Errorf("could not unmarshal json: %s", err)
		}

		pStats.Games_played = pStats.Games_played + 1
		pStats.Total_score = pStats.Total_score + gStats.Score
		pStats.Total_kills = pStats.Total_kills + gStats.Kills
		pStats.Total_deaths = pStats.Total_deaths + gStats.Deaths

		if gStats.Won {
			pStats.Games_won = pStats.Games_won + 1
		}

		if pStats.Total_deaths != 0 {
			player.Skill_level = pStats.Total_kills / pStats.Total_deaths
		} else {
			player.Skill_level = pStats.Total_kills
		}

		updatedStats, _ := json.Marshal(pStats)
		if err := player.Stats.UnmarshalJSON(updatedStats); err != nil {
			return fmt.Errorf("could not unmarshal json: %s", err)
		}

		// TODO: Modify tier

		// Update player
		cols := []string{"player_google_id", "stats", "skill_level"}

		err = txn.BufferWrite([]*spanner.Mutation{
			spanner.Update("players", cols, []interface{}{player.Player_google_id, player.Stats, player.Skill_level}),
		})

		if err != nil {
			return fmt.Errorf("could not buffer write: %s", err)
		}

		return nil
	})

	if err != nil {
		return Player{}, err
	}

	return player, nil
}
