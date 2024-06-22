package main

import (
	"fmt"
	"log"
	"net/http"
	"src/env"
	"src/handler"
	"src/middleware"
	"src/petstore"
	"src/storageservice"

	"github.com/gorilla/mux"
)

func main() {
	envs := env.GetEnvs()
	db := storageservice.GetConnection(envs)
	defer db.Close()

	ps := handler.PetHandler{
		Storage: &storageservice.StorageService{
			Db: db,
		},
	}
	r := mux.NewRouter()
	petstore.HandlerFromMux(ps, r)

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			fmt.Println("ROUTE:", pathTemplate, "Methods:", methods)
		}
		return nil
	})

	logger := log.New(log.Writer(), "", 0)
	logMiddleware := middleware.NewLogMiddleware(logger).Middleware

	wrappedRouter := logMiddleware(r)

	log.Fatal(http.ListenAndServe(":8080", wrappedRouter))

}
