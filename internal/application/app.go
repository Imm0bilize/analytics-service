package application

import (
	v1 "analytic-service/internal/adapters/http/v1"
	"analytic-service/internal/config"
	"analytic-service/pkg/auth"
	"analytic-service/pkg/httpServer"
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

	// RestAPI
	authService, err := auth.New(cfg)
	if err != nil {
		log.Fatal()
	}

	handler := v1.CreateHandler(authService.ValidateTokenStub, logger.MiddlewareLogging)
	restServer := httpServer.New(cfg, handler)
	restServer.Run()

	// Shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Print("an interrupt signal was received " + s.String())
	case err = <-restServer.Notify():
		log.Fatalf("httpServer.Notify: %s", err.Error())
	}
}
