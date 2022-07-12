package application

import (
	restServer "analytic-service/internal/adapters/http"
	"analytic-service/internal/config"
	"analytic-service/pkg/logging"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	logger, err := logging.New(cfg.Logger.Level, cfg.Logger.TsFormat)
	if err != nil {
		log.Fatal("fatal")
	}

	httpServer := restServer.New(cfg, logger.MiddlewareLogging)
	httpServer.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Print("an interrupt signal was received " + s.String())
	case err = <-httpServer.Notify():
		log.Fatalf("httpServer.Notify: %s", err.Error())
	}
}
