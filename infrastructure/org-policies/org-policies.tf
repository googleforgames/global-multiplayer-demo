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


module "gcp_org_policy_v2_requireShieldedVm" {
  source  = "terraform-google-modules/org-policy/google//modules/org_policy_v2"
  version = "~> 5.2.0"

  policy_root    = "project"
  policy_root_id = var.project
  rules = [{
    enforcement = false
    allow       = []
    deny        = []
    conditions  = []
  }]
  constraint  = "compute.requireShieldedVm"
  policy_type = "boolean"
  
  depends_on = [google_project_service.project]
}

module "gcp_org_policy_v2_disableServiceAccountKeyCreation" {
  source  = "terraform-google-modules/org-policy/google//modules/org_policy_v2"
  version = "~> 5.2.0"

  policy_root    = "project"
  policy_root_id = var.project
  rules = [{
    enforcement = false
    allow       = []
    deny        = []
    conditions  = []
  }]
  constraint  = "iam.disableServiceAccountKeyCreation"
  policy_type = "boolean"

  depends_on = [google_project_service.project]
}

module "gcp_org_policy_v2_vmCanIpForward" {
  source  = "terraform-google-modules/org-policy/google//modules/org_policy_v2"
  version = "~> 5.2.0"

  policy_root    = "project"
  policy_root_id = var.project
  rules = [{
    enforcement = false
    allow       = []
    deny        = []
    conditions  = []
  }]
  constraint  = "compute.vmCanIpForward"
  policy_type = "list"

  depends_on = [google_project_service.project]
}

module "gcp_org_policy_v2_vmExternalIpAccess" {
  source  = "terraform-google-modules/org-policy/google//modules/org_policy_v2"
  version = "~> 5.2.0"

  policy_root    = "project"
  policy_root_id = var.project
  rules = [{
    enforcement = false
    allow       = []
    deny        = []
    conditions  = []
  }]
  constraint  = "compute.vmExternalIpAccess"
  policy_type = "list"

  depends_on = [google_project_service.project]
}
