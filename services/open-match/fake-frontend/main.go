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

// The Frontend in this tutorial continuously creates Tickets in batches in Open Match.

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"open-match.dev/open-match/pkg/pb"
)

const (
	// The endpoint for the Open Match Frontend service.
	omFrontendEndpoint = "open-match-frontend.open-match.svc.cluster.local:50504"
)

func main() {
	period, ticketsPerPeriod := getConfig()
	if period == 0 || ticketsPerPeriod == 0 {
		log.Printf("Disabled, sleeping")
		for {
			time.Sleep(24 * time.Hour)
		}
	}

	// Connect to Open Match Frontend.
	conn, err := grpc.Dial(omFrontendEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Open Match, got %v", err)
	}

	defer conn.Close()
	fe := pb.NewFrontendServiceClient(conn)
	for range time.Tick(period) {
		for i := 0; i <= ticketsPerPeriod; i++ {
			req := &pb.CreateTicketRequest{
				Ticket: makeTicket(),
			}

			log.Printf("Creating Ticket %v", req.Ticket)
			resp, err := fe.CreateTicket(context.Background(), req)
			if err != nil {
				log.Printf("Failed to Create Ticket, got %s", err.Error())
				continue
			}

			log.Println("Ticket created successfully, id:", resp.Id)
			go deleteOnAssign(fe, resp)
		}
	}
}

func getConfig() (time.Duration, int) {
	period, err := time.ParseDuration(os.Getenv("FAKE_FRONTEND_PERIOD"))
	if err != nil {
		log.Printf("FAKE_FRONTEND_PERIOD not set or set to invalid duration: %v", err)
		period = 0
	}
	tickets, err := strconv.Atoi(os.Getenv("FAKE_FRONTEND_TICKETS_PER_PERIOD"))
	if err != nil {
		log.Printf("FAKE_FRONTEND_TICKETS_PER_PERIOD not set or invalid: %v", err)
		tickets = 0
	}
	log.Printf("Configured to generate %d tickets every %v", tickets, period)
	return period, tickets
}

// deleteOnAssign fetches the Ticket state periodically and deletes the Ticket
// once it has an assignment.
func deleteOnAssign(fe pb.FrontendServiceClient, t *pb.Ticket) {
	for {
		got, err := fe.GetTicket(context.Background(), &pb.GetTicketRequest{TicketId: t.GetId()})
		if err != nil {
			log.Fatalf("Failed to Get Ticket %v, got %s", t.GetId(), err.Error())
		}

		if got.GetAssignment() != nil {
			log.Printf("Ticket %v got assignment %v", got.GetId(), got.GetAssignment())
			break
		}

		time.Sleep(time.Second * 1)
	}

	_, err := fe.DeleteTicket(context.Background(), &pb.DeleteTicketRequest{TicketId: t.GetId()})
	if err != nil {
		log.Fatalf("Failed to Delete Ticket %v, got %s", t.GetId(), err.Error())
	}
}

// Ticket generates a Ticket with data using the package configuration.
func makeTicket() *pb.Ticket {
	return &pb.Ticket{
		SearchFields: &pb.SearchFields{
			DoubleArgs: map[string]float64{
				"skill":                2 * rand.Float64(),
				"latency-us-central1":  50.0 * rand.ExpFloat64(),
				"latency-europe-west1": 50.0 * rand.ExpFloat64(),
				"latency-asia-east1":   50.0 * rand.ExpFloat64(),
			},
		},
	}
}
