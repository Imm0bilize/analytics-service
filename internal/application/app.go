package application

import (
	"analytic-service/internal/adapters/database/postrgre"
	v1 "analytic-service/internal/adapters/http/v1"
	"analytic-service/internal/config"
	"analytic-service/pkg/auth"
	"analytic-service/pkg/httpServer"
	"analytic-service/pkg/logging"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	logger, err := logging.New(cfg.Logger.Level, cfg.Logger.TsFormat)
	if err != nil {
		log.Fatal(err)
	}

	// RestAPI
	authService, err := auth.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := postrgre.New(cfg)
	if err != nil {
		log.Fatal(err)
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

	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		cfg.AppParams.ShutdownTimeout,
	)
	defer cancelFn()

	if err := restServer.Shutdown(ctx); err != nil {
		logger.ErrorF("error during shutdown httpServer: %v", err)
	}
	if err := db.Shutdown(); err != nil {
		logger.ErrorF("error during shutdown DB: %v", err)
	}
	if err := authService.Shutdown(); err != nil {
		logger.Error("error while close connection with auth service")
	}

}
