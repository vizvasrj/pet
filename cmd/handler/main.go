package main

import (
	"fmt"
	"log"
	"net/http"
	"src/env"
	"src/etheus"
	"src/handler"
	"src/middleware"
	"src/petstore"
	"src/storageservice"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags | log.Ltime)
	envs := env.GetEnvs()
	db := storageservice.GetConnection(envs)
	defer db.Close()

	ps := handler.PetHandler{
		Storage: &storageservice.StorageService{
			Db: db,
		},
	}
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v3").Subrouter()
	petstore.HandlerFromMux(ps, s)

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			fmt.Println("ROUTE:", pathTemplate, "Methods:", methods)
		}
		return nil
	})
	r.Handle("/metrics", promhttp.Handler())

	logger := log.New(log.Writer(), "", 0)
	m := middleware.NewMiddleware()
	m.Logger = logger

	excludeRoutes := []string{"/metrics"}
	m.ExcludeRoutes = excludeRoutes

	r2 := m.LogMiddleware(r)

	prometheusRouter := m.PrometheusMiddleware(r2)

	log.Fatal(http.ListenAndServe(":8080", prometheusRouter))

}

func init() {
	prometheus.MustRegister(etheus.RequestCounter)
	prometheus.MustRegister(etheus.RequestDurationHistogram)
}
