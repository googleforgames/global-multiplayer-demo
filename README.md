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

### Deploy To GKE Clusters 

Replace the` _RELEASE_NAME` substitution with a unique build name. Cloudbuild
will deploy Agones using Cloud Deploy. 
```shell
$ cd deploy/
$ gcloud builds submit --config=cloudbuild.yaml --substitutions=_RELEASE_NAME=rel-1
```

## Licence

Apache 2.0

This is not an officially supported Google product
