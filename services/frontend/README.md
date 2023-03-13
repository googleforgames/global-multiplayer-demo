# Application API keys

CLIENT_ID and SECRET_ID need to be generated and fetched from https://console.cloud.google.com/apis/credentials (OAuth 2.0 Client IDs)

For the JWT_KEY, this can be any arbitrary string, but has to be consistent between deployments.
s
# For Local development

Please create a .env file in the same with the following variables:

```bash
CLIENT_ID=<CLIENT_ID>.apps.googleusercontent.com
CLIENT_SECRET=<CLIENT_SECRET>
LISTEN_PORT=8081
CALLBACK_HOSTNAME=http://localhost:8081/callback
CLIENT_LAUNCHER_PORT=8082
PROFILE_SERVICE=http://localhost:8080
PING_SERVICE=http://localhost:8083
JWT_KEY=<JWT_KEY>
API_ACCESS_KEY=<API_KEY_FOR_GAMESERVER_TO_ACCESS_FRONTEND>
LOCAL_OPENMATCH_SERVER_OVERRIDE_HOST=127.0.0.1 # in case you are testing local gameserver build and have no connection to agones nor openmatch
LOCAL_OPENMATCH_SERVER_OVERRIDE_PORT=7777 # port of the local gameserver
```

* `LISTEN_PORT` is the local port for this Docker container
* `CALLBACK_HOSTNAME` is the full URL to which authentication provider will redirect to. Should be a hostname registered in the https://console.cloud.google.com/apis/credentials (OAuth 2.0 Client IDs) or similar. Points back to this application
* `CLIENT_LAUNCHER_PORT` is the port that the launcher uses. There shouldn't be any reason to change this value.

# Building locally

`make build`

Binary will be generated in the `build/` folder

# Building container

`make build-docker`

To use container locally:

`docker run --env-file .env -p 8081:8081 frontend-api`

# Running in production

Make sure the following environment variables are set

```bash
CLIENT_ID=<CLIENT_ID>.apps.googleusercontent.com
CLIENT_SECRET=<CLIENT_SECRET>
LISTEN_PORT=8081
CLIENT_LAUNCHER_PORT=8082
PROFILE_SERVICE=<PROFILE_SERVICE_ENDPOINT>
PING_SERVICE=<PING_SERVICE_ENDPOINT>
JWT_KEY=<JWT_KEY>
API_ACCESS_KEY=<API_KEY_FOR_GAMESERVER_TO_ACCESS_FRONTEND>
```
