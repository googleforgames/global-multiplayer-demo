# ./base/kustomization.yaml
helmCharts:
- name: open-match
  namespace: open-match
  repo: https://open-match.dev/chart/stable
  version: 1.6.0
  releaseName: open-match
  valuesInline:
    open-match-override:
      enabled: true
    open-match-customize:
      enabled: true
      evaluator:
        enabled: true
    open-match-core:
      redis:
        enabled: false
        # If open-match-core.redis.enabled is set to false, have Open Match components talk to this redis address instead.
        # Otherwise the default is set to the om-redis instance.
        hostname: ${redis_host}
        port: ${redis_port}
        pool:
          maxIdle: 500
          maxActive: 500
          idleTimeout: 0
          healthCheckTimeout: 300ms

resources:
  - open-match.yaml
