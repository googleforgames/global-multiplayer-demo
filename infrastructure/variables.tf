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

variable "game_gke_clusters" {
  type        = map(any)
  description = "GKE gameclusters & associated values"
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

### Allocation Endpoint Variables ###

variable "allocation_endpoint" {
  type = object({
    name             = string
    proxy_image      = string
    weight           = number
    namespace        = string
    agones_namespace = string
  })
  description = "Allocation Endpoint Configuration Variables"
}

variable "platform_directory" {
  type        = string
  description = "Platform Directory for output to Cloud Deploy related files"
}

variable "services_directory" {
  type        = string
  description = "Services Directory for output to Cloud Deploy related files"
}
