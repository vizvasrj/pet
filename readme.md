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

## create postgres k8s secret, pvc, deployment and service
```bash
kubectl apply -f k8s/postgres/postgres-secret-pvc-deployment-service.yaml
```

---
## create petstore-storage-service deployment and service

```bash
kubectl apply -f k8s/petstore-storage/petstore-storage-deployment-service.yaml
```

---
### for migrations
```bash
kubectl apply -f k8s/migrations/petstore-storage-service-migrations.yaml
```


