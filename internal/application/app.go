package application

import (
	v1 "analytic-service/internal/adapters/http/v1"
	"analytic-service/internal/adapters/repository/idempotencyKeyRepo"
	"analytic-service/internal/adapters/repository/tasksRepo"
	"analytic-service/internal/adapters/rpc"
	"analytic-service/internal/config"
	"analytic-service/internal/domain/service"
	"analytic-service/pkg/auth"
	"analytic-service/pkg/httpServer"
	"analytic-service/pkg/kafka"
	"analytic-service/pkg/logging"
	"analytic-service/pkg/postgre"
	"context"
	"log"
	"net"
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

	// Database
	pg, err := postgre.New(
		logger,
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.DbName,
		cfg.Db.NAttemptsToConnect, cfg.Db.IsNeedMigration,
	)
	if err != nil {
		log.Fatalf("error when creating connection to db: %s", err.Error())
	}
	logger.Debugf("connection to the database was successful")
	repo := tasksRepo.New(pg)

	// Domain
	domainService := service.New(repo, logger)

	// ValidateToken
	authService, err := auth.New(cfg.Auth.Host, cfg.Auth.Port)
	if err != nil {
		log.Fatal(err)
	}

	// Rest
	handler := v1.CreateHandler(domainService)
	restServer := httpServer.New(
		handler.GetHttpHandler(authService.ValidateTokenStub, logger.MiddlewareLogging),
		cfg.Http.Port,
		cfg.Http.ReadTimeout,
		cfg.Http.WriteTimeout,
	)
	restServer.Run()
	logger.Debugf("http server started successfully")

	// idempotencyKeyValidator
	idempotencyRepo := idempotencyKeyRepo.New(pg)

	// Broker
	consumerConnection, err := kafka.NewConsumerConnection(cfg.Kafka.ClientID, net.JoinHostPort(cfg.Kafka.Host, cfg.Kafka.Port))
	if err != nil {
		logger.Fatalf("Error connecting to kafka: %s", err.Error())
	}
	broker := rpc.NewKafkaConsumer(logger, domainService, idempotencyRepo, *consumerConnection, cfg.Kafka.Topic)
	logger.Debug("connection to kafka successfully created")

	interrupt := make(chan os.Signal, 1)
	broker.Run(interrupt)

	logger.Info("the service is ready to work")

	// Shutdown
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Infof("an interrupt signal was received: %s", s.String())
	case err = <-restServer.Notify():
		logger.Fatalf("httpServer.Notify: %s", err.Error())
	case err = <-broker.Notify():
		logger.Fatalf("broker.Notify: %s", err.Error())
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
	if err := broker.Shutdown(); err != nil {
		logger.Errorf("error during closing connection with kafka: %s", err.Error())
	}

}
