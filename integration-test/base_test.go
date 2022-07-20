package integration_test

import (
	"analytic-service/internal/application"
	"analytic-service/internal/config"
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

type integrationTestSuite struct {
	suite.Suite
	cfg           *config.Config
	dbContainer   testcontainers.Container
	authContainer testcontainers.Container
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &integrationTestSuite{})
}

func (i *integrationTestSuite) SetupSuite() {
	cfg, err := config.New("../configs/test-config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	i.cfg = cfg

	ctx := context.Background()

	// DB
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.4-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "secret",
			"DB_USER":           "postgres",
			"DB_PASSWORD":       "secret",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	ip, err := dbContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err := dbContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal(err)
	}
	i.cfg.Db.Host = ip
	i.cfg.Db.Port = mappedPort.Port()
	i.dbContainer = dbContainer
	if err := dbContainer.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// auth
	req = testcontainers.ContainerRequest{
		Image:        "auth-service:latest",
		ExposedPorts: []string{"3000/tcp", "9000/tcp"},
		WaitingFor:   wait.ForLog("Grpc server started on port 3000"),
	}
	authContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	ip, err = authContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err = authContainer.MappedPort(ctx, "3000")
	if err != nil {
		log.Fatal(err)
	}
	i.cfg.Auth.Host = ip
	i.cfg.Auth.Port = mappedPort.Port()

	i.authContainer = authContainer
	if err := authContainer.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// Analytics service
	go application.Run(cfg)
	delayForRunningService := time.Second * 2
	time.Sleep(delayForRunningService)
}

func (i *integrationTestSuite) TearDownSuite() {
	ctx := context.Background()

	if err := i.dbContainer.Stop(ctx, &i.cfg.AppParams.ShutdownTimeout); err != nil {
		i.T().Errorf("error when stopping db container: %s", err.Error())
	}
	if err := i.authContainer.Stop(ctx, &i.cfg.AppParams.ShutdownTimeout); err != nil {
		i.T().Errorf("error when stopping auth container: %s", err.Error())
	}
}
