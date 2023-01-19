# Global, World Scale Multiplayer Game Demo

This multiplayer demo is a cloud first implementation of a global scale, realtime multiplayer game utilising
dedicated game servers, utlising both Google Cloud's products and open source gaming solutions.

### Provision

Get terraform set up and variables configured:

```shell
terraform init
cp terraform.tfvars.sample terraform.tfvars

# Edit terraform.tfvars, especially project
```

Create the GKE instance. This separate step is a requirement due to providers not being able to have a a resource dependency in Terraform. See the discussion in [this issue](https://github.com/hashicorp/terraform/issues/2430) for example.

```shell
terraform apply -target=google_container_cluster.game-demo-spanner-gke
```

Provision the rest of the infrastructure.

```shell
terraform apply
```

## Licence

Apache 2.0

This is not an officially supported Google product
