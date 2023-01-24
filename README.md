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

Replace the` _RELEASE_NAME` substitution with a unique build name. Cloudbuild
will deploy Agones using Cloud Deploy. 
```shell
gcloud builds submit --config=cloudbuild.yaml --substitutions=_RELEASE_NAME=rel-0001
```

## Licence

Apache 2.0

This is not an officially supported Google product
