//go:build integration

// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/stretchr/testify/assert"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"embed"
	"fmt"
	"github.com/googleforgames/global-multiplayer-demo/profile-service/models"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

//go:embed test_data/schema.sql
var SCHEMAFILE embed.FS

var TESTNETWORK = "game-sample-test"

// These integration tests run against the Spanner emulator. The emulator
// must be running and accessible prior to integration tests running.

type Emulator struct {
	testcontainers.Container
	Endpoint string
	Project  string
	Instance string
	Database string
}

type Service struct {
	testcontainers.Container
	Endpoint string
}

func teardown(ctx context.Context, emulator *Emulator, service *Service) {
	emulator.Terminate(ctx)
	service.Terminate(ctx)
}

func setupSpannerEmulator(ctx context.Context) (*Emulator, error) {
	req := testcontainers.ContainerRequest{
		Image:        "gcr.io/cloud-spanner-emulator/emulator:1.5.0",
		ExposedPorts: []string{"9010/tcp"},
		Networks: []string{
			TESTNETWORK,
		},
		NetworkAliases: map[string][]string{
			TESTNETWORK: []string{
				"emulator",
			},
		},
		Name:       "emulator",
		WaitingFor: wait.ForLog("gRPC server listening at"),
	}
	spannerEmulator, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Retrieve the container IP
	ip, err := spannerEmulator.Host(ctx)
	if err != nil {
		return nil, err
	}

	// Retrieve the container port
	port, err := spannerEmulator.MappedPort(ctx, "9010")
	if err != nil {
		return nil, err
	}

	// OS environment needed for setting up instance and database
	os.Setenv("SPANNER_EMULATOR_HOST", fmt.Sprintf("%s:%d", ip, port.Int()))

	var ec = Emulator{
		Container: spannerEmulator,
		Endpoint:  "emulator:9010",
		Project:   "test-project",
		Instance:  "test-instance",
		Database:  "test-database",
	}

	// Create instance
	err = setupInstance(ctx, ec)
	if err != nil {
		return nil, err
	}

	// Define the database and schema
	err = setupDatabase(ctx, ec)
	if err != nil {
		return nil, err
	}

	return &ec, nil
}

func setupInstance(ctx context.Context, ec Emulator) error {
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer instanceAdmin.Close()

	op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", ec.Project),
		InstanceId: ec.Instance,
		Instance: &instancepb.Instance{
			Config:      fmt.Sprintf("projects/%s/instanceConfigs/%s", ec.Project, "emulator-config"),
			DisplayName: ec.Instance,
			NodeCount:   1,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create instance %s: %v", fmt.Sprintf("projects/%s/instances/%s", ec.Project, ec.Instance), err)
	}
	// Wait for the instance creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance creation to finish failed: %v", err)
	}

	// The instance may not be ready to serve yet.
	if i.State != instancepb.Instance_READY {
		fmt.Printf("instance state is not READY yet. Got state %v\n", i.State)
	}
	fmt.Printf("Created emulator instance [%s]\n", ec.Instance)

	return nil
}

func setupDatabase(ctx context.Context, ec Emulator) error {
	// get schema statements from file
	schema, _ := SCHEMAFILE.ReadFile("test_data/schema.sql")

	// Remove trailing semi-colon/newline so Emulator can parse DDL statements
	schemaString := strings.TrimSuffix(string(schema), ";\n")

	schemaStatements := strings.Split(schemaString, ";")

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%s/instances/%s", ec.Project, ec.Instance),
		CreateStatement: "CREATE DATABASE `" + ec.Database + "`",
		ExtraStatements: schemaStatements,
	})
	if err != nil {
		fmt.Printf("Error: [%s]", err)
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		fmt.Printf("Error: [%s]", err)
		return err
	}

	fmt.Printf("Created emulator database [%s]\n", ec.Database)
	return nil

}

func setupService(ctx context.Context, ec *Emulator) (*Service, error) {
	var service = "profile-service"
	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:latest", service),
		Name:         service,
		ExposedPorts: []string{"80:80/tcp"}, // Bind to 80 on localhost to avoid not knowing about the container port
		Networks:     []string{TESTNETWORK},
		NetworkAliases: map[string][]string{
			TESTNETWORK: []string{
				service,
			},
		},
		Env: map[string]string{
			"SPANNER_PROJECT_ID":    ec.Project,
			"SPANNER_INSTANCE_ID":   ec.Instance,
			"SPANNER_DATABASE_ID":   ec.Database,
			"SERVICE_HOST":          "0.0.0.0",
			"SERVICE_PORT":          "80",
			"SPANNER_EMULATOR_HOST": ec.Endpoint,
		},
		WaitingFor: wait.ForLog("Listening and serving HTTP on 0.0.0.0:80"),
	}
	serviceContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Retrieve the container endpoint
	endpoint, err := serviceContainer.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}

	return &Service{
		Container: serviceContainer,
		Endpoint:  endpoint,
	}, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Setup the docker network so containers can talk to each other
	nr := testcontainers.NetworkRequest{
		Name:       TESTNETWORK,
		Attachable: true,
	}
	_, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: nr,
	})
	if err != nil {
		fmt.Printf("Error setting up docker test network: %s\n", err)
		os.Exit(1)
	}

	// Setup the emulator container and default instance/database
	spannerEmulator, err := setupSpannerEmulator(ctx)
	if err != nil {
		fmt.Printf("Error setting up emulator: %s\n", err)
		os.Exit(1)
	}

	// Run service
	service, err := setupService(ctx, spannerEmulator)
	if err != nil {
		fmt.Printf("Error setting up service: %s\n", err)
		os.Exit(1)
	}

	defer teardown(ctx, spannerEmulator, service)

	os.Exit(m.Run())
}

