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

##### Spanner Pipelines #####

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
      target_id = google_clouddeploy_target.spanner.target_id
    }
  }
}

##### Agones Pipelines #####

resource "google_clouddeploy_target" "agones" {
  for_each = var.game_gke_clusters

  location = var.clouddeploy_config.location
  name     = "agones-deploy-target-${each.key}"

  annotations = {
    my_first_annotation = "agones-annotation-1"

    my_second_annotation = "agones-annotation-2"
  }

  description = "Global Game: Agones Deploy Target - ${each.key}"

  gke {
    cluster = data.google_container_cluster.game-demo-agones-gke[each.key].id
  }

  labels = {
    my_first_label = "global-game-demo"

    my_second_label = "agones"
  }

  project          = var.project
  require_approval = false

  depends_on = [google_project_service.project]
}

resource "google_clouddeploy_delivery_pipeline" "agones" {
  location = var.clouddeploy_config.location
  name     = "agones-deploy-pipeline"

  annotations = {
    my_first_annotation = "agones-annotation-1"

    my_second_annotation = "agones-annotation-2"
  }

  description = "Global Game: Agones Deploy Pipeline"

  labels = {
    my_first_label = "global-game-demo"

    my_second_label = "agones"
  }

  project = var.project

  serial_pipeline {
    dynamic "stages" {
      for_each = var.game_gke_clusters
      content {
        target_id = google_clouddeploy_target.agones[stages.key].target_id
	      profiles = [stages.key]
      }
    }
  }
}

##### Open Match Pipelines #####

resource "google_clouddeploy_target" "open-match-target" {
  location = var.open-match_gke_config.location
  name     = "global-game-open-match-deploy-target"

  annotations = {
    my_first_annotation = "open-match-annotation-1"

    my_second_annotation = "open-match-annotation-2"
  }

  description = "Global Game: Open Match Deploy Target"

  gke {
    cluster = data.google_container_cluster.game-demo-open-match-gke.id
  }

  labels = {
    my_first_label = "global-game-demo"

    my_second_label = "open-match"
  }

  project          = var.project
  require_approval = false

  depends_on = [google_project_service.project]
}

resource "google_clouddeploy_delivery_pipeline" "open-match" {
  location = var.open-match_gke_config.location
  name     = "global-game-open-match-deploy-pipeline"

  annotations = {
    my_first_annotation = "open-match-annotation-1"

    my_second_annotation = "open-match-annotation-2"
  }

  description = "Global Game: Open Match Deploy Pipeline"

  labels = {
    my_first_label = "global-game-demo"

    my_second_label = "open-match"
  }

  project = var.project

  serial_pipeline {
    stages {
      target_id = google_clouddeploy_target.open-match-target.target_id
    }
  }
}

##### Cloud Deploy IAM #####

resource "google_project_iam_member" "clouddeploy-container" {
  project = var.project
  role    = "roles/container.developer"
  member  = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"

  depends_on = [google_project_service.project]
}

resource "google_project_iam_member" "clouddeploy-build" {
  project = var.project
  role    = "roles/cloudbuild.workerPoolUser"
  member  = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"

  depends_on = [google_project_service.project]
}

resource "google_project_iam_member" "clouddeploy-logs" {
  project = var.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"

  depends_on = [google_project_service.project]
}

resource "google_project_iam_member" "clouddeploy-jobrunner" {
  project = var.project
  role    = "roles/clouddeploy.jobRunner"
  member  = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"

  depends_on = [google_project_service.project]
}
