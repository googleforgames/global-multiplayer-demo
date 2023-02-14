# Profile Service

The Profile Service provides a REST API to interact with Cloud Spanner to manage Player Profiles. The service runs on GKE Autopilot.

## Prerequisites
Cloud Spanner must be set up using the infrastructure steps before this service will work.

Local testing requires Docker to be installed.

## Schema management

This service uses the [wrench migration tool](https://github.com/cloudspannerecosystem/wrench) to perform Cloud Spanner schema migrations.

Basic usage is as follows.

- Create the database with the initial schema:

```
export SPANNER_PROJECT_ID=your-project-id
export SPANNER_INSTANCE_ID=your-instance-id
export SPANNER_DATABASE_ID=your-database-id

wrench create --directory ../infrastructure/schema
```

- Apply migrations

```
wrench migrate up --directory ./schema
```

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
