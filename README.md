# Global, World Scale Multiplayer Game Demo

This multiplayer demo is a cloud first implementation of a global scale, realtime multiplayer game utilising
dedicated game servers, utlising both Google Cloud's products and open source gaming solutions.

## Deployment

### Enable APIs

Before you set up the infrastructure, it is important to enable the appropriate APIs using the gcloud command line.

You must [install and configure gcloud](https://cloud.google.com/sdk/docs/install-sdk).

When that's complete, ensure your gcloud project is set correctly.

```shell
gcloud config set project <PROJECT_ID>
```

> **NOTE:** You can find your PROJECT_ID in [Cloud Console](https://cloud.google.com/resource-manager/docs/creating-managing-projects#identifying_projects).

Then, run the following `gcloud` command to enable the appropriate APIs in your project:

```shell
gcloud services enable compute.googleapis.com \
    cloudbuild.googleapis.com \
    container.googleapis.com \
    artifactregistry.googleapis.com
```

### Provision

Get terraform set up and variables configured:

```shell
terraform init
cp terraform.tfvars.sample terraform.tfvars

# Edit terraform.tfvars, especially project
```

Create the GKE instance. This separate step is a requirement due to providers not being able to have a a resource dependency in Terraform. See the discussion in [this issue](https://github.com/hashicorp/terraform/issues/2430) for example.

```shell
terraform apply -target=google_container_cluster.demo-game-gke
```

Provision the rest of the infrastructure.

```shell
terraform apply
```

## Licence

Apache 2.0

This is not an officially supported Google product
