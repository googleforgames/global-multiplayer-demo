# Global, World Scale Multiplayer Game Demo

This multiplayer demo is a cloud first implementation of a global scale, realtime multiplayer game utilising
dedicated game servers, utlising both Google Cloud's products and open source gaming solutions.

### Provision

Get terraform set up and variables configured:

```shell
$ terraform init
$ cp terraform.tfvars.sample terraform.tfvars

# Edit terraform.tfvars, especially project
```

Provision the infrastructure.

```shell
$ terraform apply
```

### Deploy To GKE Clusters 

The below will list all GKE clusters in your project:
 
```shell
$ gcloud container clusters list
```

The below will list all Cloud Deploy Pipelines for a region:
 
```shell
$ gcloud deploy delivery-pipelines list --region=us-central1|grep name|awk -F\/ '{print $6}'
```

The below will deploy release `release-v1` from a provided K8s manifest file to a GKE cluster through it's delCloud Deploy Delivery Pipeline:
 
```shell
$ gcloud deploy releases create release-v1 --from-k8s-manifest=release-v1.yaml --region=us-central1 --delivery-pipeline=global-game-agones-deploy-pipeline-us-central1 
```

## Licence

Apache 2.0

This is not an officially supported Google product
