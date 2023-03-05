# Optionally apply these Org Policies, as specified in terraform.tfvars file

module "gcp_org_policy_v2_requireShieldedVm" {
  source  = "terraform-google-modules/org-policy/google//modules/org_policy_v2"
  version = "~> 5.2.0"

  count          = var.apply_org_policies == true ? 1 : 0
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
}

module "gcp_org_policy_v2_disableServiceAccountKeyCreation" {
  source  = "terraform-google-modules/org-policy/google//modules/org_policy_v2"
  version = "~> 5.2.0"

  count          = var.apply_org_policies == true ? 1 : 0
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
}

module "gcp_org_policy_v2_vmCanIpForward" {
  source  = "terraform-google-modules/org-policy/google//modules/org_policy_v2"
  version = "~> 5.2.0"

  count          = var.apply_org_policies == true ? 1 : 0
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
}

module "gcp_org_policy_v2_vmExternalIpAccess" {
  source  = "terraform-google-modules/org-policy/google//modules/org_policy_v2"
  version = "~> 5.2.0"

  count          = var.apply_org_policies == true ? 1 : 0
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
}
