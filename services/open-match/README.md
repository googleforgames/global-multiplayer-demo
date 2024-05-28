# Open Match Integration

## Ticket Format

The Match Function expects the following `SearchFields` on every ticket:
* `skill` is a `float64` that represents the aggregate skill of the player.
* `latency-$REGION` is the ping time of the player to `$REGION` in milliseconds, for each `$REGION` configured.

## Match Function

Our goal with the Match Function is to demonstrate something rudimentary but still interesting: Match
players by skill and by latency to a given region.

To align player latencies, the Match Function uses a separate [`MatchProfile`](https://pkg.go.dev/open-match.dev/open-match@v1.7.0/pkg/pb#MatchProfile) per region, but each regional `MatchProfile` evaluates every incoming ticket.

For each regional `MatchProfile`, we:
* Score each ticket based roughly on `skill-latency_to_region`, i.e. higher skill is better, lower latency to that region is better.
* Sort the incoming tickets by score
* Create 3 ticket matches from the sorted tickets, thereby grouping scores
* Assign a match score that is simply the sum of the scores of each ticket, for use by the [Default Evaluator](https://open-match.dev/site/docs/tutorials/defaultevaluator/).

The [Evaluator](https://open-match.dev/site/docs/guides/evaluator/) (part of Open Match Core) then chooses
a match from the overlapping matches returned by the different profiles.

## Director

The Director allocates a GameServer from an GKE Standard/Autopilot and Agones cluster hosted in the target region for a 
given set of match player's latencies.

It does this by providing the `region` HTTP header to an Anthos Service Mesh Allocation Service - where the `region` 
header will route the allocation request to one of the Agones GKE clusters in that region.

## Credit

This integration is based on the [Open Match Matchmaker 101 tutorial](https://open-match.dev/site/docs/tutorials/matchmaker101/frontend/) [(source)](https://github.com/googleforgames/open-match/tree/release-1.7/tutorials).
