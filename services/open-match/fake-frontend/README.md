This folder provides a fake frontend that generates fake tickets for partial testing
of the Open Match services in services/open-match. Note that if the Director is tied
to a real Agones Fleet, this will generate arbitrary load, which could be bad.

The fake frontend is always deployed, but does nothing unless you change the configuration:

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
