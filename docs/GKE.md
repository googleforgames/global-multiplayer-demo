# GKE

This repository provides support for running the backend applications on GKE Autopilot utilizing Workload Identity.

## Setup

### Terraform
The provided Terraform [gke.tf](../gke.tf) will provision a [GKE Autopilot](https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-overview) cluster. This method of provisioning GKE will automatically manage the nodes for the cluster as the backend applications are added.

Additionally, there is a service account created for the backend services.

### Kubectl
To interact with the GKE cluster, ensure kubectl is installed.

Once that is done, authenticate to GKE with the following commands:

```
export USE_GKE_GCLOUD_AUTH_PLUGIN=True
export GKE_CLUSTER=global-game-gke # change this based on the terraform configuration
gcloud container clusters get-credentials $GKE_CLUSTER --region us-central1
kubectl get namespaces
```

If there are no issues with the kubectl commands, kubectl is properly authenticated.
