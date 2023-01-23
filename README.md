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

Provision the infrastructure.

```shell
terraform apply
```

## Licence

Apache 2.0

This is not an officially supported Google product
