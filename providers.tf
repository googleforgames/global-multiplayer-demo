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

provider "google" {
  project = var.project
}

data "google_client_config" "provider" {}

data "google_container_cluster" "spanner-gke-provider" {
  name     = var.spanner_gke_config.cluster_name
  location = var.spanner_gke_config.location
}

data "google_container_cluster" "gke-provider" {
  name     = var.spanner_gke_config.cluster_name
  location = var.spanner_gke_config.location
}

provider "kubernetes" {
  host  = "https://${data.google_container_cluster.gke-provider.endpoint}"
  token = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.gke-provider.master_auth[0].cluster_ca_certificate,
  )
}
