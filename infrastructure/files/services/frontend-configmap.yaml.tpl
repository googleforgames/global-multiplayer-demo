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

apiVersion: v1
kind: ConfigMap
metadata:
  name: frontend-service
data:
  CLIENT_ID: ${client_id}
  CLIENT_SECRET: ${client_secret}
  LISTEN_PORT: "8080"
  CLIENT_LAUNCHER_PORT: "8082"
  PROFILE_SERVICE: http://profile
  PING_SERVICE: http://ping-discovery
  JWT_KEY: ${jwt_key}
