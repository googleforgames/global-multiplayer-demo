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

package mmf

import (
	"fmt"
	"log"
	"sort"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"open-match.dev/open-match/pkg/matchfunction"
	"open-match.dev/open-match/pkg/pb"
)

const (
	matchName = "match-skill-and-latency"

	ticketsPerMatch = 4
)

// Run is this match function's implementation of the gRPC call defined in api/matchfunction.proto.
func (s *MatchFunctionService) Run(req *pb.RunRequest, stream pb.MatchFunction_RunServer) error {
	// Fetch tickets for the single pool specified in the Match Profile.
	p := req.GetProfile()
	log.Printf("Generating proposals for profile %v", p.GetName())

	tickets, err := matchfunction.QueryPool(stream.Context(), s.queryServiceClient, p.GetPools()[0])
	if err != nil {
		log.Printf("Failed to query tickets pool %q: %v", p.GetPools()[0], err)
		return err
	}

	// Generate proposals.
	idPrefix := fmt.Sprintf("profile-%v-time-%v", p.GetName(), time.Now().Format("2006-01-02T15:04:05.00"))
	proposals, err := makeMatches(p.GetName(), idPrefix, tickets)
	if err != nil {
		log.Printf("Failed to generate matches: %v", err)
		return err
	}

	// Stream the generated proposals back to Open Match.
	log.Printf("Streaming %v proposals to Open Match", len(proposals))
	for _, proposal := range proposals {
		if err := stream.Send(&pb.RunResponse{Proposal: proposal}); err != nil {
			log.Printf("Failed to stream proposals to Open Match: %v", err)
			return err
		}
	}

	return nil
}

// Find all matches for the given profile.
func makeMatches(profileName, idPrefix string, tickets []*pb.Ticket) ([]*pb.Match, error) {
	if len(tickets) < ticketsPerMatch {
		return nil, nil
	}

	// Score each ticket and sort by score.
	ticketScores := make(map[string]float64) // map of Ticket.Id -> fitness score
	for _, ticket := range tickets {
		ticketScores[ticket.Id] = score(ticket.SearchFields.DoubleArgs["skill"], ticket.SearchFields.DoubleArgs["latency-"+profileName])
	}
	sort.Slice(tickets, func(i, j int) bool {
		return ticketScores[tickets[i].Id] > ticketScores[tickets[j].Id]
	})

	var matches []*pb.Match
	count := 0
	for len(tickets) >= ticketsPerMatch {
		matchTickets := tickets[:ticketsPerMatch]
		tickets = tickets[ticketsPerMatch:]

		var matchScore float64
		for _, ticket := range matchTickets {
			matchScore += ticketScores[ticket.Id]
		}

		eval, err := anypb.New(&pb.DefaultEvaluationCriteria{Score: matchScore})
		if err != nil {
			log.Printf("Failed to marshal DefaultEvaluationCriteria into anypb: %v", err)
			return nil, fmt.Errorf("Failed to marshal DefaultEvaluationCriteria into anypb: %w", err)
		}

		matches = append(matches, &pb.Match{
			MatchId:       fmt.Sprintf("%s-%d", idPrefix, count),
			MatchProfile:  profileName,
			MatchFunction: matchName,
			Tickets:       matchTickets,
			Extensions:    map[string]*anypb.Any{"evaluation_input": eval},
		})
		count++
	}
	return matches, nil
}

func score(skill, latency float64) float64 {
	// skill is kill/death, latency is in milliseconds - aggregate in a way that the higher the score, the better
	// (so we subtract latency, since lower latency is better).
	return skill - (latency / 1000.0)
}
