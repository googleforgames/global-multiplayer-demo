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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/ini.v1"
)

var (
	myToken  string
	myApp    fyne.App
	myWindow fyne.Window
	iniCfg   *ini.File
)

func main() {
	// Load configs
	var err error
	iniCfg, err = ini.Load("app.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
		os.Exit(1)
	}

	myToken = ""

	// Callback handling from the frontend api
	http.HandleFunc("/callback", handleGoogleCallback)
	go func() {
		log.Printf("Google for Games Launcher is listening for callbacks on :%s", iniCfg.Section("").Key("callback_listen_port").String())
		log.Println(http.ListenAndServe(":"+iniCfg.Section("").Key("callback_listen_port").String(), nil))
	}()

	// UI
	myApp = app.New()
	myWindow = myApp.NewWindow("Google for Games Launcher")
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.NewSize(320, 260))
	myWindow.CenterOnScreen()

	image := canvas.NewImageFromFile("assets/header.png")
	image.FillMode = canvas.ImageFillContain

	buttonSignIn := widget.NewButtonWithIcon("Sign-in with Google", theme.HomeIcon(), func() {
		openBrowser(iniCfg.Section("").Key("frontend_api").String() + "/login")
	})

	buttonExit := widget.NewButtonWithIcon("Exit", theme.CancelIcon(), func() {
		log.Println("Tapped exit")
		myApp.Quit()
	})

	subGrid := container.New(layout.NewGridLayout(1), layout.NewSpacer(), buttonSignIn, buttonExit)
	grid := container.New(layout.NewGridLayout(1), image, subGrid)

	myWindow.SetContent(grid)

	// If we have a valid token, let's use it and update the UI right away
	if loadToken() {
		playerName := getPlayerName()
		updateUI(playerName)
	}

	myWindow.ShowAndRun()
}

func handleGoogleCallback(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)

			fmt.Fprintf(rw, "{\"error\": \"%s\"}", err)
			log.Printf("panic occurred: %s", err)
		}
	}()

	// Save my token
	myToken = req.FormValue("token")
	if len(myToken) == 0 {
		log.Fatal("No token received!")
	}

	saveToken(myToken)

	// Update UI with profile info and launch game button
	playerName := getPlayerName()
	updateUI(playerName)

	// Close the browser window
	closeScript := `<script> 
		setTimeout("window.close()",3000) 
	</script>
	<p>
		<h2>Authenticated successfully. Please return to your application. This tab will close in 3 seconds.</h2>
	</p>`
	fmt.Fprintf(rw, closeScript)
}

func updateUI(playerName string) {
	// Update UI with profile info and launch game button
	image := canvas.NewImageFromFile("assets/header.png")
	image.FillMode = canvas.ImageFillContain

	label1 := widget.NewLabel(fmt.Sprintf("Welcome %s!", playerName))
	label1.Alignment = fyne.TextAlignCenter
	label2 := widget.NewLabel("Are you ready to play again?!")
	label2.Alignment = fyne.TextAlignCenter

	buttonPlay := widget.NewButtonWithIcon("Open Droidshooter", theme.MediaPlayIcon(), func() {
		log.Println("Tapped Play!")
		handlePlay()
	})

	buttonExit := widget.NewButtonWithIcon("Exit", theme.CancelIcon(), func() {
		log.Println("Tapped exit")
		myApp.Quit()
	})

	infoGrid := container.New(layout.NewGridLayout(1), label1, label2)
	subGrid := container.New(layout.NewGridLayout(1), infoGrid, buttonPlay, buttonExit)
	grid := container.New(layout.NewGridLayout(1), image, subGrid)
	myWindow.SetContent(grid)
}

func handlePlay() {
	params := fmt.Sprintf("-token=%s", myToken)

	// Get the binary file from the ini
	cmd := exec.Command(iniCfg.Section(runtime.GOOS).Key("binary").String(), params)
	log.Printf("Launching: %s %s", iniCfg.Section(runtime.GOOS).Key("binary").String(), params)

	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error: %s", err)
	}
}

func getPlayerName() string {
	log.Printf("Getting player info")

	req, err := http.NewRequest("GET", iniCfg.Section("").Key("frontend_api").String()+"/profile", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", myToken))

	if err != nil {
		log.Fatal("Unable to initiate request to game api. Connection issues?")
	}

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != 200 {
		log.Fatalf("Unable to fetch user information. Expired token?: %s", data)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatalf("Unable to decode json: %s", err)
	}

	return result["player_name"].(string)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func saveToken(token string) {
	dirname, err := os.UserHomeDir()
	filename := dirname + "/droidshooter.jwt"

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filename, []byte(token), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func loadToken() bool {
	dirname, err := os.UserHomeDir()
	filename := dirname + "/droidshooter.jwt"

	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return false
	}

	file, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}

	modifiedtime := file.ModTime()
	in30Days := time.Now().Add(24 * time.Hour * 30)

	if modifiedtime.After(in30Days) {
		log.Printf("Token is old. Deleting.")
		os.Remove(filename)
		return false
	}

	myToken = string(data)
	log.Printf("Token loaded from file")
	return true
}