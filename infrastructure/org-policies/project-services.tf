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

resource "google_project_service" "project" {
  project = var.project
  for_each = toset(var.gcp_project_services)
  service  = each.value

  timeouts {
    create = "30m"
    update = "40m"
  }
  
  # Ensure service is truly active before continuing onward
  provisioner "local-exec" {
    command = <<EOF
        while [ ! $(gcloud services list --project=${var.project} | grep ${each.value} | wc -l ) ];
        do
            sleep 1s 
        done
        EOF
  }

  disable_dependent_services = false
  disable_on_destroy = false
}