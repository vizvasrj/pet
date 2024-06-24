package main

import (
	"src/app"
	"src/etheus"
)

func main() {
	// Initialize Prometheus metrics with App
	// storageServiceURL := "http://localhost:8081"

	app := &app.App{
		Metrics: etheus.NewMetrics(),
	}
	app.Initialize()
	app.SetupRoutes()

	app.Run(":8080")
}
