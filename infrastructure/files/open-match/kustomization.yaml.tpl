# Copyright 2023 Google LLC All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
  - agones-allocator-vs.yaml
