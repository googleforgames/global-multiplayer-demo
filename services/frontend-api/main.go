package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const htmlIndex = `<html><body>
<a href="/login">Log in with Google</a>
</body></html>
`

type GoogleOauthToken struct {
	AccessToken  string
	RefreshToken string
	Expiry       string
	TokenType    string
	IdToken      string
}

type UserInfo struct {
	Id            string `json:"id"`
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

type OMServerResponse struct {
	IP     string
	Port   int
	Region string
}

type RegionalLatency struct {
	Region  string
	Latency int
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL: "http://localhost:8080/callback",
		Scopes:      []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:    google.Endpoint,
	}
	// Some random string, random for each request
	oauthStateString = "random"
)

func main() {
	godotenv.Load()

	googleOauthConfig.ClientID = os.Getenv("APP_CLIENT_ID")
	googleOauthConfig.ClientSecret = os.Getenv("APP_CLIENT_SECRET")

	log.Printf("%+v\n", googleOauthConfig)

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	http.HandleFunc("/profile", handleProfileInfo)
	http.HandleFunc("/play", handlePlay)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil || response.StatusCode != 200 {
		fmt.Printf("Failed getting user info: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer response.Body.Close()
	// Use response.Body to get user information.

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var result UserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Token: %s\n", token.AccessToken)
	fmt.Fprintf(w, "Login Successful!")
}

func handleProfileInfo(rw http.ResponseWriter, req *http.Request) {
	token := req.FormValue("token")
	token = refreshToken(token)

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil || response.StatusCode != 200 {
		fmt.Printf("Failed getting user info: %s\n", err)
		http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
		return
	}

	defer response.Body.Close()
	// Use response.Body to get user information.

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var result UserInfo
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}

	jsonResponse, _ := json.Marshal(result)

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")

	fmt.Fprintln(rw, string(jsonResponse))
}

func handlePlay(rw http.ResponseWriter, req *http.Request) {
	/*defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "{\"error\": \"%s\"}", err)
			log.Println("panic occurred:", err)
		}
	}()*/

	/*token := req.FormValue("token")
	token = refreshToken(token)

	//Get profile here (from Cloud Spanner via token/id??)
	*/

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")

	LatencyList := parseRegionalLatencies(rw, req)
	for _, LatencyItem := range LatencyList {
		log.Printf("Latency for %s is %d\n", LatencyItem.Region, LatencyItem.Latency)
	}

	// Sort by latency
	sort.SliceStable(LatencyList, func(i, j int) bool {
		return LatencyList[i].Latency < LatencyList[j].Latency
	})

	mOMResponse, _ := json.Marshal(findMatchingServer(LatencyList)) // add profile param for finding the server

	fmt.Fprint(rw, string(mOMResponse))
}

func refreshToken(t string) string {
	token := oauth2.Token{
		AccessToken: t,
		TokenType:   "bearer",
	}

	tokenSource := googleOauthConfig.TokenSource(context.TODO(), &token)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalln(err)
	}

	if newToken.AccessToken != token.AccessToken {
		log.Println("Saved new token: ", newToken.AccessToken)
	} else {
		log.Println("Old token still good: ", token.AccessToken)
	}

	return newToken.AccessToken
}

func findMatchingServer(LatencyList []RegionalLatency) OMServerResponse {
	if len(LatencyList) > 0 {
		log.Printf("Looking for a server in the %s region.\n", LatencyList[0].Region)
	} else {
		LatencyList = append(LatencyList, RegionalLatency{Region: "n/a", Latency: -1})
	}

	// TODO: Query OpenMatch on `OMFrontendEndpoint` in a specified region for a server.

	IP := "127.0.0.1"
	Port := 7777

	return OMServerResponse{
		IP:     IP,
		Port:   Port,
		Region: LatencyList[0].Region}

}

func parseRegionalLatencies(rw http.ResponseWriter, req *http.Request) []RegionalLatency {
	decoder := json.NewDecoder(req.Body)
	var LatencyList []RegionalLatency
	err := decoder.Decode(&LatencyList)

	if err != nil {
		return LatencyList
	}

	return LatencyList
}
