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

resource "google_secret_manager_secret" "secret_github_packages" {
  secret_id = "github-auth"

  labels = {
    "environment" = var.resource_env_label
  }

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "pat_1" {
  secret  = google_secret_manager_secret.secret_github_packages.id
  enabled = true

  secret_data = "${var.github_username}:${var.github_pat}"
}

resource "google_secret_manager_secret_iam_binding" "cloud_build_binding" {
  project   = var.project
  secret_id = google_secret_manager_secret.secret_github_packages.id
  role      = "roles/secretmanager.secretAccessor"
  members = [
    "serviceAccount:cloudbuild-cicd@${var.project}.iam.gserviceaccount.com",
  ]
}
