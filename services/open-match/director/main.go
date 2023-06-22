// Copyright 2023 Google LLC
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
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	allocation "github.com/googleforgames/global-multiplayer-demo/services/open-match/director/agones/swagger"
	"google.golang.org/grpc"
	"open-match.dev/open-match/pkg/pb"
)

// The Director in this tutorial continuously polls Open Match for the Match
// Profiles and makes random assignments for the Tickets in the returned matches.

const (
	// The endpoint for the Open Match Backend service.
	omBackendEndpoint = "open-match-backend.open-match.svc.cluster.local:50505"

	// The Host and Port for the Match Function service endpoint.
	functionHostName       = "open-match-matchfunction.open-match.svc.cluster.local"
	functionPort     int32 = 50502

	// Agones Allocation Service base path
	agonesAllocationService = "http://agones-allocator.agones-system.svc.cluster.local:8000"

	// Namespace to allocate from
	gameNamespace = "default"
)

// TODO: This should be an environment variable.
var regions = []string{"us-central1", "europe-west1", "asia-east1"}

func main() {
	ctx := context.Background()

	// Connect to Open Match Backend.
	conn, err := grpc.Dial(omBackendEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Open Match Backend, got %s", err.Error())
	}

	defer conn.Close()
	be := pb.NewBackendServiceClient(conn)

	// Create a client per region, using "region" header to route via the ASM
	// VirtualService. (Each of these clients is accessing the same endpoint,
	// but using a different header.)
	aas := make(map[string]*allocation.APIClient)
	for _, region := range regions {
		aas[region] = allocation.NewAPIClient(&allocation.Configuration{
			BasePath:      agonesAllocationService,
			DefaultHeader: map[string]string{"region": region},
			UserAgent:     "global-multiplayer-demo/open-match/director",
		})
	}

	// Generate the profiles to fetch matches for.
	profiles := generateProfiles()
	log.Printf("Fetching matches for %v profiles", len(profiles))

	for range time.Tick(time.Second * 5) {
		// Fetch matches for each profile and make random assignments for Tickets in
		// the matches returned.
		var wg sync.WaitGroup
		for _, p := range profiles {
			wg.Add(1)
			go func(wg *sync.WaitGroup, p *pb.MatchProfile) {
				defer wg.Done()
				matches, err := fetch(be, p)
				if err != nil {
					log.Printf("Failed to fetch matches for profile %v, got %s", p.GetName(), err.Error())
					return
				}

				log.Printf("Generated %v matches for profile %v", len(matches), p.GetName())
				for _, match := range matches {
					assignMatch(ctx, be, aas[match.GetMatchProfile()], match)
				}
			}(&wg, p)
		}

		wg.Wait()
	}
}

func fetch(be pb.BackendServiceClient, p *pb.MatchProfile) ([]*pb.Match, error) {
	req := &pb.FetchMatchesRequest{
		Config: &pb.FunctionConfig{
			Host: functionHostName,
			Port: functionPort,
			Type: pb.FunctionConfig_GRPC,
		},
		Profile: p,
	}

	stream, err := be.FetchMatches(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error on fetch match request from backend: %w", err)
	}

	var result []*pb.Match
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error on fetch match stream recieving: %w", err)
		}

		result = append(result, resp.GetMatch())
	}

	return result, nil
}

// assignMatch assigns `match`. If we fail, abandon the tickets - we'll catch it next loop.
func assignMatch(ctx context.Context, be pb.BackendServiceClient, aas *allocation.APIClient, match *pb.Match) {
	aar, _, err := aas.AllocationServiceApi.Allocate(ctx, allocation.AllocationAllocationRequest{Namespace: gameNamespace})
	if err != nil {
		if swErr, ok := err.(allocation.GenericSwaggerError); ok {
			log.Printf("Could not allocate game server from Agones for match %s: %s: %s", match.GetMatchId(), swErr.Error(), swErr.Body())
		} else {
			log.Printf("Could not allocate game server from Agones for match %s (generic error): %v", match.GetMatchId(), err)
		}
		return
	}
	if len(aar.Ports) < 1 {
		log.Printf("Expecting at least one port from Agones allocation. Response: %v", aar)
		return
	}
	// blindly assume port[0] is the one we want
	conn := fmt.Sprintf("%s:%d", aar.Address, aar.Ports[0].Port)
	log.Printf("Allocated %s for match %s. Payload: %v", conn, match.GetMatchId(), aar)

	if err := assignConnToTickets(be, conn, match.GetTickets()); err != nil {
		log.Printf("Could not assign connection %s to match %s: %v", conn, match.GetMatchId(), err)
	}

	log.Printf("Assigned %s to match %s", conn, match.GetMatchId())
}

func assignConnToTickets(be pb.BackendServiceClient, conn string, tickets []*pb.Ticket) error {
	ticketIDs := []string{}
	for _, t := range tickets {
		ticketIDs = append(ticketIDs, t.Id)
	}

	req := &pb.AssignTicketsRequest{
		Assignments: []*pb.AssignmentGroup{
			{
				TicketIds: ticketIDs,
				Assignment: &pb.Assignment{
					Connection: conn,
				},
			},
		},
	}

	_, err := be.AssignTickets(context.Background(), req)
	return err
}

func generateProfiles() []*pb.MatchProfile {
	var profiles []*pb.MatchProfile
	for _, region := range regions {
		profiles = append(profiles, &pb.MatchProfile{
			Name: region,
			Pools: []*pb.Pool{{
				Name: region,
			}},
		})
	}
	return profiles
}
