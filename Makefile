codegen-generate:
	@echo "Generating code..."
	oapi-codegen -package petstore -generate types,client,gorilla-server,spec swagger/openapi.yaml > petstore/petstore.gen.go

build:
	@echo "Building..."
	go build -o petstore cmd/main.go