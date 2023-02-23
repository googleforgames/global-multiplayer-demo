// Copyright 2023 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

# Generate random string to endpoint service name to avoid the following error when deleting/recreating endpoint service
#   Error: googleapi: Error 400: Service global-game-us-central1-02.endpoints.global-game-sandbox.cloud.goog 
#   has been deleted and will be purged after 30 days. To reuse this service, please undelete the service
#   following https://cloud.google.com/service-infrastructure/docs/create-services#undeleting., failedPrecondition ...
resource "random_string" "endpoint_random_string" {
  length  = 4
  special = false
  upper   = false
}

resource "google_endpoints_service" "endpoints_service" {
  for_each     = var.game_gke_clusters
  service_name = "${each.key}-${random_string.endpoint_random_string.result}.endpoints.${var.project}.cloud.goog"
  grpc_config = templatefile(
    "${path.module}/files/agones/api_config.yaml.tpl", {
      service-name    = "${each.key}-${random_string.endpoint_random_string.result}.endpoints.${var.project}.cloud.goog"
      service-account = google_service_account.ae_sa.email
    }
  )
  protoc_output_base64 = filebase64("files/agones/agones_allocation_api_descriptor.pb")
}

resource "google_endpoints_service_iam_binding" "endpoints_service_binding" {
  for_each = var.game_gke_clusters

  service_name = google_endpoints_service.endpoints_service[each.key].service_name
  role         = "roles/servicemanagement.serviceController"
  members = [
    google_service_account.ae_sa.member
  ]
  depends_on = [google_project_service.allocator-service]
}

resource "google_service_account_iam_binding" "workload-identity-binding" {
  service_account_id = google_service_account.ae_sa.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:${var.project}.svc.id.goog[${var.allocation_endpoint.agones_namespace}/agones-allocator]",
  ]

  depends_on = [module.agones_gke_clusters]
}

resource "google_service_account" "ae_sa" {
  account_id   = "allocation-endpoint-esp-sa"
  display_name = "Service Account for Allocation Endpoint"
}

resource "google_service_account_key" "ae_sa_key" {
  service_account_id = google_service_account.ae_sa.name
}

resource "google_cloud_run_service_iam_binding" "binding" {
  for_each = var.game_gke_clusters

  service  = google_cloud_run_service.aep_cloud_run[each.key].name
  project  = google_cloud_run_service.aep_cloud_run[each.key].project
  location = google_cloud_run_service.aep_cloud_run[each.key].location
  role     = "roles/run.invoker"
  members = [
    "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com",
    "serviceAccount:${google_service_account.cloudbuild-sa.email}",
    "serviceAccount:${google_service_account.ae_sa.email}"
  ]
}


resource "google_cloud_run_service" "aep_cloud_run" {
  for_each = var.game_gke_clusters

  project  = var.project
  name     = "allocation-endpoint-proxy-${each.key}"
  location = each.value.region

  template {
    spec {
      container_concurrency = 80
      timeout_seconds       = 30
      containers {
        image = var.allocation_endpoint.proxy_image
        env {
          name = "CLUSTERS_INFO"
          value = templatefile(
            "${path.module}/files/agones/clusters_info.tpl", {
              name      = data.google_container_cluster.game-demo-agones-gke[each.key].name
              ip        = google_compute_address.allocation-endpoint[each.key].address
              weight    = var.allocation_endpoint.weight
              namespace = var.allocation_endpoint.agones_namespace
          })
        }
        env {
          name  = "AUDIENCE"
          value = "${each.key}-${random_string.endpoint_random_string.result}.endpoints.${var.project}.cloud.goog"
        }
        env {
          name = "SA_KEY"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.ae-sa-key.secret_id
              key  = "latest"
            }
          }
        }
        ports {
          container_port = 8080
          # this enables the http/2 support. h2c: https://cloud.google.com/run/docs/configuring/http2
          name = "h2c"
        }
        resources {
          limits = {
            "cpu"    = "2000m"
            "memory" = "256Mi"
          }
        }
      }
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1000"
        "autoscaling.knative.dev/minScale" = "0"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  metadata {
    annotations = {
      "run.googleapis.com/ingress"     = "all"
      "run.googleapis.com/client-name" = "terraform"
    }
  }

  lifecycle {
    ignore_changes = [
      # Ignore changes for the values set by GCP
      metadata[0].annotations,
      # This is currently not working and the fix is available in TF 0.14
      # https://github.com/hashicorp/terraform/pull/27141
      template[0].metadata[0].annotations["run.googleapis.com/sandbox"],
    ]
  }

  depends_on = [
    google_secret_manager_secret_version.ae-sa-key-secret,
    google_secret_manager_secret_iam_member.secret-access,
    google_project_service.project
  ]
}

resource "google_secret_manager_secret" "ae-sa-key" {
  secret_id = "allocation-endpoint-sa-key"

  replication {
    automatic = true
  }
  depends_on = [google_project_service.project]
}

resource "google_secret_manager_secret_version" "ae-sa-key-secret" {
  secret      = google_secret_manager_secret.ae-sa-key.id
  secret_data = base64decode(google_service_account_key.ae_sa_key.private_key)
}

resource "google_secret_manager_secret_iam_member" "secret-access" {
  secret_id  = google_secret_manager_secret.ae-sa-key.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
  depends_on = [google_project_service.project]
}

resource "google_project_service" "allocator-service" {
  for_each = var.game_gke_clusters

  service                    = google_endpoints_service.endpoints_service[each.key].id
  disable_dependent_services = true
}

resource "google_compute_address" "allocation-endpoint" {
  for_each = var.game_gke_clusters
  region   = each.value.region

  name = "allocator-endpoint-ip-${each.key}"
}

# Make Skaffold file for Cloud Deploy into each GKE Cluster
resource "local_file" "agones-skaffold-file" {
  content = templatefile(
    "${path.module}/files/agones/skaffold.yaml.tpl", {
      gke_clusters = var.game_gke_clusters
  })
  filename = "${path.module}/${var.platform_directory}/agones/install/skaffold.yaml"
}

# Make cluster specific helm value for LB IP
resource "local_file" "agones-ae-lb-file" {
  for_each = var.game_gke_clusters

  content = templatefile(
    "${path.module}/files/agones/ae-lb-ip-patch.yaml.tpl", {
      lb_ip = google_compute_address.allocation-endpoint[each.key].address
  })
  filename = "${path.module}/${var.platform_directory}/agones/install/${each.key}/kustomization.yaml"
}

# Create agones-system ns manifest as resource referenced by kustomization.yaml
resource "local_file" "agones-ns-file" {
  for_each = var.game_gke_clusters

  content  = file("${path.module}/files/agones/agones-system.yaml")
  filename = "${path.module}/${var.platform_directory}/agones/install/${each.key}/agones-system.yaml"
}

# Make Kubernetes manifest files to patch the Agones deployment for Allocation Endpoint
resource "local_file" "patch-agones-manifest" {
  for_each = var.game_gke_clusters

  content = templatefile(
    "${path.module}/files/agones/patch-agones-allocator.yaml.tpl", {
      project_id   = var.project
      location     = each.value.region
      cluster_name = each.key
      service_name = "${each.key}-${random_string.endpoint_random_string.result}.endpoints.${var.project}.cloud.goog"
      sa_email     = google_service_account.ae_sa.email
  })
  filename = "${path.module}/${var.platform_directory}/agones/endpoint-patch/patch-agones-allocator-${each.key}.yaml"
}
