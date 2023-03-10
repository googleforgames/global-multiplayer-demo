## Cross-cluster communication.

Anthos Service Mesh (ASM) is installed and configured across all agones clusters and the game services cluster. This allows cross cluster communication between services in the mesh. To add a service to the mesh, a sidecar container can be injected into any pod containing the `istio.io/rev: asm-managed` label. This label can be seen as a patch to the agones-allocator pod as part of the agones deployment manifest. The mesh is used to establish communication between the open-match director and agones allocator.

In addition to cross-cluster communication, the mesh is configured to route requests to the regional agones allocator endpoints based on request header from the open match director. To make a client request from the global game services GKE cluster to one of the regional endpoints, you would ensure:

1. That the sidecar is injected appropriately by applying the `istio.io/rev: asm-managed` label to the pod
1. Make a request with a `region` header having a value matching the region of the target cluster. The following is a request made from a pod on the global game services GKE cluster:

`curl agones-allocator.agones-system:8000/gameserverallocation -X POST -H "region: asia-east1"`

> Note: Requests without the `region` header are redirected to a faulty endpoint and will error out. This is to avoid accidental load balancing between the regional endpoints. 