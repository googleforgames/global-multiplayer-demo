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

resource "google_clouddeploy_target" "spanner" {
  location = var.spanner_gke_config.location
  name     = "global-game-spanner-deploy-target"

  annotations = {
    my_first_annotation = "spanner-annotation-1"

    my_second_annotation = "spanner-annotation-2"
  }

  description = "Global Game: Spanner Deploy Target"

  gke {
    cluster = data.google_container_cluster.game-demo-spanner-gke.id
  }

  labels = {
    my_first_label = "global-game-demo"

    my_second_label = "spanner"
  }

  project          = var.project
  require_approval = false

  depends_on = [google_project_service.project]
}

resource "google_clouddeploy_delivery_pipeline" "spanner" {
  location = var.spanner_gke_config.location
  name     = "global-game-spanner-deploy-pipeline"

  annotations = {
    my_first_annotation = "spanner-annotation-1"

    my_second_annotation = "spanner-annotation-2"
  }

  description = "Global Game: Spanner Deploy Pipeline"

  labels = {
    my_first_label = "global-game-demo"

    my_second_label = "spanner"
  }

  project = var.project

  serial_pipeline {
    stages {
      profiles  = ["spanner-profile-one"]
      target_id = google_clouddeploy_target.spanner.target_id
    }
  }
}
