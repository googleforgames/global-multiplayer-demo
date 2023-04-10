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

serviceAccount: projects/$${PROJECT_ID}/serviceAccounts/cloudbuild-cicd@$${PROJECT_ID}.iam.gserviceaccount.com
steps:

  #
  # Building of the images
  #

  # Login to Github to be able to access Unreal container
  - name: gcr.io/cloud-builders/docker
    id: github-login
    script: |
      echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
    secretEnv:
      - CR_PAT

  # Build Linux Client & Server Binaries
  - name: ghcr.io/epicgames/unreal-engine:dev-slim-5.1.0
    id: build-linux-binaries
    args: [ "/home/ue4/UnrealEngine/Engine/Build/BatchFiles/RunUAT.sh", "BuildCookRun", 
            "--Server", "-ServerConfig=Development",
            "-Project=/workspace/Droidshooter.uproject", "-UTF8Output", "-NoDebugInfo", "-AllMaps",
            "-NoP4", "-Build", "-Cook", "-Stage", "-Pak", "-Package", "-Archive",
            "-ArchiveDirectory=/workspace/Packaged",
            "-Platform=Linux" ]

  # Publish Linux Server binaries to project Artifacts Repo
  - name: gcr.io/cloud-builders/docker
    id: publish-linux-server-container
    args: [ "build", ".", "-t", "$${_UNREAL_SERVER_IMAGE}" ]

  # Copy over Binaries to GCS
  - name: "gcr.io/cloud-builders/gcloud-slim"
    id: copy-binaries-to-gcs
    args: [
      "storage",
      "cp",
      "--recursive",
      "/workspace/Packaged",
      "gs://${CLIENT_BUCKET}"
  ]

  #
  # Deployment
  #

  # Run job for deployment of game server container to game GKE clusters
  - name: gcr.io/google.com/cloudsdktool/cloud-sdk
    id: deploy-linux-binaries-to-gke-clusters
    entrypoint: gcloud
    args:
      [
        "deploy", "releases", "create", "$${_RELEASE_NAME}",
        "--delivery-pipeline", "global-game-agones-gameservers",
        "--skaffold-file", "skaffold.yaml",
        "--images", "droidshooter-server=$${_UNREAL_SERVER_IMAGE}",
        "--region", "us-central1"
      ]

artifacts:
  images:
    - $${_REGISTRY}/droidshooter-server
substitutions:
  _UNREAL_SERVER_IMAGE: $${_REGISTRY}/droidshooter-server:$${BUILD_ID}
  _REGISTRY: us-docker.pkg.dev/$${PROJECT_ID}/global-game-images
  _RELEASE_NAME: rel-01
availableSecrets:
  secretManager:
    - versionName: projects/$${PROJECT_ID}/secrets/github-packages/versions/latest
      env: CR_PAT
options:
  dynamic_substitutions: true
  machineType: E2_HIGHCPU_32
  logging: CLOUD_LOGGING_ONLY
