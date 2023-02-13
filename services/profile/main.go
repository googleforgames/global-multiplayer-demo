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

// Package main exposes the REST endpoints for the profile-service.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	spanner "cloud.google.com/go/spanner"
	"github.com/googleforgames/global-multiplayer-demo/profile-service/config"
	"github.com/googleforgames/global-multiplayer-demo/profile-service/models"

	"github.com/gin-gonic/gin"
)

// setSpannerConnection is a mutator to create spanner context and client, and set them in gin
func setSpannerConnection(c config.Config) gin.HandlerFunc {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, c.Spanner.DB())

	if err != nil {
		log.Fatal(err)
	}

	return func(c *gin.Context) {
		c.Set("spanner_client", *client)
		c.Set("spanner_context", ctx)
		c.Next()
	}
}

// getSpannerConnection is a helper function to retrieve spanner client and context
func getSpannerConnection(c *gin.Context) (context.Context, spanner.Client) {
	return c.MustGet("spanner_context").(context.Context),
		c.MustGet("spanner_client").(spanner.Client)

}

// getPlayerID responds to the GET /players/:id endpoint
// Returns a player's information when provided a valid player_google_id
func getPlayerByID(c *gin.Context) {
	var playerGoogleId = c.Param("id")

	ctx, client := getSpannerConnection(c)

	player, err := models.GetPlayerByGoogleId(ctx, client, playerGoogleId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "player not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, player)
}

// createPlayer responds to the POST /players endpoint
// When provided the required fields of player_name, profile_image and region, creates a player.
func createPlayer(c *gin.Context) {
	var player models.Player

	if err := c.BindJSON(&player); err != nil {
		if err := c.AbortWithError(http.StatusBadRequest, err); err != nil {
			fmt.Printf("could not abort: %s", err)
		}
		return
	}

	ctx, client := getSpannerConnection(c)
	err := player.AddPlayer(ctx, client)
	if err != nil {
		if err := c.AbortWithError(http.StatusBadRequest, err); err != nil {
			fmt.Printf("could not abort: %s", err)
		}
		return
	}

	c.IndentedJSON(http.StatusCreated, player.Player_google_id)
}

// updatePlayer responds to the PUT /players endpoint
// Updates the player's profile information. Does not include updating stats, which is handled
// by a different endpoint.
func updatePlayer(c *gin.Context) {
	var player models.Player

	if err := c.BindJSON(&player); err != nil {
		if err := c.AbortWithError(http.StatusBadRequest, err); err != nil {
			fmt.Printf("could not abort: %s", err)
		}
		return
	}

	ctx, client := getSpannerConnection(c)
	err := player.UpdatePlayer(ctx, client)
	if err != nil {
		if err := c.AbortWithError(http.StatusBadRequest, err); err != nil {
			fmt.Printf("could not abort: %s", err)
		}
		return
	}

	c.IndentedJSON(http.StatusOK, player)
}

// getPlayerStats responds to the GET /players/:id/stats endpoint
// Returns a player's stats when provided a valid player_google_id
func getPlayerStats(c *gin.Context) {
	var playerGoogleId = c.Param("id")

	ctx, client := getSpannerConnection(c)

	player, err := models.GetPlayerStats(ctx, client, playerGoogleId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "player not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, player)
}

// updatePlayerStats responds to the PUT /player/stats endpoint
// Updates the player's stats based on provided payload
func updatePlayerStats(c *gin.Context) {
	game_stats := models.SingleGameStats{}

	err := c.BindUri(&game_stats)

	// Error binding the google_id from the URI
	if err != nil {
		if err := c.AbortWithError(http.StatusBadRequest, err); err != nil {
			fmt.Printf("could not abort: %s", err)
		}
		return
	}

	if err := c.BindJSON(&game_stats); err != nil {
		if err := c.AbortWithError(http.StatusBadRequest, err); err != nil {
			fmt.Printf("could not abort: %s", err)
		}
		return
	}

	ctx, client := getSpannerConnection(c)

	player, err := models.UpdateStats(ctx, client, game_stats)
	if err != nil {
		if err := c.AbortWithError(http.StatusBadRequest, err); err != nil {
			fmt.Printf("could not abort: %s", err)
		}
		return
	}

	c.IndentedJSON(http.StatusOK, player)
}

// main initializes the gin router and configures the endpoints
func main() {
	configuration, _ := config.NewConfig()

	router := gin.Default()
	// TODO: Better configuration of trusted proxy
	if err := router.SetTrustedProxies(nil); err != nil {
		fmt.Printf("could not set trusted proxies: %s", err)
		return
	}

	router.Use(setSpannerConnection(configuration))

	router.POST("/players", createPlayer)
	router.GET("/players/:id", getPlayerByID)
	router.PUT("/players", updatePlayer)
	router.GET("/players/:id/stats", getPlayerStats)
	router.PUT("/players/:id/stats", updatePlayerStats)

	if err := router.Run(configuration.Server.URL()); err != nil {
		fmt.Printf("could not run gin router: %s", err)
		return
	}
}
