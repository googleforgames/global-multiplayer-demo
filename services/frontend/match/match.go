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

// Match creates/watches Open Match tickets for a PlayRequest
package match

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	om "open-match.dev/open-match/pkg/pb"

	"github.com/googleforgames/global-multiplayer-demo/services/frontend-api/models"
)

const (
	// The endpoint for the Open Match Frontend service.
	omFrontendEndpoint = "open-match-frontend.open-match.svc.cluster.local:50504"
)

// Matcher interfaces with Open Match to create a ticket and watch for an assignment.
type Matcher struct {
	conn   *grpc.ClientConn
	client om.FrontendServiceClient
}

// NewMatcher returns a new Matcher. Close() should be deferred after a successful return.
func NewMatcher() (*Matcher, error) {
	conn, err := grpc.Dial(omFrontendEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Open Match: %w", err)
	}
	return &Matcher{client: om.NewFrontendServiceClient(conn)}, nil
}

// Close releases any resources of the Matcher.
func (m *Matcher) Close() {
	m.conn.Close()
}

// FindMatchingServer takes a PlayRequest, constructs an Open Match ticket, and waits for an assignment.
func (m *Matcher) FindMatchingServer(ctx context.Context, pr *models.PlayRequest) (*models.OMServerResponse, error) {
	log.Printf("Creating Open Match ticket for /play request: %#v", pr)

	req := &om.CreateTicketRequest{
		Ticket: makeTicket(pr),
	}
	resp, err := m.client.CreateTicket(ctx, req)
	if err != nil {
		log.Printf("CreateTicket failed for ticket %#v: %v", req.Ticket, err)
		return nil, fmt.Errorf("CreateTicket failed: %w", err)
	}
	tid := resp.Id
	log.Printf("Ticket %s: created: %v", tid, req.Ticket)

	stream, err := m.client.WatchAssignments(ctx, &om.WatchAssignmentsRequest{TicketId: tid})
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("Ticket %s: WatchAssignments failed: %v", tid, err)
			return nil, fmt.Errorf("WatchAssignments failed: %w", err)
		}

		omsr, err := hostPortToModel(resp.Assignment.Connection)
		if err != nil {
			log.Printf("Ticket %s: can't parse connection: %v", tid, err)
			return nil, fmt.Errorf("can't parse connection: %w", err)
		}
		return omsr, nil
	}
}

func hostPortToModel(hostPort string) (*models.OMServerResponse, error) {
	pieces := strings.Split(hostPort, ":")
	if len(pieces) != 2 {
		return nil, fmt.Errorf("host:port %q has %d pieces, want 2", hostPort, len(pieces))
	}
	port, err := strconv.Atoi(pieces[1])
	if err != nil {
		return nil, fmt.Errorf("can't parse port of host:port %q as int: %w", hostPort, err)
	}
	return &models.OMServerResponse{
		IP:   pieces[0],
		Port: port,
	}, nil
}

func makeTicket(pr *models.PlayRequest) *om.Ticket {
	t := &om.Ticket{
		SearchFields: &om.SearchFields{
			DoubleArgs: map[string]float64{
				"skill": 0.0, // TODO: Add skill!
			},
		},
	}
	// TODO: validate against known regions
	for region, ping := range pr.PingByRegion {
		t.SearchFields.DoubleArgs["latency-"+region] = float64(ping)
	}
	return t
}