var test_player = models.Player{
	Player_google_id: "123456",
	Player_name:      "test player",
	Profile_image:    "default",
	Region:           "amer",
}

var test_stats = []string{
	"{'won': false, 'score': 100, 'kills': 5, 'deaths': 20 }",
	"{'won': true, 'score': 1000, 'kills': 20, 'deaths': 1 }"
}

func TestAddPlayers(t *testing.T) {
	pJson, err := json.Marshal(test_player)
	assert.Nil(t, err)

	bufferJson := bytes.NewBuffer(pJson)

	// Test adding non-existing players
	response, err := http.Post("http://localhost/players", "application/json", bufferJson)
	assert.Nil(t, err)

	assert.Equal(t, 201, response.StatusCode)

	// Validate response is the original player_google_id
	var data string
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	json.Unmarshal(body, &data)
	assert.Equal(t, test_player.Player_google_id, data)

	// Test adding same player, should be statuscode 400 since player exists from previous call
	response, err = http.Post("http://localhost/players", "application/json", bufferJson)
	assert.Nil(t, err)
	assert.Equal(t, 400, response.StatusCode)
}

func TestGetPlayers(t *testing.T) {
	// Get the testPlayer's data and validate response code (assuming the result was not empty)
	response, err := http.Get(fmt.Sprintf("http://localhost/players/%s", test_player.Player_google_id))
	if err != nil {
		t.Fatal(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	var pData models.Player
	json.Unmarshal(body, &pData)

	assert.Equal(t, test_player.Player_google_id, pData.Player_google_id)
	assert.Equal(t, test_player.Player_name, pData.Player_name)
	assert.Equal(t, test_player.Profile_image, pData.Profile_image)
	assert.Equal(t, test_player.Region, pData.Region)
	assert.Equal(t, 0, pData.Skill_level)
	assert.Equal(t, "U", pData.Tier)
}

func TestUpdatePlayer(t *testing.T) {
	test_player.region = "asia"

	pJson, err := json.Marshal(test_player)
	assert.Nil(t, err)

	bufferJson := bytes.NewBuffer(pJson)

	response, err := http.Put("http://localhost/players", "application/json", bufferJson)
	assert.Nil(t, err)

	assert.Equal(t, 200, response.StatusCode)

	assert.Equal(t, test_player.Player_google_id, pData.Player_google_id)
	assert.Equal(t, test_player.Player_name, pData.Player_name)
	assert.Equal(t, test_player.Profile_image, pData.Profile_image)
	assert.Equal(t, test_player.Region, pData.Region)
	assert.Equal(t, 0, pData.Skill_level)
	assert.Equal(t, "U", pData.Tier)
}

func TestUpdatePlayerStats(t *testing.T) {
	for _, new_stats := range test_stats {
		bufferJson := bytes.NewBuffer(new_stats)

		response, err := http.Put(fmt.Sprintf("http://localhost/players/%s/stats", test_player.Player_google_id), "application/json", bufferJson)
		assert.Nil(t, err)

		assert.Equal(t, 200, response.StatusCode)
	}
}

func TestGetPlayerStats(t *testing.T) {
	// Get the testPlayer's stats and validate response code (assuming the result was not empty)
	response, err := http.Get(fmt.Sprintf("http://localhost/players/%s/stats", test_player.Player_google_id))
	if err != nil {
		t.Fatal(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	var pData models.Player
	json.Unmarshal(body, &pData)

	assert.Equal(t, test_player.Player_google_id, pData.Player_google_id)
	assert.Equal(t, 2, pData.Stats.Games_played)
	assert.Equal(t, 1, pData.Stats.Games_won)
	assert.Equal(t, 1100, pData.Stats.Total_score)
	assert.Equal(t, 25, pData.Stats.Total_kills)
	assert.Equal(t, 21, pData.Stats.Total_deaths)
	assert.Equal(t, 1, pData.Skill_level)
}
