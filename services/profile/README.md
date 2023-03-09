# Profile Service

The Profile Service provides a REST API to interact with Cloud Spanner to manage Player Profiles. The service runs on GKE Autopilot.

## API

<table>
    <thead>
        <td>Endpoint</td>
        <td>Input</td>
        <td>Return</td>
        <td>Description</td>
    </thead>
    <tbody>
        <tr>
            <td><code>GET /players/:player_id:</code></td>
            <td> None </td>
            <td>
                <pre>{
"player_google_id": "[string]",
"player_name": "[string]",
"profile_image": "[string]",
"region": "[string]"
}</pre>
            </td>
            <td>
                Retrieve information about the player
            </td>
        </tr>
    </tbody>
    <tbody>
        <tr>
            <td><code>GET /players/:player_id:/stats</code></td>
            <td> None </td>
            <td>
                <pre>{
"player_google_id": "[string]",
"stats": "[json]",
"skill_level": [int64],
"tier": "[string]" # currently unused
}</pre>
            </td>
            <td>
                Retrieve stats and skill info about the player
            </td>
        </tr>
    </tbody>
    <tbody>
        <tr>
            <td><code>POST /players</code></td>
            <td>
                <pre>{
"player_google_id": "[string]",
"player_name": "[string]",
"profile_image": "default", # Currently unused
"region": "[amer,eur,apac]"
}</pre>
            </td>
            <td>
                <pre>
                    "[player_google_id]"
                </pre>
            </td>
            <td>
                Create a new player
            </td>
        </tr>
    </tbody>
    <tbody>
        <tr>
            <td><code>PUT /players</code></td>
            <td>
                <pre>{
"player_google_id": "[string]",
"player_name": "[string]",
"profile_image": "default", # Currently unused
"region": "[amer,eur,apac]"
}</pre>
            </td>
            <td>
                <pre>{
"player_google_id": "[string]",
"player_name": "[string]",
"profile_image": "[string]",
"region": "[string]"
}</pre>
            </td>
            <td>
                Update player information
            </td>
        </tr>
    </tbody>
    <tbody>
        <tr>
            <td><code>PUT /players/:player_id:/stats</code></td>
            <td>
                <pre>{
"won": [true, false],
"score": [int64],
"kills": [int64],
"deaths": [int64]
}</pre>
            </td>
            <td>
                <pre>{
"player_google_id": "[string]",
"stats": "[json]",
"skill_level": [int64],
"tier": "[string]" # currently unused
}</pre>
            </td>
            <td>
                Update player stats
            </td>
        </tr>
    </tbody>


</table>

## Prerequisites
Cloud Spanner must be set up using the infrastructure steps before this service will work.

Local testing requires Docker to be installed.

## Schema management

This service uses the [Liquibase Spanner extension](https://github.com/cloudspannerecosystem/liquibase-spanner) to perform Cloud Spanner schema migrations.

## Local deployment

This service provides a Makefile to build a binary as well as run various tests. These tests require Docker to work.

The following commands should be run from the `./services/profile` directory.
Build a local binary:

```
make build
```

> ***NOTE:*** This build command does not build a docker container.

Build the docker container:

```
make build-docker
```

Run unit tests:

```
make test-unit
```

Run integration tests:

```
make test-integration
```

Cleanup binary and docker images:

```
make clean
```
