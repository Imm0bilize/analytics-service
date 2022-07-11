package application

import (
	restServer "analytic-service/internal/adapters/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	httpServer := restServer.New("8080", time.Second*10, time.Second*10)
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
