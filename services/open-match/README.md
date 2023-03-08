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
* Create 4 ticket matches from the sorted tickets, thereby grouping scores
* Assign a match score that is simply the sum of the scores of each ticket, for use by the [Default Evaluator](https://open-match.dev/site/docs/tutorials/defaultevaluator/).

The [Evaluator](https://open-match.dev/site/docs/guides/evaluator/) (part of Open Match Core) then chooses
a match from the overlapping matches returned by the different profiles.

## Director

TBD, write more when the Director is wired to Agones.

## Fake Frontend

The `openmatch-fake-frontend` deployment generates fake tickets for testing of the Open Match services
in thie directory. When enabled, the fake frontend will generate a batch of tickets every period. For
each ticket, the fake frontend waits on an Open Match assignment and then deletes the ticket.

Note that if the Director is tied to a real Agones Fleet, this will generate arbitrary game server
allocations. You may want to use a game server binary that exits after a specific period of time,
otherwise, when enabled, the fake frontend will continue to allocate game servers.

The fake frontend is always deployed, but does nothing unless configured:

```shell
# Edit the `period` and `tickets_per_period` to what you want, e.g. "20s" and "20"
# The fake frontend will generate `tickets_per_period` tickets every `period`, where
# period is expressed as a Go Duration (e.g. "20s", "1m", etc.)
kubectl edit configmap/open-match-fake-frontend

# Restart the deployment to pick up the config
kubectl rollout restart deployment/open-match-fake-frontend

# Verify load has started
kubectl logs deployment/open-match-fake-frontend
```

## Credit

This integration is based on the [Open Match Matchmaker 101 tutorial](https://open-match.dev/site/docs/tutorials/matchmaker101/frontend/) [(source)](https://github.com/googleforgames/open-match/tree/release-1.7/tutorials).
