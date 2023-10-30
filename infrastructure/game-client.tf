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


resource "google_service_account" "game_client_vm" {
  count = var.enable_game_client_vm ? 1 : 0

  project = var.project

  account_id   = "game-client-vm"
  display_name = "Custom SA for Game Client VM"
}

resource "google_compute_address" "game_client_vm_static_ip" {
  count = var.enable_game_client_vm ? 1 : 0

  project = var.project
  name    = "game-client-vm-static-ip"
  region  = var.game_client_vm_region
}

data "google_compute_image" "game_client_vm_os" {
  count = var.enable_game_client_vm ? 1 : 0

  family  = var.game_client_vm_os_family
  project = var.game_client_vm_os_project
}

resource "google_compute_instance" "game_client_vm" {
  count = var.enable_game_client_vm ? 1 : 0

  project = var.project

  name         = "game-client-vm"
  machine_type = var.game_client_vm_machine_type
  zone         = "${var.game_client_vm_region}-a"

  tags = ["game-client-vm-ssh"]

  scheduling {
    on_host_maintenance = "TERMINATE"
  }

  boot_disk {
    initialize_params {
      image = data.google_compute_image.game_client_vm_os[0].self_link
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "NVME"
  }

  network_interface {
    subnetwork = google_compute_subnetwork.subnet["${var.game_client_vm_region}"].self_link
    # network = google_compute_network.vpc.id

    access_config {
      // Ephemeral public IP
      nat_ip = google_compute_address.game_client_vm_static_ip[0].address
    }
  }

  metadata = {
    serial-port-logging-enable = "TRUE"
  }

  metadata_startup_script = file("${path.root}/game-client-startup.sh")

  service_account {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    email  = google_service_account.game_client_vm[0].email
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_firewall" "game-client-vm-ssh" {
  project = var.project

  name    = "game-client-vm-ssh"
  network = google_compute_network.vpc.id

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  target_tags   = ["game-client-vm-ssh"]
  source_ranges = var.game_client_vm_allowed_cidr
}
