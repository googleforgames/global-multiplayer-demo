// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

resource "google_artifact_registry_repository" "container_registry" {
  repository_id = "global-game-images"
  location      = var.artifact_registry_config.location
  description   = "Repository for container images for the global game"
  format        = "Docker"

  labels = {
    "environment" = var.resource_env_label
  }

  depends_on = [google_project_service.project]
}
