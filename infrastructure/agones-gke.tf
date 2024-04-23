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

##------------------------------##
## Google Terraform: Agones GKE ##
##------------------------------##

data "google_container_engine_versions" "regions" {
  for_each = merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters)

  location = each.value.region
}

module "agones_gke_standard_clusters" {
  for_each = var.game_gke_standard_clusters

  source = "git::https://github.com/googleforgames/agones.git//install/terraform/modules/gke/?ref=v1.40.0"

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

module "agones_gke_autopilot_clusters" {
  for_each = var.game_gke_autopilot_clusters

  source = "git::https://github.com/googleforgames/agones.git//install/terraform/modules/gke-autopilot/?ref=v1.40.0"

  cluster = {
    name     = each.key
    location = each.value.region
    project  = var.project

    # Install Current GKE default version
    kubernetesVersion = data.google_container_engine_versions.regions[each.key].default_cluster_version

    network    = google_compute_network.vpc.id
    subnetwork = "global-game-${each.value.region}-subnet"
  }
  udpFirewall = false

  depends_on = [google_compute_subnetwork.subnet, google_project_service.project]
}

data "google_container_cluster" "game-demo-agones" {
  for_each = merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters)

  name     = each.key
  location = each.value.region

  depends_on = [module.agones_gke_standard_clusters, module.agones_gke_autopilot_clusters]
}

resource "google_gke_hub_membership" "membership" {
  for_each      = merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters)
  provider      = google-beta
  project       = var.project
  membership_id = "${each.key}-membership"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${data.google_container_cluster.game-demo-agones[each.key].id}"
    }
  }

  depends_on = [google_project_service.project]
}

resource "google_gke_hub_feature" "mesh" {
  name     = "servicemesh"
  project  = var.project
  location = "global"
  provider = google-beta

  depends_on = [google_project_service.project]
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

# Make Skaffold file for Cloud Deploy into each GKE Cluster
resource "local_file" "agones-skaffold-file" {
  content = templatefile(
    "${path.module}/files/agones/skaffold.yaml.tpl", {
      gke_clusters = merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters)
  })
  filename = "${path.module}/${var.platform_directory}/agones/skaffold.yaml"
}

# Make cluster specific helm value for LB IP
resource "local_file" "agones-ae-lb-file" {
  for_each = merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters)

  content = templatefile(
    "${path.module}/files/agones/agones-install.yaml.tpl", {
      location = each.value.region
  })
  filename = "${path.module}/${var.platform_directory}/agones/${each.key}/kustomization.yaml"
}

# Create agones-system ns manifest as resource referenced by kustomization.yaml
resource "local_file" "agones-ns-file" {
  for_each = merge(var.game_gke_standard_clusters, var.game_gke_autopilot_clusters)

  content  = file("${path.module}/files/agones/agones-system.yaml")
  filename = "${path.module}/${var.platform_directory}/agones/${each.key}/agones-system.yaml"
}

