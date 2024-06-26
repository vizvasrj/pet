package main

import (
	"src/pkg/handler/app"
	"src/pkg/handler/etheus"
)

func main() {
	app := &app.App{
		Metrics: etheus.NewMetrics(),
	}
	app.Initialize()
	app.SetupRoutes()
	app.Run(":8080")
}
