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


resource "google_redis_instance" "open-match" {
  name           = "global-game-open-match"
  tier           = "STANDARD_HA"
  memory_size_gb = 1
  region         = var.services_gke_config.location

  location_id             = "${var.services_gke_config.location}-a"
  alternative_location_id = "${var.services_gke_config.location}-f"

  authorized_network      = google_compute_network.vpc.id
  transit_encryption_mode = "DISABLED"
  connect_mode            = "PRIVATE_SERVICE_ACCESS"

  redis_version = "REDIS_6_X"
  display_name  = "Global Game Demo: Open Match"

  labels = {
    "environment" = var.resource_env_label
  }

  depends_on = [google_service_networking_connection.private_service_connection, google_project_service.project, google_container_cluster.services-gke]
}

resource "google_compute_global_address" "private_service_range" {
  name          = "private-service-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc.id

  depends_on = [google_project_service.project]
}

resource "google_service_networking_connection" "private_service_connection" {
  network                 = google_compute_network.vpc.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_service_range.name]

  depends_on = [google_project_service.project]
}

# Add Redis Host & IP to Open Match kustomization.yaml
resource "local_file" "open-match-kustomization-file" {

  content = templatefile(
    "${path.module}/files/open-match/kustomization.yaml.tpl", {
      redis_host = google_redis_instance.open-match.host
      redis_port = google_redis_instance.open-match.port
  })
  filename = "${path.module}/${var.platform_directory}/open-match/base/kustomization.yaml"
}
