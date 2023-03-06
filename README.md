# Global, World Scale Multiplayer Game Demo

This multiplayer demo is a cloud first implementation of a global scale, realtime multiplayer game utilising
dedicated game servers, utlising both Google Cloud's products and open source gaming solutions.

## OAuth Authentication

We need to manually set up the OAuth authentication, as unfortunately this cannot be automated.

The details, such as name and email address of both of these steps don't matter, so feel free to use something
arbitrary for any part not specified.

Open the [Google OAuth consent screen](https://console.cloud.google.com/apis/credentials/consent) for your project,
and create an "External" App, and allowlist any users you wish to be able to login to your deployment of this game.

Open the [Google Credentials](https://console.cloud.google.com/apis/credentials) screen for your project, and click 
"+ CREATE CREDENTIALS", and create an "OAuth Client ID" of type "Web Application".

Leave this page open, as we'll need the Client ID and Client secret of the ID you just created shortly. 

## Infrastructure and Services

### Prerequisites

To run the Game Demo install, you will need the following applications installed on your workstation:

* A Google Cloud Project
* [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
* [Google Cloud CLI](https://cloud.google.com/sdk/docs/install)

You can also click on the following icon to open this repository in a 'batteries-included' [Google Cloud Shell](https://cloud.google.com/shell) web development environment.

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://ssh.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fgoogleforgames%2Fglobal-multiplayer-demo.git&cloudshell_git_branch=main&cloudshell_open_in_editor=README.md&cloudshell_workspace=.)

### Google Cloud Auth

Once you have Google Cloud CLI installed, you will need to set your [GCP Project ID](https://support.google.com/googleapi/answer/7014113?hl=en#:~:text=The%20project%20ID%20is%20a,ID%20or%20create%20your%20own.):

```shell
export PROJECT_ID=<PROJECT_ID>
gcloud config set project ${PROJECT_ID}
```

and then authenticate to generate [Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/application-default-credentials) that can be leveraged by Terraform
```shell
gcloud auth application-default login
```

Clone this directory locally and, we'll also set an environment variable to it's root directory, for easy navigation:

```shell
git clone https://github.com/googleforgames/global-multiplayer-demo.git
cd global-multiplayer-demo
export GAME_DEMO_HOME=$(pwd)
```

## Provision

### Optional: GCS Backend

Normally Terraform stores the current state in the `terraform.tfstate` file locally. However, if you would like to have Terraform store the state file in a GCS Bucket, you can:

- [ ] Edit `backend.tf.sample`
- [ ] Change the `bucket =` line to an already created GCS bucket
- [ ] Rename `backend.tf.sample` to `backend.tf`.

NOTE: The GCS bucket does not have to exist in the same Google project as the Global Game but the Google user/service account running Terraform must have write & read access to that bucket.

### Initialize Terraform & configure variables

```shell
cd $GAME_DEMO_HOME/infrastructure
terraform init
cp terraform.tfvars.sample terraform.tfvars
```

You will need to now edit `terraform.tfvars`

* Update <PROJECT_ID> with the ID of your Google Cloud Project, 
* Updated <CLIENT_ID> and <CLIENT_SECRET> with the Client ID and Client secret created in the above step.

You can edit other variables in this file, but we recommend leaving the default values for your first run before 
experimenting.

### Provision the infrastructure.

> **Warning**  
> This demo in its default state creates multiple Kubernetes clusters around the world, 
> Spanner instances, and more. Running this demo for an extended amount of time may incur significant costs.


```shell
terraform apply
```

### OAuth Authentication

We now need to update our OAuth authentication configuration with the address of our authenticating frontend API.

Let's grab the IP for that API, by running:

```shell
gcloud compute addresses list --filter=name=frontend-service --format="value(address)"
```

This should give you back an IP, such as `35.202.107.204`.

1. Click "+ ADD URI" under "Authorised JavaScript origins" and add "http://[IP_ADDRESS].sslip.io".
2. Click "+ ADD URI" under "Authorised redirect URIs" and add "http://[IP_ADDRESS].sslip.io/callback"
3. Click "Save".

Since OAuth needs a domain to authenticate against, we'll use [sslip.io](https://sslip.io) for development purposes. 

### Deploy Agones To Agones GKE Clusters

The Agones deployment is in two steps: The Initial Install and the Allocation Endpoint Patch.

### Initial Install
Replace the` _RELEASE_NAME` substitution with a unique build name. Cloudbuild will deploy Agones using Cloud Deploy.

```shell
cd $GAME_DEMO_HOME/platform/agones/
gcloud builds submit --config=cloudbuild.yaml --substitutions=_RELEASE_NAME=rel-1
```

Navigate to the [agones-deploy-pipeline](https://console.cloud.google.com/deploy/delivery-pipelines/us-central1/agones-deploy-pipeline) delivery pipeline to review the rollout status. Cloudbuild will create a Cloud Deploy release which automatically deploys Agones to the first game server cluster. Agones can be deployed to subsequent clusters by clicking on the `promote` button within the Pipeline visualization or by running the following gcloud command:

```shell
# Replace RELEASE_NAME with the unique build name
gcloud deploy releases promote --release=RELEASE_NAME --delivery-pipeline=agones-deploy-pipeline --region=us-central1`
```

Continue the promotion until Agones has been deployed to all clusters. 

You can monitor the status of the deployment through the Cloud Logging URL returned by the `gcloud builds` command as well as the Kubernetes Engine/Worloads panel in the GCP Console. Once the Worloads have been marked as OK, you can proceed to apply the Allocation Endpoint Patch.

### Deploy Open Match to Services GKE Cluster

Replace the` _RELEASE_NAME` substitution with a unique build name. Cloudbuild will deploy Open Match using Cloud Deploy.

```shell
cd $GAME_DEMO_HOME/platform/open-match/
gcloud builds submit --config=cloudbuild.yaml --substitutions=_RELEASE_NAME=rel-1
```

## Install Game Backend Services

To install all the backend services, submit the following Cloud Build command, and replace the` _RELEASE_NAME` 
substitution with a unique build name.

```shell
cd $GAME_DEMO_HOME/services
gcloud builds submit --config=cloudbuild.yaml --substitutions=_RELEASE_NAME=rel-1
```

This will:

* Build all the images required for all services.
* Store those image in [Artifact Registry](https://cloud.google.com/artifact-registry)
* Deploy them via Cloud Build to a Autopilot cluster.

## Game Client

To build the Game Client for your host machine, you will need:

* [Unreal Engine 5.1.0](https://www.unrealengine.com/en-US/download) for your platform.

Open the [`game`](./game) folder in Unreal Engine. Once finished opening, you can run the game client directly within 
the editor (Hit the ▶️ button), or we can package the project via: Platforms > {your host platform} > Package Project,
and execute the resultant package.

## Run the Game Launcher

{TODO: still requires more details}

```shell
cd $GAME_DEMO_HOME/game/GameLauncher

# Grab the IP Address again of our frontend service, so we can use it
gcloud compute addresses list --filter=name=frontend-service --format="value(address)"
```

Edit the app.ini, and replace the `frontend_api` value with http://[IP_ADDRESS].sslip.io

The run:

```shell
go run main.go
```

### Troubleshooting

##### This project was made with a different version of the Unreal Engine.

If you hit this issue, it may be that you are building on a different host platform than the original. To solve, 
click: More Options > Convert in-place.

The project should open as normal now.

## Licence

Apache 2.0

This is not an officially supported Google product
