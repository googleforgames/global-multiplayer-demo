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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/googleforgames/global-multiplayer-demo/services/frontend-api/match"
	"github.com/googleforgames/global-multiplayer-demo/services/frontend-api/models"
	"github.com/googleforgames/global-multiplayer-demo/services/frontend-api/shared"
	"github.com/googleforgames/global-multiplayer-demo/services/frontend-api/shared/auth"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	oauthStateString = "r4nd0ms7r1ng"
)

func main() {
	// Load local .env
	godotenv.Load()
	shared.ValidateEnvVars()

	// Oauth config from env vars
	googleOauthConfig.ClientID = os.Getenv("CLIENT_ID")
	googleOauthConfig.ClientSecret = os.Getenv("CLIENT_SECRET")
	googleOauthConfig.RedirectURL = os.Getenv("CALLBACK_HOSTNAME")

	r := gin.Default()

	// TODO: Better configuration of trusted proxy
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("could not set trusted proxies: %s", err)
	}

	m, err := match.NewMatcher()
	if err != nil {
		log.Fatalf("could not initialize matcher: %v", err)
	}
	defer m.Close()

	r.GET("/login", handleGoogleLogin)
	r.GET("/callback", handleGoogleCallback)

	// JWT protected endpoint handlers
	r.POST("/play", auth.VerifyJWT(func(id string, c *gin.Context) { handlePlay(id, c, m) }))
	r.GET("/profile", auth.VerifyJWT(handleProfile))
	r.GET("/stats", auth.VerifyJWT(handleGetStats))
	r.PUT("/stats", auth.VerifyJWT(handleUpdateStats))
	r.GET("/ping", auth.VerifyJWT(handlePingServers))

	log.Printf("Google for Games Frontend API is listening on :%s\n", os.Getenv("LISTEN_PORT"))

	if err := r.Run(":" + os.Getenv("LISTEN_PORT")); err != nil {
		log.Fatalf("could not run gin router: %s", err)
		return
	}
}

// Generates a redirect to google's login
func handleGoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

// Callback handler that gets the access token for further profile querying of Google APIs
func handleGoogleCallback(c *gin.Context) {

	state := c.Request.FormValue("state")
	if state != oauthStateString {
		err := fmt.Errorf("invalid oauth state, expected '%s', got '%s'", oauthStateString, state)
		if shared.HandleError(c, http.StatusBadRequest, "auth callback", err) {
			return
		}
	}

	code := c.Request.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if shared.HandleError(c, http.StatusBadRequest, "auth exchange", err) {
		return
	}

	// Call our profile service to create a new entry in Spanner DB
	id, err := createProfileIfNotExists(token.AccessToken)
	if shared.HandleError(c, http.StatusInternalServerError, "creating profile if not exists", err) {
		return
	}

	// Generate our own jwt token (1 month validity) that will be used in the game launcher / game client
	jwtToken, err := auth.GenerateJWT(id, 31)
	if shared.HandleError(c, http.StatusInternalServerError, "token generation", err) {
		return
	}

	// Redirect to the launcher callback port
	http.Redirect(c.Writer, c.Request, fmt.Sprintf("http://localhost:%s/callback?token=%s", os.Getenv("CLIENT_LAUNCHER_PORT"), jwtToken), http.StatusTemporaryRedirect)
}

// Profile handling endpoint
func handleProfile(id string, c *gin.Context) {
	// Query our own profile service to get player data
	response, err := http.Get(fmt.Sprintf("%s/players/%s", os.Getenv("PROFILE_SERVICE"), id))
	if shared.HandleError(c, http.StatusInternalServerError, "fetching profile", err) {
		return
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		// Profile found, decode and show
		var p models.Player
		err := json.NewDecoder(response.Body).Decode(&p)
		if shared.HandleError(c, http.StatusInternalServerError, "decoding profile", err) {
			return
		}

		c.JSON(http.StatusOK, p)
	} else if response.StatusCode == 404 { // If not found, return an error
		err := fmt.Errorf("profile not found: %s", id)
		if shared.HandleError(c, http.StatusBadRequest, "profile lookup", err) {
			return
		}
	} else {
		err := fmt.Errorf("unable to fetch profile, error code: %d", response.StatusCode)
		if shared.HandleError(c, http.StatusBadRequest, "profile lookup", err) {
			return
		}
	}
}

