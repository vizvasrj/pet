codegen-generate:
	@echo "Generating code..."
	oapi-codegen -package petstore -generate types,client,gorilla-server,spec swagger/openapi.yaml > petstore/petstore.gen.go

build:
	@echo "Building..."
	go build -o petstore cmd/main.go


protoc-storage-service:
	@echo "Generating Go files"
	cd protoStorageService && protoc --go_out=. --go-grpc_out=. \
		--go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto
