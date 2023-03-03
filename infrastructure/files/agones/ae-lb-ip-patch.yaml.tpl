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
            port: 443
            targetPort: 9443
          grpc:
            enabled: false
            loadBalancerIP: "${lb_ip}"

resources:
  - agones-system.yaml

patches:
  - target:
      kind: Service
      name: agones-allocator
    patch: |-
      apiVersion: v1
      kind: Service
      metadata:
        name: agones-allocator
        namespace: agones-system
      spec:
        ports:
        - name: https
          port: 443
          targetPort: 9443
        selector:
          $patch: replace
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
          spec:
            containers:
            - args:
              - --listener_port=9443
              - --generate_self_signed_cert
              - --backend=grpc://127.0.0.1:8443
              - --service=${service_name}
              - --rollout_strategy=managed
              image: gcr.io/endpoints-release/endpoints-runtime:2
              imagePullPolicy: IfNotPresent
              name: esp
              ports:
              - containerPort: 9443
