helmCharts:
  - name: agones
    repo: https://agones.dev/chart/stable
    version: 1.36.0
    releaseName: agones
    namespace: agones-system
    valuesInline:
      agones:
        crds:
          cleanupOnDelete: false
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
          labels:
            istio.io/rev: asm-managed  #ASM managed dataplane channel
            region: ${location}        #Region to identify the POD and send traffic

resources:
  - agones-system.yaml
