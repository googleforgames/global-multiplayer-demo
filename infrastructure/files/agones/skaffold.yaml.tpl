apiVersion: skaffold/v4beta1
kind: Config
manifests:
  kustomize:
    paths:
      - ./path_to_cluster_specific_skaffold
    buildArgs:
      - --enable-helm
deploy:
  kubectl:
    flags:
      apply:
        - --server-side # Avoid the "Too long: must have at most 262144 bytes" problem
  tolerateFailuresUntilDeadline: true # Fixes startup timeouts
profiles:
%{ for cluster_name, values in gke_clusters ~}
  - name: ${cluster_name}
    patches:
      - op: replace
        path: /manifests/kustomize/paths/0
        value: "./${cluster_name}"
%{ endfor ~}
