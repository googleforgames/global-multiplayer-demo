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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/ini.v1"
)

type GoogleOauthToken struct {
	AccessToken  string
	RefreshToken string
	Expiry       string
	TokenType    string
	IdToken      string
}

type UserInfo struct {
	Id    string `json:"id"`
	Sub   string `json:"sub"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	myToken        string
	myRefreshToken string
	myApp          fyne.App
	myWindow       fyne.Window
	iniCfg         *ini.File
)

func main() {
	// Load configs
	var err error
	iniCfg, err = ini.Load("app.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// Callback handling from the frontend api
	http.HandleFunc("/callback", handleGoogleCallback)
	go func() {
		fmt.Println("Google for Games Launcher is listening for callbacks on :" + iniCfg.Section("").Key("callback_listen_port").String())
		fmt.Println(http.ListenAndServe(":"+iniCfg.Section("").Key("callback_listen_port").String(), nil))
	}()

	// UI
	myApp = app.New()
	myWindow = myApp.NewWindow("Google for Games Launcher")
	myWindow.Resize(fyne.NewSize(320, 260))

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
	myWindow.ShowAndRun()
}

func handleGoogleCallback(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusInternalServerError)

			fmt.Fprintf(rw, "{\"error\": \"%s\"}", err)
			log.Println("panic occurred:", err)
		}
	}()

	// Save my token
	myToken = req.FormValue("access_token")
	if len(myToken) == 0 {
		panic("No token received!")
	}
	myRefreshToken = req.FormValue("refresh_token")
	if len(myRefreshToken) == 0 {
		panic("No refresh received!")
	}

	// Update UI with profile info and launch game button
	myProfile := getProfileInfo()
	fmt.Printf("My name is " + myProfile.Name)

	image := canvas.NewImageFromFile("assets/header.png")
	image.FillMode = canvas.ImageFillContain

	label1 := widget.NewLabel(fmt.Sprintf("Welcome %s!", myProfile.Name))
	label1.Alignment = fyne.TextAlignCenter
	label2 := widget.NewLabel(fmt.Sprintf("Are you ready to play again?!"))
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

	// Close the browser window
	closeScript := `<script> 
		setTimeout("window.close()",3000) 
	</script>
	<p>
		<h2>Authenticated successfully. Please return to your application. This tab will close in 3 seconds.</h2>
	</p>`
	fmt.Fprintf(rw, closeScript)
}

func handlePlay() {
	params := fmt.Sprintf("-token=%s -refresh_token=%s", myToken, myRefreshToken)

	// Get the binary file from the ini
	cmd := exec.Command(iniCfg.Section(runtime.GOOS).Key("binary").String(), params)
	fmt.Printf("Launching: %s %s\n", iniCfg.Section(runtime.GOOS).Key("binary").String(), params)

	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error: %s", err)
	}
}

func getProfileInfo() UserInfo {
	fmt.Printf("Getting profile info\n")

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + myToken)
	if response.StatusCode != 200 {
		panic("Unable to fetch user information. Expired token?")
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

	return result
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
