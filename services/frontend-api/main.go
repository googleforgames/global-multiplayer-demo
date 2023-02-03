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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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

	// Oauth config from env vars
	googleOauthConfig.ClientID = os.Getenv("CLIENT_ID")
	googleOauthConfig.ClientSecret = os.Getenv("CLIENT_SECRET")
	googleOauthConfig.RedirectURL = "http://localhost:" + os.Getenv("LISTEN_PORT") + "/callback"

	// Endpoint handlers
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	// JWT protected endpoint handlers
	http.HandleFunc("/play", auth.VerifyJWT(handlePlay))
	http.HandleFunc("/profile", auth.VerifyJWT(handleProfile))
	http.HandleFunc("/stats", auth.VerifyJWT(handleStats))
	http.HandleFunc("/ping", auth.VerifyJWT(handlePingServers))

	fmt.Println("Google for Games Frontend API is listening on :" + os.Getenv("LISTEN_PORT"))
	fmt.Println(http.ListenAndServe(":"+os.Getenv("LISTEN_PORT"), nil))
}

// Generates a redirect to google's login
func handleGoogleLogin(rw http.ResponseWriter, req *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	http.Redirect(rw, req, url, http.StatusTemporaryRedirect)
}

// Callback handler that gets the access token for further profile querying of Google APIs
func handleGoogleCallback(rw http.ResponseWriter, req *http.Request) {
	// Generic panic recovery that spits out 500 response code and a json formatted error
	defer shared.RecoverFromPanic(rw, req)

	state := req.FormValue("state")
	if state != oauthStateString {
		panic(fmt.Sprintf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state))
	}

	code := req.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		panic(fmt.Sprintf("Code exchange failed with '%s'\n", err))
	}

	// Call our profile service to create a new entry in Spanner DB
	id := createProfileIfNotExists(token.AccessToken)

	// Generate our own jwt token (1 month validity) that will be used in the game launcher / game client
	jwtToken, err := auth.GenerateJWT(id, 31)
	if err != nil {
		panic(fmt.Sprintf("Unable to generate JWT token '%s'\n", err))
	}

	// Redirect to the launcher callback port
	http.Redirect(rw, req, fmt.Sprintf("http://localhost:%s/callback?token=%s", os.Getenv("CLIENT_LAUNCHER_PORT"), jwtToken), http.StatusTemporaryRedirect)
}

// Profile handling endpoint
func handleProfile(id string, rw http.ResponseWriter, req *http.Request) {
	defer shared.RecoverFromPanic(rw, req)

	// Query our own profile service to get player data
	response, err := http.Get(fmt.Sprintf("%s/players/%s", os.Getenv("PROFILE_SERVICE"), id))
	if err != nil {
		panic(err.Error())
	}

	defer response.Body.Close()

	// If not found, return an error
	if response.StatusCode == 404 {
		panic("Profile not found: " + id)
	} else if response.StatusCode == 200 {
		// Profile found, decode and show
		var p models.Player
		err := json.NewDecoder(response.Body).Decode(&p)
		if err != nil {
			fmt.Printf("Failed getting user info: %s\n", err)
			panic(err)
		}

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(p)
		return
	}

	panic("Unable to process player request")
}

// Stats handler for GET and PUT methods
func handleStats(id string, rw http.ResponseWriter, req *http.Request) {
	defer shared.RecoverFromPanic(rw, req)

	if req.Method == "PUT" {
		handleUpdateStats(id, rw, req)
	} else if req.Method == "GET" {
		handleGetStats(id, rw, req)
	}
}

func handleGetStats(id string, rw http.ResponseWriter, req *http.Request) {
	response, err := http.Get(fmt.Sprintf("%s/players/%s/stats", os.Getenv("PROFILE_SERVICE"), id))
	if err != nil {
		panic(err.Error())
	}

	defer response.Body.Close()

	if response.StatusCode == 404 {
		panic("Profile stats not found")
	} else if response.StatusCode == 200 {

		var p models.Player
		err := json.NewDecoder(response.Body).Decode(&p)
		if err != nil {
			fmt.Printf("Failed getting user info: %s\n", err)
			panic(err)
		}

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(p.Stats)
		return
	}

	panic("Unable to process player request")
}

func handleUpdateStats(id string, rw http.ResponseWriter, req *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/players/%s/stats", os.Getenv("PROFILE_SERVICE"), id), req.Body)
	if err != nil {
		panic(err.Error())
	}

	response, err := client.Do(req)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Content-Type", "application/json")

		json.NewEncoder(rw).Encode("ok")
		return
	}

	panic("Unable to process player request")
}

// WIP: Needs an endpoint to fetch the ping servers
func handlePingServers(id string, rw http.ResponseWriter, req *http.Request) {
	defer shared.RecoverFromPanic(rw, req)

	// TODO: fetch servers from some tbd endpoint

	var pingServers []models.PingServer = []models.PingServer{
		{IP: "216.239.32.53", Region: "Tokyo"},
		{IP: "216.239.38.53", Region: "London"},
		{IP: "216.239.34.53", Region: "North Virginia"},
	}

	PingResponse, _ := json.Marshal(pingServers)

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(PingResponse))

}

// WIP: Handles the play request from the game client
func handlePlay(id string, rw http.ResponseWriter, req *http.Request) {
	defer shared.RecoverFromPanic(rw, req)

	// Get regions by preferred order
	preferredRegions := strings.Split(req.FormValue("preferred_regions"), ",")
	for _, region := range preferredRegions {
		log.Println(region)
	}

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")

	// TODO #1: Get profile here (from Cloud Spanner via token/id??)
	// TODO #2: Add profile parameter for finding the server (besides the preferred region)
	mOMResponse, _ := json.Marshal(models.FindMatchingServer(preferredRegions))

	fmt.Fprint(rw, string(mOMResponse))
}

// Function responsible for checking if profile is not yet created in our own profile service
func createProfileIfNotExists(token string) string {

	// Fetch the data from Google's API by access token
	response, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + token)
	if err != nil || response.StatusCode != 200 {
		panic("Failed getting user info: %s\n")
	}

	defer response.Body.Close()

	udata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var userInfo models.UserInfo
	if err := json.Unmarshal(udata, &userInfo); err != nil {
		fmt.Printf("Failed getting user info: %s\n", err)
		panic(err)
	}

	// Fetch profile data from our profile service
	response, err = http.Get(fmt.Sprintf("%s/players/%s", os.Getenv("PROFILE_SERVICE"), userInfo.Sub))

	if err != nil {
		panic(err.Error())
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

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

	}

	// Return google's user ID that we use as a primary in our profile service
	return userInfo.Sub
}
