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

##### Services Pipelines #####

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
          project = var.project

          # Spanner config
          spanner_service_account = google_service_account.spanner-sa.email
          spanner_instance_id     = google_spanner_instance.global-game-spanner.name
          spanner_database_id     = google_spanner_database.spanner-database.name

          # Ping Service config
          ping_service_account = google_service_account.ping_discovery_sa.email

          # Frontend config
          frontend_client_id         = var.frontend-service.client_id
          frontend_client_secret     = var.frontend-service.client_secret
          frontend_jwt_key           = var.frontend-service.jwt_key
          frontend_service_address   = google_compute_address.frontend-service.address
          frontend_callback_hostname = "http://${google_compute_address.frontend-service.address}.sslip.io/callback"

          # Open Match config
          players_per_match = var.open-match-matchfunction.players_per_match
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

  deploy_parameters = {
    "agones.allocator.labels.region" = each.value.region
  }

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
      }
    }
  }
}

resource "google_clouddeploy_automation" "agones-gke" {
  name              = "agones-deploy-automation-${each.key}"
  project           = var.project
  location          = var.clouddeploy_config.location
  delivery_pipeline = google_clouddeploy_delivery_pipeline.agones-gke.name
  service_account   = google_service_account.cloudbuild-sa.email
  for_each          = merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters)

  description = "Agones Deploy Automation - ${each.key}"
  selector {
    targets {
      id = google_clouddeploy_target.agones-gke[each.key].target_id
    }
  }

  suspended = false
  rules {
    promote_release_rule {
      id                    = "promote-release"
      wait                  = var.clouddeploy_config.pipeline_promotion_wait
      destination_target_id = "@next"
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

resource "google_clouddeploy_automation" "gameservers_gke" {
  name              = "gameserver-deploy-automation-${each.key}"
  project           = var.project
  location          = var.clouddeploy_config.location
  delivery_pipeline = google_clouddeploy_delivery_pipeline.gameservers_gke.name
  service_account   = google_service_account.cloudbuild-sa.email
  for_each          = local.targets_by_region

  description = "Gameserver Deploy Automation - ${each.key}"
  selector {
    targets {
      id = google_clouddeploy_target.agones_regional_targets[each.key].target_id
    }
  }

  suspended = false
  rules {
    promote_release_rule {
      id                    = "promote-release"
      wait                  = var.clouddeploy_config.pipeline_promotion_wait
      destination_target_id = "@next"
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
      deploy_parameters {
        values = {
          # values are passed into Cloud Deploy Helm Charts by convention.
          "open-match-core.redis.hostname" = google_redis_instance.open-match.host
          "open-match-core.redis.port"     = google_redis_instance.open-match.port
        }
      }
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
