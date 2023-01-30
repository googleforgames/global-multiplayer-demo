# Copyright 2023 Google LLC All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: agones
    app.kubernetes.io/managed-by: google-cloud-deploy
    deploy.cloud.google.com/delivery-pipeline-id: global-game-agones-deploy-pipeline-${location}
    deploy.cloud.google.com/location: ${location}
    deploy.cloud.google.com/project-id: ${project_id}
    deploy.cloud.google.com/release-id: rel-1
    deploy.cloud.google.com/target-id: global-game-agones-deploy-target-${location}
    heritage: Helm
    multicluster.agones.dev/role: allocator
    release: my-agones
  name: agones-allocator
  namespace: agones-system
spec:
  progressDeadlineSeconds: 600
  replicas: 3
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: agones
      heritage: Helm
      multicluster.agones.dev/role: allocator
      release: my-agones
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        prometheus.io/path: /metrics
        prometheus.io/port: "8080"
        prometheus.io/scrape: "true"
        revision/tls-cert: "1"
      labels:
        app: agones
        app.kubernetes.io/managed-by: google-cloud-deploy
        deploy.cloud.google.com/delivery-pipeline-id: global-game-agones-deploy-pipeline-${location}
        deploy.cloud.google.com/location: ${location}
        deploy.cloud.google.com/project-id: ${project_id}
        deploy.cloud.google.com/release-id: rel-1
        deploy.cloud.google.com/target-id: global-game-agones-deploy-target-${location}
        heritage: Helm
        multicluster.agones.dev/role: allocator
        release: my-agones
    spec:
      affinity:
        nodeAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - preference:
              matchExpressions:
              - key: agones.dev/agones-system
                operator: Exists
            weight: 1
      containers:
      - args:
        - --listener_port=9443
        - --generate_self_signed_cert
        - --backend=grpc://127.0.0.1:8443
        - --service=agones-allocation-endpoint-${location}.endpoints.${project_id}.cloud.goog
        - --rollout_strategy=managed
        image: gcr.io/endpoints-release/endpoints-runtime:2
        imagePullPolicy: IfNotPresent
        name: esp
        ports:
        - containerPort: 9443
      - env:
        - name: GRPC_PORT
          value: "8443"
        - name: API_SERVER_QPS
          value: "400"
        - name: API_SERVER_QPS_BURST
          value: "500"
        - name: PROMETHEUS_EXPORTER
          value: "true"
        - name: STACKDRIVER_EXPORTER
          value: "false"
        - name: GCP_PROJECT_ID
        - name: STACKDRIVER_LABELS
        - name: DISABLE_MTLS
          value: "true"
        - name: DISABLE_TLS
          value: "true"
        - name: REMOTE_ALLOCATION_TIMEOUT
          value: 10s
        - name: TOTAL_REMOTE_ALLOCATION_TIMEOUT
          value: 30s
        - name: POD_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
        - name: CONTAINER_NAME
          value: agones-allocator
        - name: LOG_LEVEL
          value: info
        - name: FEATURE_GATES
        - name: ALLOCATION_BATCH_WAIT_TIME
          value: 500ms
        image: us-docker.pkg.dev/agones-images/release/agones-allocator:1.29.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /live
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 3
          periodSeconds: 3
          successThreshold: 1
          timeoutSeconds: 1
        name: agones-allocator
        ports:
        - containerPort: 8443
          name: grpc
          protocol: TCP
        - containerPort: 8080
          name: http
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /ready
            port: 8080
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: agones-allocator
      serviceAccountName: agones-allocator
      terminationGracePeriodSeconds: 30
      tolerations:
      - effect: NoExecute
        key: agones.dev/agones-system
        operator: Equal
        value: "true"
---
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
---
apiVersion: v1
kind: ServiceAccount
metadata:
  annotations:
    iam.gke.io/gcp-service-account: allocation-endpoint-esp-sa@${project_id}.iam.gserviceaccount.com
  name: agones-allocator
  namespace: agones-system
---
