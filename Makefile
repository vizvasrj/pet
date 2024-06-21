codegen-generate:
	oapi-codegen -package petstore -generate types,client,gorilla-server,spec swagger/openapi.yaml > petstore/petstore.gen.go

build:
	go build -o petstore cmd/main.go