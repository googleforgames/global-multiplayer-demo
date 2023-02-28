# For Local development

Please create a .env file in the same with the following variables:

```bash
CLIENT_ID=<CLIENT_ID>.apps.googleusercontent.com
CLIENT_SECRET=<CLIENT_SECRET>
LISTEN_PORT=8081
CLIENT_LAUNCHER_PORT=8082
PROFILE_SERVICE=http://localhost:8080
PING_SERVICE=http://localhost:8083
JWT_KEY=<JWT_KEY>
```

CLIENT_ID and SECRET_ID need to be generated in the https://console.cloud.google.com/apis/credentials (OAuth 2.0 Client IDs)

# Building locally

`make build`

Binary will be generated in the `build/` folder

# Building container

`make build-docker`

To use container locally:

`docker run --env-file .env -p 8081:8081 frontend-api`
