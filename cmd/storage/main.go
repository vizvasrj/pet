package main

import (
	"log"
	"net"
	"src/env"
	"src/proto_storage"
	"src/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	envs := env.GetEnvs()
	db := storage.GetConnection(envs)
	s := &storage.StorageService{
		Db: db,
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile("./cert/server.crt", "./cert/server.key")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)

	}
	server := grpc.NewServer(grpc.Creds(creds))
	proto_storage.RegisterStorageServiceServer(server, s)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
