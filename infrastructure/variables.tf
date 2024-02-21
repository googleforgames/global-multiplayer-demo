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

### Project Variables ###

variable "project" {
  type        = string
  description = "GCP Project Name"
}

variable "gcp_project_services" {
  type        = list(any)
  description = "GCP Service APIs (<api>.googleapis.com) to enable for this project"
  default     = []
}

variable "resource_env_label" {
  type        = string
  description = "Label/Tag to apply to resources"
}

### VPC Variables ###

variable "vpc_name" {
  type        = string
  description = "VPC Name"
}

variable "vpc_regions" {
  type        = map(any)
  description = "Regions for VPC Subnets to be created"
}

### Agones Variables ###

variable "game_gke_standard_clusters" {
  type        = map(any)
  description = "GKE Standard Game Clusters & Associated values"
}

variable "game_gke_autopilot_clusters" {
  type        = map(any)
  description = "GKE Autopilot Game Clusters & Associated values"
}

### Cloud Deploy Variables ###

variable "clouddeploy_config" {
  type = object({
    location = string
  })
}

### Artifact Registry Variables ###

variable "artifact_registry_config" {
  type = object({
    location = string
  })
}

### Spanner Variables ###

variable "spanner_config" {
  type = object({
    instance_name = string
    db_name       = string
    location      = string
    num_nodes     = number
  })

  description = "Configuration specs for Spanner"
}

variable "schema_directory" {
  type        = string
  description = "Schema directory where schema definition is found"
}

### Services GKE Variables ###

variable "services_gke_config" {
  type = object({
    cluster_name    = string
    location        = string
    resource_labels = map(string)
  })

  description = "Configuration specs for GKE Autopilot cluster that hosts all backend services"
}

variable "app_service_account_config" {
  type = object({
    name        = string
    description = string
  })
  description = "The configuration specifications for the backend service account"
}

variable "k8s_service_account_id" {
  type        = string
  description = "The kubernetes service account that will impersonate the IAM service account to access Cloud Spanner. This account will be created."
}

### Frontend Service Variables ###

variable "frontend-service" {
  type = object({
    client_id     = string
    client_secret = string
    jwt_key       = string
  })
  description = "Configuration for the frontend service that provides oAuth authentications"
}

variable "platform_directory" {
  type        = string
  description = "Platform Directory for output to Cloud Deploy related files"
}

variable "services_directory" {
  type        = string
  description = "Services Directory for output to Cloud Deploy related files"
}

### Open Match Match Function Variables ###

variable "open-match-matchfunction" {
  type = object({
    players_per_match = number
  })
  description = "Configuration for the Open Match Match Function"
}

### Dedicated Game Server Variables

variable "github_username" {
  type        = string
  description = "The GitHub username that matches to the `github_pat` personal access token"
}

variable "github_pat" {
  type        = string
  description = "A GitHub personal access token (classic) with at least repo scope"
}

### Game Client VM Variables

variable "enable_game_client_vm" {
  type        = bool
  description = "Whether to create or not a Linux Game Client VM"
  default     = false
}

variable "game_client_vm_machine_type" {
  type        = string
  description = "Game Client VM Machine Type"
}

variable "game_client_vm_allowed_cidr" {
  type        = list(any)
  description = "Game Client VM Allowed CIDRs"
}

variable "game_client_vm_region" {
  type        = string
  description = "Game Client VM Region"
}

variable "game_client_vm_storage" {
  type        = number
  description = "Game Client VM Storage Size"
}

variable "game_client_vm_os_family" {
  type        = string
  description = "Game Client VM OS Image Family"
}

variable "game_client_vm_os_project" {
  type        = string
  description = "Game Client OS Image Project"
}
