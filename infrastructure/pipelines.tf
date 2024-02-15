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

resource "google_clouddeploy_target" "services_deploy_target" {
  location    = var.services_gke_config.location
  name        = "global-game-services-target"
  description = "Global Game: Backend Services Deploy Target"

  gke {
    cluster = data.google_container_cluster.services-gke.id
  }

  project          = var.project
  require_approval = false

  labels = {
    "environment" = var.resource_env_label
  }

  depends_on = [google_project_service.project]
}

resource "google_clouddeploy_delivery_pipeline" "services_pipeline" {
  location = var.services_gke_config.location
  name     = "global-game-services"

  description = "Global Game: Backend Services Pipeline"

  project = var.project

  labels = {
    "environment" = var.resource_env_label
  }

  serial_pipeline {
    stages {
      target_id = google_clouddeploy_target.services_deploy_target.target_id
      deploy_parameters {
        values = {
          project               = var.project
          spanner_service_email = google_service_account.spanner-sa.email
          spanner_instance_id   = google_spanner_instance.global-game-spanner.name
          spanner_database_id   = google_spanner_database.spanner-database.name
        }
      }
    }
  }
}

##### Agones Pipelines #####

resource "google_clouddeploy_target" "agones-gke" {
  for_each = merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters)

  location = var.clouddeploy_config.location
  name     = "${each.value.short_name}-agones-deploy"

  description = "Global Game: Agones Deploy Target - ${each.key}"

  gke {
    cluster = data.google_container_cluster.game-demo-agones[each.key].id
  }

  labels = {
    "environment" = var.resource_env_label
  }

  project          = var.project
  require_approval = false

  depends_on = [google_project_service.project]
}

resource "google_clouddeploy_delivery_pipeline" "agones-gke" {
  location    = var.clouddeploy_config.location
  name        = "agones-deploy-pipeline"
  description = "Global Game: Agones Deploy Pipeline"

  labels = {
    "environment" = var.resource_env_label
  }

  project = var.project

  serial_pipeline {
    dynamic "stages" {
      for_each = merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters)
      content {
        target_id = google_clouddeploy_target.agones-gke[stages.key].target_id
        profiles  = [stages.key]
      }
    }
  }
}

##### Game Server Pipeline

# Sort as a map of regions, with a list of Cloud Deploy targets within them.
# Amusingly, it comes out in exactly the order we want, so no need to get extra fancy.
locals {
  targets_by_region = {
    for name, cluster in merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters) : cluster.region =>
    google_clouddeploy_target.agones-gke[name].target_id...
  }
}

resource "google_clouddeploy_target" "agones_regional_targets" {
  for_each = local.targets_by_region
  provider = google-beta
  name     = each.key

  location    = var.clouddeploy_config.location
  description = "Global Game: Agones Game Servers - ${each.key}"
  labels = {
    "environment" = var.resource_env_label
  }

  project          = var.project
  require_approval = false

  multi_target {
    target_ids = local.targets_by_region[each.key]
  }

  depends_on = [google_project_service.project]
}

resource "google_clouddeploy_delivery_pipeline" "gameservers_gke" {
  location    = var.clouddeploy_config.location
  name        = "global-game-agones-gameservers"
  description = "Global Game: Agones GameServer Deploy Pipeline"
  provider    = google-beta

  labels = {
    "environment" = var.resource_env_label
  }

  project = var.project

  serial_pipeline {
    dynamic "stages" {
      for_each = local.targets_by_region
      content {
        target_id = google_clouddeploy_target.agones_regional_targets[stages.key].target_id
      }
    }
  }
}


##### Open Match Pipelines #####

resource "google_clouddeploy_target" "open-match-target" {
  location = var.services_gke_config.location
  name     = "global-game-open-match-target"

  description = "Global Game: Open Match Deploy Target"

  gke {
    cluster = data.google_container_cluster.services-gke.id
  }

  labels = {
    "environment" = var.resource_env_label
  }

  project          = var.project
  require_approval = false

  depends_on = [google_project_service.project]
}

resource "google_clouddeploy_delivery_pipeline" "open-match" {
  location = var.services_gke_config.location
  name     = "global-game-open-match"

  description = "Global Game: Open Match Deploy Pipeline"
  project     = var.project

  labels = {
    "environment" = var.resource_env_label
  }

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
