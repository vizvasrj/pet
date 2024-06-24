package app

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
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	Router     *mux.Router
	Storage    *storageservice.StorageService
	Log        *log.Logger
	httpServer *http.Server
	Metrics    *etheus.Metrics
}

func (a *App) Initialize() {
	envs := env.GetEnvs()
	db := storageservice.GetConnection(envs)

	a.Storage = &storageservice.StorageService{
		Db: db,
	}
	a.Router = mux.NewRouter()
	a.Log = log.New(log.Writer(), "", 0)
}

// Setup routes and middleware
func (a *App) SetupRoutes() {
	m := middleware.NewMiddleware()
	m.Logger = a.Log
	m.RequestCounter = a.Metrics.RequestCounter
	m.RequestDurationHistogram = a.Metrics.RequestDurationHistogram
	// Route for metrics
	a.Router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	a.Router.Use(m.LogMiddleware) // Log all requests

	// Create a subrouter for the API
	api := a.Router.PathPrefix("/api/v3").Subrouter()

	// Initialize your handler
	// petHandler := handler.PetHandler{
	// 	Storage: a.Storage,
	// }
	petHandler, err := handler.NewPetHandler("example.com:8081")
	if err != nil {
		a.Log.Fatalf("Error initializing PetHandler: %v", err)
	}
	petstore.HandlerFromMux(petHandler, api)

	// Middleware (order matters!)
	api.Use(
		m.PrometheusMiddleware, // Prometheus
	)
	printRoutes(a.Router)
}

// Run the application
func (a *App) Run(addr string) {
	a.httpServer = &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}

	a.Log.Printf("Starting server on %s", addr)
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.Log.Fatalf("Server error: %v", err)
	}
}

func printRoutes(r *mux.Router) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			if methods != nil {
				fmt.Println("ROUTE:", pathTemplate, "Methods:", methods)
			}
		}
		return nil
	})
}
