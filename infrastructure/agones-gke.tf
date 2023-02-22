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


data "google_container_engine_versions" "regions" {
  for_each = var.game_gke_clusters

  location = each.value.region
}

module "agones_gke_clusters" {
  for_each = var.game_gke_clusters

  source = "git::https://github.com/googleforgames/agones.git//install/terraform/modules/gke/?ref=main"

  cluster = {
    name             = each.key
    location         = each.value.region
    project          = var.project
    autoscale        = true
    workloadIdentity = true
    machineType      = each.value.machine_type

    # Install Current GKE default version
    kubernetesVersion = data.google_container_engine_versions.regions[each.key].default_cluster_version

    network    = google_compute_network.vpc.id
    subnetwork = "global-game-${each.value.region}-subnet"
  }
  udpFirewall = false

  depends_on = [google_compute_subnetwork.subnet, google_project_service.project]
}

data "google_container_cluster" "game-demo-agones-gke" {
  for_each = var.game_gke_clusters

  name     = each.key
  location = each.value.region

  depends_on = [module.agones_gke_clusters]
}

resource "google_compute_firewall" "agones-gameservers" {
  name    = "agones-gameservers"
  project = var.project
  network = google_compute_network.vpc.id

  allow {
    protocol = "udp"
    ports    = ["7000-8000"]
  }

  target_tags   = ["game-server"]
  source_ranges = ["0.0.0.0/0"]
}
