package api

import (
	"context"
	"errors"
	"fmt"

	_ "cardcheck/docs"
	"cardcheck/internal/app/api/handler"
	"cardcheck/internal/config"
	"cardcheck/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/gofiber/swagger"
)

type Server struct {
	HTTPServer *fiber.App
}

func New(cfg *config.Config, handler *handler.Handler) *Server {
	server := new(Server)

	fconfig := fiber.Config{
		ReadTimeout:  cfg.Server.AppReadTimeout,
		WriteTimeout: cfg.Server.AppWriteTimeout,
		IdleTimeout:  cfg.Server.AppIdleTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "internal server error"

			var localError *fiber.Error
			if errors.As(err, &localError) {
				code = localError.Code
				if localError.Message != "" {
					message = localError.Message
				}
			}

			return c.Status(code).JSON(domain.ResponseMessage{Message: message})
		},
	}
	server.HTTPServer = fiber.New(fconfig)
	server.HTTPServer.Use(recover.New())
	server.HTTPServer.Use(logger.New())
	server.initRoutes(server.HTTPServer, handler, cfg)

	return server
}

func (s *Server) Shutdown(ctx context.Context) error {
	const op = "api.Server.Shutdown"

	return fmt.Errorf("%s: %w", op, s.HTTPServer.ShutdownWithContext(ctx))
}

func (s Server) initRoutes(app *fiber.App, handler *handler.Handler, cfg *config.Config) {
	// Swagger documentation
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "/swagger/doc.json",
		DeepLinking:  true,
		DocExpansion: "none",
		Title:        "Cardcheck API",
	}))

	// API routes
	app.Post("/check", timeout.NewWithContext(handler.CheckCard.Validate, cfg.Server.AppWriteTimeout))
}
