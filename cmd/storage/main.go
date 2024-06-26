package main

import (
	"log"
	"net"
	"src/env"
	"src/pkg/storage/database"
	"src/pkg/storage/middleware"
	"src/pkg/storage/service"
	"src/proto_storage"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	envs := env.GetEnvs()
	db := database.GetConnection(envs)
	defer db.Close()

	s := &service.StorageService{
		Db:     db,
		Logger: logger,
	}
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile("./cert/server.crt", "./cert/server.key")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)

	}
	serverCreds := grpc.Creds(creds)
	recoveryMiddleware := grpc.UnaryInterceptor(middleware.RecoveryMiddleware(logger))
	server := grpc.NewServer(serverCreds, recoveryMiddleware)
	proto_storage.RegisterStorageServiceServer(server, s)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
