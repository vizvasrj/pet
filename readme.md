## create docker images
```bash
docker build -t us.icr.io/vizvasrj/petstore-storage-service:1 -f dockerfiles/Dockerfile.petstore-storage-service .
docker build -t us.icr.io/vizvasrj/petstore-storage-service-migrations:1 -f dockerfiles/Dockerfile.petstore-storage-service-migrations .
docker build -t us.icr.io/vizvasrj/petstore-handler-service:1 -f dockerfiles/Dockerfile.petstore-handler-service .
```

## copy images to kind
```bash
kind load docker-image us.icr.io/vizvasrj/petstore-storage-service:1
kind load docker-image us.icr.io/vizvasrj/petstore-storage-service-migrations:1
kind load docker-image us.icr.io/vizvasrj/petstore-handler-service:1
```

---

## create postgres k8s deployment and service
### create secret
```bash
kubectl apply -f k8s/postgres/postgres-secret.yaml
```
### create pvc, deployment and service
```bash
kubectl apply -f k8s/postgres/postgres-pvc.yaml
```
### create deployment and service
```bash
kubectl apply -f k8s/postgres/postgres-deployment.yaml
```

### create service
```bash
kubectl apply -f k8s/postgres/postgres-service.yaml
```
---
## create petstore-storage-service deployment and service

### create petstore-storage-service deployment
```bash
kubectl apply -f k8s/petstore-storage/petstore-storage-deployment.yaml
```

### create petstore-storage-service service
```bash
kubectl apply -f k8s/petstore-storage/petstore-storage-service.yaml
```
---
### for migrations
```bash
kubectl apply -f k8s/migrations/petstore-storage-service-migrations.yaml
```


