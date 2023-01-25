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
}

resource "google_spanner_database" "spanner-database" {
  instance                 = google_spanner_instance.global-game-spanner.name
  name                     = var.spanner_config.db_name
  version_retention_period = "3d"
  deletion_protection      = false

  depends_on = [google_project_service.project]
}
