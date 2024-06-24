package main

import (
	"src/app"
	"src/etheus"
)

func main() {
	// Initialize Prometheus metrics with App
	app := &app.App{
		Metrics: etheus.NewMetrics(),
	}
	app.Initialize()
	app.SetupRoutes()
	app.Run(":8080")
}
