package main

import (
	_ "analytic-service/docs"
	"analytic-service/internal/application"
	"analytic-service/internal/config"
	"log"
)

// @title        Analytics service
// @version      1.0
// @description  Service for collecting analytics about working with tasks from clients

// @host      localhost:8080
// @BasePath  /
func main() {
	pathToCfg, err := config.GetPathToCfgFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.New(pathToCfg)
	if err != nil {
		log.Fatal(err)
	}

	application.Run(cfg)
}
