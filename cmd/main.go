package main

import (
	"cardcheck/internal/app/api"
	"cardcheck/internal/app/api/handler"
	"cardcheck/internal/app/service"
	"cardcheck/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"
)

const (
	timoutLimit = 5
)

// @title Cardcheck API
// @version	1.0
// @description	This is an API for validating credit cards.
// @contact.name Mark Raiter
// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info("Starting application...")
	log.Info("port: " + cfg.Server.AppAddress)

	validate := validator.New()
	service := service.New(log)
	handler := handler.New(service, validate, log)

	server := api.New(cfg, handler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.HTTPServer.Listen(cfg.Server.AppAddress); err != nil {
			log.Error("HTTPServer.Listen", "error", err)
		}
	}()

	<-stop

	if err := server.HTTPServer.ShutdownWithTimeout(timoutLimit * time.Second); err != nil {
		log.Error("ShutdownWithTimeout", "error", err)
	}

	if err := server.HTTPServer.Shutdown(); err != nil {
		log.Error("Shutdown", "error", err)
	}

	log.Info("server stopped")
}
