#!/usr/bin/bash

openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 365 -key ca.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out ca.crt
openssl req -newkey rsa:2048 -nodes -keyout server.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=*.example.com" -out server.csr
# openssl x509 -req -extfile <(printf "subjectAltName=DNS:example.com,DNS:www.example.com") -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
openssl x509 -req -extfile <(printf "subjectAltName=DNS:petstore-storage-service,DNS:petstore-storage-service.default.svc.cluster.local,DNS:petstore-handler-service,DNS:petstore-handler-service.default.svc.cluster.local,DNS:example.com") -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt

#use server.crt, server.key in grpc server
#use server.crt in grpc client
#do not forget to add example.com to your /etc/hosts
