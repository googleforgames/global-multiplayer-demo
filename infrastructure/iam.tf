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

resource "google_service_account" "cloudbuild-sa" {
  project      = var.project
  account_id   = "cloudbuild-cicd"
  display_name = "Cloud Build - CI/CD service account"
}

resource "google_project_iam_member" "cloudbuild-sa-cloudbuild" {
  project = var.project
  role    = "roles/cloudbuild.builds.builder"
  member  = "serviceAccount:${google_service_account.cloudbuild-sa.email}"
}

resource "google_project_iam_member" "cloudbuild-sa-gke-admin" {
  project = var.project
  role    = "roles/container.admin"
  member  = "serviceAccount:${google_service_account.cloudbuild-sa.email}"
}

resource "google_project_iam_member" "cloudbuild-sa-cloudbuild-roles" {
  project = var.project
  for_each = toset([
    "roles/serviceusage.serviceUsageAdmin",
    "roles/clouddeploy.operator",
    "roles/cloudbuild.builds.builder",
    "roles/container.admin",
    "roles/storage.admin",
    "roles/iam.serviceAccountUser",
    "roles/spanner.databaseUser",
    "roles/gkehub.editor",
    "roles/compute.viewer"
  ])
  role   = each.key
  member = "serviceAccount:${google_service_account.cloudbuild-sa.email}"
}

resource "google_project_iam_member" "clouddeploy-iam" {
  project = var.project
  for_each = toset([
    "roles/container.admin",
    "roles/artifactregistry.reader",
    "roles/storage.admin"
  ])
  role   = each.key
  member = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"

  depends_on = [google_project_service.project]
}
