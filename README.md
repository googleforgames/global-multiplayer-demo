# Global, World Scale Multiplayer Game Demo

This multiplayer demo is a cloud first implementation of a global scale, realtime multiplayer game utilising
dedicated game servers, utlising both Google Cloud's products and open source gaming solutions.

### Prerequisites

To run the Game Demo install, you will need the following applications installed on your workstation:

* [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
* [Google Cloud CLI](https://cloud.google.com/sdk/docs/install)

## Google Cloud Auth

Once you have Google Cloud CLI installed, you will need to authenticate against Google Cloud:

```shell
$ gcloud auth application-default login
```

and then set your Google Cloud Project project name/PROJECT_ID:

```shell
$ gcloud config set project <PROJECT_ID>
```


### Provision

Initialize Terraform  & configure variables

```shell
$ terraform init
$ cp terraform.tfvars.sample terraform.tfvars

# Edit terraform.tfvars, especially <PROJECT_ID>
```

Provision the infrastructure.

```shell
$ terraform apply
```

### Deploy Agones To GKE Clusters 

The Agones deployment is in two steps: The Initial Install and the Allocation Endpoint Patch.

#### Initial Install
Replace the` _RELEASE_NAME` substitution with a unique build name. Cloudbuild will deploy Agones using Cloud Deploy. 

```shell
$ cd deploy/agones/install
$ gcloud builds submit --config=cloudbuild.yaml --substitutions=_RELEASE_NAME=rel-1
```

You can monitor the status of the deployment through the Cloud Logging URL returned by the `gcloud builds` command as well as the Kubernetes Engine/Worloads panel in the GCP Console. Once the Worloads have been marked as OK, you can proceed to apply the Allocation Endpoint Patch.

#### Allocation Endpoint Patch
After the Agones install has completed and the GKE Workloads show complete, run the Allocation Endpoint Patch Cloud Deploy to apply the appropriate endpoint patches to each cluster: 

```shell
$ cd deploy/agones/endpoint-patch/
$ gcloud builds submit --config=cloudbuild.yaml
```

***NOTE*** - The cloudbuild.yaml, kustomization.yaml & skaffold.yaml files will not exist until Terraform runs for the first time! The templates used for these files are stored in `files/agones/`.

You can monitor the status of the deployment through the Cloud Logging URL returned by the `gcloud builds` comma
nd as well as the Kubernetes Engine/Worloads panel in the GCP Console. Once the Worloads have been marked as O
K, Agones should be avaialable. 

### Deploy Spanner Applications to GKE Cluster

#### Initial Deploy
Replace the` _RELEASE_NAME` substitution with a unique build name. Cloudbuild will deploy Spanner applications using Cloud Deploy.

```shell
$ cd deploy/spanner/install
$ gcloud builds submit --config=cloudbuild.yaml --substitutions=_RELEASE_NAME=rel-1
```

## Licence

Apache 2.0

This is not an officially supported Google product
