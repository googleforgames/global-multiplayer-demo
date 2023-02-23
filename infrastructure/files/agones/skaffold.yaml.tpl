apiVersion: skaffold/v2beta29
kind: Config
deploy:
  kustomize:
    paths:
      - "./path_to_cluster_specific_skaffold"
    buildArgs: ["--enable-helm"]
    flags:
      apply: ['--server-side'] # Avoid the "Too long: must have at most 262144 bytes" problem
profiles:
%{ for cluster_name, values in gke_clusters ~}
- name: ${cluster_name}
  patches:
  - op: replace
    path: /deploy/kustomize/paths/0
    value: "./${cluster_name}"
%{ endfor ~}