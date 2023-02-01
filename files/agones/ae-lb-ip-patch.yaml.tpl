helmCharts:
- name: agones
  repo: https://agones.dev/chart/stable
  version: 1.29.0
  releaseName: my-agones
  namespace: agones-system
  valuesInline:
    agones:
       allocator:
         disableMTLS: true
         disableTLS: true
         service:
           http:
             enabled: false
           loadBalancerIP: "${lb_ip}"
