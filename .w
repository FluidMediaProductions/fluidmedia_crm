apiVersion: v1
kind: Endpoints
metadata:
  name: glusterfs-cluster
  namespace: mailu
subsets:
  - addresses:
    - ip: 138.201.93.184
    ports:
    - port: 1
  - addresses:
    - ip: 88.99.122.54
    ports:
    - port: 1
  - addresses:
    - ip: 195.201.58.86
    ports:
    - port: 1
---
apiVersion: v1
kind: Service
metadata:
  name: glusterfs-cluster
  namespace: mailu
spec:
  ports:
  - port: 1
