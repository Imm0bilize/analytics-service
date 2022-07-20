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
	logger.Debugf("logger was successfully created")

	// ValidateToken
	authService, err := auth.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Database
	pg, err := postgre.New(logger, cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.NAttemptsToConnect, cfg.Db.IsNeedMigration)
	if err != nil {
		log.Fatalf("error when creating connection to auth service: %s", err.Error())
	}
	logger.Debugf("connection to the database was successful")
	repo := repository.NewPgRepo(pg)

	domainService := service.New(repo, logger)

	// Rest
	handler := v1.CreateHandler(domainService)
	restServer := httpServer.New(cfg, handler.GetHttpHandler(authService.ValidateTokenStub, logger.MiddlewareLogging))
	restServer.Run()
	logger.Debugf("http server started successfully")

	// grpc/broker
	grpcServer, err := rpc.New(cfg, logger, domainService)
	if err != nil {
		log.Fatalf("error when creating grpc server: %s", err.Error())
	}
	grpcServer.Run()
	logger.Debugf("grpc server started successfully")

	logger.Info("the service is ready to work")

	// Shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Infof("an interrupt signal was received: %s", s.String())
	case err = <-restServer.Notify():
		logger.Fatalf("httpServer.Notify: %s", err.Error())
	case err = <-grpcServer.Notify():
		logger.Fatalf("grpcServer.Notify: %s", err.Error())
	}

	ctx, cancelFn := context.WithTimeout(
		context.Background(),
		cfg.AppParams.ShutdownTimeout,
	)
	defer cancelFn()

	if err := restServer.Shutdown(ctx); err != nil {
		logger.Errorf("error during shutdown httpServer: %s", err.Error())
	}
	if err := pg.Shutdown(); err != nil {
		logger.Errorf("error during shutdown DB: %s", err.Error())
	}
	if err := authService.Shutdown(); err != nil {
		logger.Errorf("error while close connection with auth service: %s", err.Error())
	}
	if err := grpcServer.Shutdown(); err != nil {
		logger.Errorf("error during shutdown gRPC server: %s", err.Error())
	}

}
