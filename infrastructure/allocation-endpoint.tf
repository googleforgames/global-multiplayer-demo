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

resource "google_service_account_iam_binding" "workload-identity-binding" {
  service_account_id = google_service_account.ae_sa.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:${var.project}.svc.id.goog[agones-system/agones-allocator]",
  ]

  depends_on = [module.agones_gke_standard_clusters, module.agones_gke_autopilot_clusters]
}

resource "google_service_account" "ae_sa" {
  account_id   = "allocation-endpoint-esp-sa"
  display_name = "Service Account for Allocation Endpoint"
}

resource "google_service_account_key" "ae_sa_key" {
  service_account_id = google_service_account.ae_sa.name
}

resource "google_secret_manager_secret" "ae-sa-key" {
  secret_id = "allocation-endpoint-sa-key"

  replication {
    automatic = true
  }

  labels = {
    "environment" = var.resource_env_label
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