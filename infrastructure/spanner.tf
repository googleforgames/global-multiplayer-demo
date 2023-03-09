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

resource "google_spanner_instance" "global-game-spanner" {
  config       = var.spanner_config.location
  display_name = var.spanner_config.instance_name
  num_nodes    = var.spanner_config.num_nodes

  labels = {
    "environment" = var.resource_env_label
  }

  depends_on = [google_project_service.project]
}

resource "google_spanner_database" "spanner-database" {
  instance                 = google_spanner_instance.global-game-spanner.name
  name                     = var.spanner_config.db_name
  version_retention_period = "3d"
  deletion_protection      = false
}


resource "google_service_account_iam_binding" "spanner-workload-identity-binding" {
  service_account_id = google_service_account.spanner-sa.name
  role               = "roles/iam.workloadIdentityUser"

  members = [
    "serviceAccount:${var.project}.svc.id.goog[default/profile]"
  ]

  depends_on = [google_container_cluster.services-gke]
}

resource "google_service_account" "spanner-sa" {
  project      = var.project
  account_id   = "spanner-sa"
  display_name = "Spanner Service Account"
}

resource "google_project_iam_member" "spanner-sa" {
  project = var.project
  role    = "roles/spanner.databaseUser"
  member  = "serviceAccount:${google_service_account.spanner-sa.email}"
}

# Make liquibase.properties for schema management
resource "local_file" "liquibase-properties" {
  content = templatefile(
    "${path.module}/files/spanner/liquibase.properties.tpl", {
      project_id    = var.project
      instance_id   = google_spanner_instance.global-game-spanner.name
      database_id   = google_spanner_database.spanner-database.name
  })
  filename = "${path.module}/${var.schema_directory}/liquibase.properties"
}

# Make Config file for deploy with Cloud Deploy
resource "local_file" "services-profile-config" {
  content = templatefile(
    "${path.module}/files/services/profile-service-config.yaml.tpl", {
      service_email = google_service_account.spanner-sa.email
      project_id    = var.project
      instance_id   = google_spanner_instance.global-game-spanner.name
      database_id   = google_spanner_database.spanner-database.name
  })
  filename = "${path.module}/${var.services_directory}/profile/spanner_config.yaml"
}