// Getting the stats from profile api
func handleGetStats(id string, c *gin.Context) {
	response, err := http.Get(fmt.Sprintf("%s/players/%s/stats", os.Getenv("PROFILE_SERVICE"), id))
	if shared.HandleError(c, http.StatusInternalServerError, "fetching profile", err) {
		return
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		// Profile found, decode and show
		var p models.Player
		err := json.NewDecoder(response.Body).Decode(&p)
		if shared.HandleError(c, http.StatusInternalServerError, "decoding profile", err) {
			return
		}

		c.JSON(http.StatusOK, p.Stats)
	} else if response.StatusCode == 404 { // If not found, return an error
		err := fmt.Errorf("profile not found: %s", id)
		if shared.HandleError(c, http.StatusBadRequest, "profile lookup", err) {
			return
		}
	} else {
		err := fmt.Errorf("unable to fetch profile, error code: %d", response.StatusCode)
		if shared.HandleError(c, http.StatusBadRequest, "profile lookup", err) {
			return
		}
	}
}

// Updating the stats in the profile api
func handleUpdateStats(id string, c *gin.Context) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/players/%s/stats", os.Getenv("PROFILE_SERVICE"), id), c.Request.Body)
	if shared.HandleError(c, http.StatusInternalServerError, "stats update", err) {
		return
	}

	response, err := client.Do(req)
	if shared.HandleError(c, http.StatusInternalServerError, "stats update", err) {
		return
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		c.JSON(http.StatusOK, "OK")
		return
	} else {
		err := fmt.Errorf("unable to update profile stats, error code: %d", response.StatusCode)
		if shared.HandleError(c, http.StatusBadRequest, "stats update", err) {
			return
		}

	}
}

// WIP: Needs an endpoint to fetch the ping servers
func handlePingServers(id string, c *gin.Context) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/list", os.Getenv("PING_SERVICE")), c.Request.Body)
	if shared.HandleError(c, http.StatusInternalServerError, "fetch ping servers", err) {
		return
	}

	response, err := client.Do(req)
	if shared.HandleError(c, http.StatusInternalServerError, "fetch ping servers", err) {
		return
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		var pingServers map[string]models.PingServer
		err := json.NewDecoder(response.Body).Decode(&pingServers)
		if shared.HandleError(c, http.StatusInternalServerError, "decoding ping servers", err) {
			return
		}

		c.JSON(http.StatusOK, pingServers)
		return
	} else {
		err := fmt.Errorf("unable to update profile stats, error code: %d", response.StatusCode)
		if shared.HandleError(c, http.StatusBadRequest, "stats update", err) {
			return
		}

	}
}

// WIP: Handles the play request from the game client
func handlePlay(id string, c *gin.Context, m *match.Matcher) {

	host, hok := os.LookupEnv("LOCAL_OPENMATCH_SERVER_OVERRIDE_HOST")
	port, pok := os.LookupEnv("LOCAL_OPENMATCH_SERVER_OVERRIDE_PORT")
	if hok && pok {
		port, _ := strconv.Atoi(port)
		fmt.Printf("Overriding openmatch response with %s:%d\n", host, port)
		c.JSON(http.StatusOK, models.OMServerResponse{IP: host, Port: port})
		return
	}

	pr := &models.PlayRequest{}
	if err := c.Bind(pr); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// TODO #1: Get profile here (from Cloud Spanner via token/id??)
	// TODO #2: Add profile parameter for finding the server (besides the preferred region)

	conn, err := m.FindMatchingServer(c.Request.Context(), pr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, conn)
}

// Function responsible for checking if profile is not yet created in our own profile service
func createProfileIfNotExists(token string) (string, error) {

	// Fetch the data from Google's API by access token
	response, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("unable to fetch google's user profile, error code: %d", response.StatusCode)
	}

	var udata bytes.Buffer
	_, err = io.Copy(&udata, response.Body)

	if err != nil {
		return "", err
	}

	var userInfo models.UserInfo
	err = json.Unmarshal(udata.Bytes(), &userInfo)
	if err != nil {
		return "", err
	}

	// Fetch profile data from our profile service
	response, err = http.Get(fmt.Sprintf("%s/players/%s", os.Getenv("PROFILE_SERVICE"), userInfo.Sub))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// If we get a 404, means that the player doesn't exist and we need to create it via POST call
	if response.StatusCode == 404 {
		// Init the model with data fro mgoogle's api
		p := models.Player{
			Player_google_id: userInfo.Sub,
			Player_name:      userInfo.Name,
			Profile_image:    userInfo.Picture,
			Region:           userInfo.Locale,
			Stats:            models.PlayerStats{},
			Skill_level:      0,
			Tier:             "U",
		}

		profileData, _ := json.Marshal(p)
		request, err := http.NewRequest("POST", os.Getenv("PROFILE_SERVICE")+"/players", bytes.NewBuffer(profileData))
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			return "", err
		}

		client := &http.Client{}
		response, err := client.Do(request)

		if err != nil {
			return "", err
		}

		defer response.Body.Close()

	}

	// Return google's user ID that we use as a primary in our profile service
	return userInfo.Sub, nil
}
