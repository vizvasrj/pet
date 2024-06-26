package main

import (
	"src/pkg/handler/app"
	"src/pkg/handler/etheus"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	app := &app.App{
		Metrics: etheus.NewMetrics(),
		Logger:  logger,
	}
	app.Initialize()
	app.SetupRoutes()
	app.Run(":8080")
}
