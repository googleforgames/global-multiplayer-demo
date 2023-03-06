## Generic game launcher with "Sign-in with google" authentication.

To run the launcher locally:

```shell
go run main.go
```

To build locally after installing necessary dependencies just run:

```shell
go build .
```

Move the built binary together with `app.ini` and `assets/` to the game client folder. Launch there.

`app.ini` contains the configuration endpoint for the Frontend API as well as executable names for the game client.

If you want to fully package the launcher:

For prerequisites check here:
https://developer.fyne.io/started/

For packaging check here:
https://developer.fyne.io/started/packaging
