# Global, World Scale Multiplayer Game Demo

This multiplayer demo is a cloud first implementation of a global scale, realtime multiplayer game utilising
dedicated game servers, utlising both Google Cloud's products and open source gaming solutions.

## Infrastructure and Services

### Prerequisites

To run the Game Demo install, you will need the following applications installed on your workstation:

* [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
* [Google Cloud CLI](https://cloud.google.com/sdk/docs/install)

You can also click on the following icon to open this repository in a 'batteries-included' [Google Cloud Shell](https://cloud.google.com/shell) web development environment.

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://ssh.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fgoogleforgames%2Fglobal-multiplayer-demo.git&cloudshell_git_branch=main&cloudshell_open_in_editor=README.md&cloudshell_workspace=.)

### Google Cloud Auth

Once you have Google Cloud CLI installed, you will need to authenticate against Google Cloud:

```shell
gcloud auth application-default login
```

and then set your Google Cloud Project to name/PROJECT_ID:

```shell
gcloud config set project <PROJECT_ID>
```

Clone this directory locally and, we'll also set an environment variable to it's root directory, for easy navigation:

```shell
git clone https://github.com/googleforgames/global-multiplayer-demo.git
cd global-multiplayer-demo
export GAME_DEMO_HOME=$(pwd)
```

### Provision

# Optional: GCS Backend

Normally Terraform stores the current state in the `terraform.tfstate` file locally. However, if you would like to have Terraform store the state file in a GCS Bucket, you can edit the `backend.tf.sample` file, change the `bucket = ` line to your already created GCS bucket, and rename the file to `backend.tf`.

Note: The GCS bucket does not have to exist in the same Google project as the Global Game but the account runnint terraform must have write & read access to the bucket.

# Initialize Terraform & configure variables

```shell
cd $GAME_DEMO_HOME/infrastructure
terraform init
cp terraform.tfvars.sample terraform.tfvars

# Edit terraform.tfvars, especially <PROJECT_ID>
```


# Provision the infrastructure.

```shell
terraform apply
```

#### Deploy Agones To Agones GKE Clusters

The Agones deployment is in two steps: The Initial Install and the Allocation Endpoint Patch.


#### Initial Install
Replace the` _RELEASE_NAME` substitution with a unique build name. Cloudbuild will deploy Agones using Cloud Deploy.

```shell
cd $GAME_DEMO_HOME/platform/agones/install
gcloud builds submit --config=cloudbuild.yaml --substitutions=_RELEASE_NAME=rel-1
```

Navigate to the [agones-deploy-pipeline](https://console.cloud.google.com/deploy/delivery-pipelines/us-central1/agones-deploy-pipeline) delivery pipeline to review the rollout status. Cloudbuild will create a Cloud Deploy release which automatically deploys Agones to the first game server cluster. Agones can be deployed to subsequent clusters by clicking on the `promote` button within the Pipeline visualization or by running the following gcloud command:

```shell
# Replace RELEASE_NAME with the unique build name
$ gcloud deploy releases promote --release=RELEASE_NAME --delivery-pipeline=agones-deploy-pipeline --region=us-central1`
```

Continue the promotion until Agones has been deployed to all clusters. 

You can monitor the status of the deployment through the Cloud Logging URL returned by the `gcloud builds` command as well as the Kubernetes Engine/Worloads panel in the GCP Console. Once the Worloads have been marked as OK, you can proceed to apply the Allocation Endpoint Patch.

#### Deploy Open Match to Services GKE Cluster

Replace the` _RELEASE_NAME` substitution with a unique build name. Cloudbuild will deploy Open Match using Cloud Deploy.

```shell
cd $GAME_DEMO_HOME/platform/open-match/
gcloud builds submit --config=cloudbuild.yaml --substitutions=_RELEASE_NAME=rel-1
```

## Install Game Backend Services

TODO: fill in once we have services.

## Game Client

To build the Game Client for your host machine, you will need:

* [Unreal Engine 5.1.0](https://www.unrealengine.com/en-US/download) for your platform.

Open the [`game`](./game) folder in Unreal Engine. Once finished opening, you can run the game client directly within 
the editor (Hit the ▶️ button), or we can package the project via: Platforms > {your host platform} > Package Project,
and execute the resultant package.

### Troubleshooting

##### This project was made with a different version of the Unreal Engine.

If you hit this issue, it may be that you are building on a different host platform than the original. To solve, 
click: More Options > Convert in-place.

The project should open as normal now.

## Licence

Apache 2.0

This is not an officially supported Google product
