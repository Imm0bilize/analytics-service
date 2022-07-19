package application

import (
	v1 "analytic-service/internal/adapters/http/v1"
	"analytic-service/internal/adapters/repository"
	"analytic-service/internal/adapters/rpc"
	"analytic-service/internal/config"
	"analytic-service/internal/domain/service"
	"analytic-service/pkg/auth"
	"analytic-service/pkg/httpServer"
	"analytic-service/pkg/logging"
	"analytic-service/pkg/postgre"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	logger, err := logging.New(cfg.Logger.Level, cfg.Logger.TsFormat)
	if err != nil {
		log.Fatalf("error when creating logger: %s", err.Error())
	}

	// ValidateToken
	authService, err := auth.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Database
	pg, err := postgre.New(logger, cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.NAttemptsToConnect)
	if err != nil {
		log.Fatalf("error when creating connection to auth service: %s", err.Error())
	}
	repo := repository.NewPgRepo(pg)

	domainService := service.New(repo, logger)

	// Rest
	handler := v1.CreateHandler(domainService)
	restServer := httpServer.New(cfg, handler.GetHttpHandler(authService.ValidateTokenStub, logger.MiddlewareLogging))
	restServer.Run()

	// grpc/broker
	grpcServer, err := rpc.New(cfg, domainService)
	if err != nil {
		log.Fatalf("error when creating grpc server: %s", err.Error())
	}
	grpcServer.Run()

	// Shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Print("an interrupt signal was received " + s.String())
	case err = <-restServer.Notify():
		log.Fatalf("httpServer.Notify: %s", err.Error())
	case err = <-grpcServer.Notify():
		log.Fatalf("grpcServer.Notify: %s", err.Error())
	}

	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		cfg.AppParams.ShutdownTimeout,
	)
	defer cancelFn()

	if err := restServer.Shutdown(ctx); err != nil {
		logger.Errorf("error during shutdown httpServer: %v", err)
	}
	if err := pg.Shutdown(); err != nil {
		logger.Errorf("error during shutdown DB: %v", err)
	}
	if err := authService.Shutdown(); err != nil {
		logger.Errorf("error while close connection with auth service")
	}
	if err := grpcServer.Shutdown(); err != nil {
		logger.Errorf("error while close connection with auth service")
	}

}
