package application

import (
	restServer "analytic-service/internal/adapters/http"
	"analytic-service/pkg/logging"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {

	logger, err := logging.New("DEBUG", time.RFC3339)
	if err != nil {
		log.Fatal("fatal")
	}

	httpServer := restServer.New("8080", time.Second*10, time.Second*10, logger.MiddlewareLogging)
	httpServer.Run()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Print("an interrupt signal was received " + s.String())
	case err := <-httpServer.Notify():
		log.Fatalf("httpServer.Notify: %s", err.Error())
	}
}
