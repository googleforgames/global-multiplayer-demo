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

resource "google_compute_network" "vpc" {
  name                    = var.vpc_name
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  for_each      = var.vpc_regions
  name          = "global-game-${each.key}-subnet"
  ip_cidr_range = each.value.vpc_subnet_cidr
  region        = each.key
  network       = google_compute_network.vpc.id
}

resource "google_compute_router" "vpc_router" {
  for_each = var.vpc_regions
  name     = "global-game-${each.key}-router"
  region   = each.key
  network  = google_compute_network.vpc.id

  bgp {
    asn = 64514
  }
}

resource "google_compute_router_nat" "vpc_nat" {
  for_each                           = var.vpc_regions
  name                               = "global-game-${each.key}-nat"
  router                             = google_compute_router.vpc_router[each.key].name
  region                             = each.key
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}

resource "google_compute_firewall" "asm-multicluster-pods" {
  name    = "asm-multicluster-pods"
  project = var.project
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
  }

  allow {
    protocol = "icmp"
  }

  source_ranges = ["10.0.0.0/8"]
}
