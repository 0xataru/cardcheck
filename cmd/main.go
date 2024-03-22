package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/markraiter/cardcheck/internal/app/api"
	"github.com/markraiter/cardcheck/internal/app/api/handler"
	"github.com/markraiter/cardcheck/internal/app/service"
	"github.com/markraiter/cardcheck/internal/config"
)

const (
	timoutLimit = 5
)

// @title Cardcheck API
// @version	1.0
// @description	This is an API for validating credit cards.
// @contact.name Mark Raiter
// @contact.email raitermark@proton.me
// host localhost:5555
// @BasePath /api
func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info("Starting application...")
	log.Info("port: " + cfg.Server.AppAddress)

	service := service.New(log)

	handler := handler.New(service, log)

	server := api.New(cfg, handler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.HTTPServer.Listen(cfg.Server.AppAddress); err != nil {
			log.Error("HTTPServer.Listen", err)
		}
	}()

	<-stop

	if err := server.HTTPServer.ShutdownWithTimeout(timoutLimit * time.Second); err != nil {
		log.Error("ShutdownWithTimeout", err)
	}

	if err := server.HTTPServer.Shutdown(); err != nil {
		log.Error("Shutdown", err)
	}

	log.Info("server stopped")
}
