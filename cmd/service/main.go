package main

import (
	"analytic-service/internal/application"
	"analytic-service/internal/config"
	"log"
)

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
