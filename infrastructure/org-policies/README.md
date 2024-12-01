# Troubleshooting Org Policy Issues

Have you encountered any error messages related to `compute.requireShieldedVm`, `iam.disableServiceAccountKeyCreation`, `compute.vmCanIpForward`, or `compute.vmExternalIpAccess` when following the instructions in this repo's [quickstart guide](https://github.com/googleforgames/global-multiplayer-demo/blob/main/README.md)?  If so, please try the steps outlined in this README to update your GCP Org Policies in a manner that meets this demo's minimum requirements.

**BELOW:** Example error arising from a conflict with default Org Policy settings.
```shell
│ Error: Error creating service account key: googleapi: Error 400: Key creation is not allowed on this service account.
│ Details:
│ [
│   {
│     "@type": "type.googleapis.com/google.rpc.PreconditionFailure",
│     "violations": [
│       {
│         "description": "Key creation is not allowed on this service account.",
│         "subject": "projects/gameops-001/serviceAccounts/allocation-endpoint-esp-sa@gameops-001.iam.gserviceaccount.com?configvalue=allocation-endpoint-esp-sa%40gameops-001.iam.gserviceaccount.com",
│         "type": "constraints/iam.disableServiceAccountKeyCreation"
│       }
│     ]
│   }
│ ]
│ , failedPrecondition
│ 
│   with google_service_account_key.ae_sa_key,
│   on allocation-endpoint.tf line 54, in resource "google_service_account_key" "ae_sa_key":
│   54: resource "google_service_account_key" "ae_sa_key" {
│   omitted...
```
### Double-check the prerequisites 

Please read this repo's [primary README](https://github.com/googleforgames/global-multiplayer-demo/blob/main/README.md) and ensure that  all its main prerequisites have been met, such as creating a GCP Project

### Set up your environment

Navigate to this repo's `infrastructure/org-policies` subdirectory and then please run the following commands.

```shell
# Set your GCP Project ID
export PROJECT_ID=<PROJECT_ID>
gcloud config set project ${PROJECT_ID}

# Add this value to your default tfvars file
sed -i 's/PROJECT_ID/${PROJECT_ID}/g' terraform.tfvars.sample
cp terraform.tfvars.sample terraform.tfvar

# Set your GCP Organization ID
export ORGANIZATION_ID=<ORGANIZATION_ID>

# Enter the email address that you use to authenticate with GCP
export GCP_EMAIL_ADRESS=<GCP_EMAIL_ADRESS>
```

### Authenticate your user identity and generate Application Default Credentials (ADC)

Run the following command to simultanously authenticate you as an end user with Google Cloud and to obtain a set of Application Default Credentials.  Later on in this README, Terraform will perform a `user_project_override` and leverage your user identity to provision GCP org policies.  This is neccessary because org policies are a protected resource that (usually) cannot be edited with standard ADC credentials. 

```shell
gcloud auth login --no-launch-browser --brief --update-adc --quiet
```

### Update any missing org-level GCP IAM permissions

If you are already a super user in your GCP organization, run the following commands to ensure that you have the IAM permissions neccessary to edit GCP Org Policy resources.

```shell
gcloud organizations add-iam-policy-binding ${ORGANIZATION_ID} --condition=None --member="user:${GCP_EMAIL_ADRESS}" --role="roles/orgpolicy.PolicyAdmin"
```

If you do not have such elevated permissions, please share this README with your platform administrator and request that they either (1) grant you this IAM permission, or (2) run the provisioning steps in this guide on your behalf.

### Apply project-scoped overrides to your default GCP org policies

Now that you are authenticated and have the neccessary IAM permissions in place, run the following commands to create a few project-scoped overrides to your default org policies.

```shell
terraform init
terraform apply
```
More specifically, this step will edit the four [organization policy constraints](https://cloud.google.com/resource-manager/docs/organization-policy/org-policy-constraints) listed below.

- `compute.requireShieldedVm`
- `iam.disableServiceAccountKeyCreation`
- `compute.vmCanIpForward`
- `compute.vmExternalIpAccess`

### Final Step

Once you have completed this README, return to this repo's [quickstart guide](https://github.com/googleforgames/global-multiplayer-demo/blob/main/README.md) and try running through its instructions once again.  With any luck, your org policy issues will be resolved and you will be able to successfully finish provisioning the demo.

**BELOW:** Example output shown when Terraform successfully updates the organization policy constraints.
```shell
...omitted
google_project_service.project["cloudresourcemanager.googleapis.com"]: Creation complete after 7s [id=planet-scale-demo-044/cloudresourcemanager.googleapis.com]
module.gcp_org_policy_v2_vmCanIpForward.google_org_policy_policy.project_policy[0]: Creating...
module.gcp_org_policy_v2_disableServiceAccountKeyCreation.google_org_policy_policy.project_policy_boolean[0]: Creating...
module.gcp_org_policy_v2_vmExternalIpAccess.google_org_policy_policy.project_policy[0]: Creating...
module.gcp_org_policy_v2_requireShieldedVm.google_org_policy_policy.project_policy_boolean[0]: Creating...
module.gcp_org_policy_v2_vmCanIpForward.google_org_policy_policy.project_policy[0]: Creation complete after 3s [id=projects/planet-scale-demo-044/policies/compute.vmCanIpForward]
module.gcp_org_policy_v2_vmExternalIpAccess.google_org_policy_policy.project_policy[0]: Creation complete after 3s [id=projects/planet-scale-demo-044/policies/compute.vmExternalIpAccess]
module.gcp_org_policy_v2_requireShieldedVm.google_org_policy_policy.project_policy_boolean[0]: Creation complete after 4s [id=projects/planet-scale-demo-044/policies/compute.requireShieldedVm]
module.gcp_org_policy_v2_disableServiceAccountKeyCreation.google_org_policy_policy.project_policy_boolean[0]: Creation complete after 4s [id=projects/planet-scale-demo-044/policies/iam.disableServiceAccountKeyCreation]
Apply complete! Resources: 6 added, 0 changed, 0 destroyed.
```

## Licence

Apache 2.0

This is not an officially supported Google product
