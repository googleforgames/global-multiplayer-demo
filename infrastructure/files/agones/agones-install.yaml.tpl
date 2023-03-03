helmCharts:
  - name: agones
    repo: https://agones.dev/chart/stable
    version: 1.30.0
    releaseName: agones
    namespace: agones-system
    valuesInline:
      agones:
        crds:
          cleanupOnDelete: false
        featureGates: "SplitControllerAndExtensions=true"
        allocator:
          disableMTLS: true
          disableTLS: true
          service:
            serviceType: ClusterIP
            http:
              port: 8000
              targetPort: 8000
              portName: http-alloc
            grpc:
              enabled: false

resources:
  - agones-system.yaml

patches:
  - target:
      kind: ServiceAccount
      name: agones-allocator
    patch: |-
      apiVersion: v1
      kind: ServiceAccount
      metadata:
        annotations:
          iam.gke.io/gcp-service-account: ${sa_email}
        name: agones-allocator
        namespace: agones-system
  - target:
      kind: Deployment
      name: agones-allocator
    patch: |-
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: agones-allocator
        namespace: agones-system
      spec:
        template:
          metadata:
            labels:
              istio.io/rev: asm-managed  #ASM managed dataplane channel
              region: ${location}        #Region to identify the POD and send traffic
