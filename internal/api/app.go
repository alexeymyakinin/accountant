package api

import (
	"accountant/internal/dependencies"
	"accountant/internal/handlers"
	"accountant/internal/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
	"log/slog"
)

func New(container *dependencies.Container) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:    ErrorHandler,
		StructValidator: NewStructValidator(validator.New()),
	})

	setMiddlewares(app)
	setHandlers(app, container)
	return app
}

func setMiddlewares(app *fiber.App) {
	app.Use(recover.New())
	app.Use(slogfiber.New(slog.Default()))
	app.Use(middlewares.NewJWTMiddleware(middlewares.JWTMiddlewareConfig{Skip: []string{"/login", "/register"}}))
}

func setHandlers(app *fiber.App, container *dependencies.Container) {
	handlers.SetupAuth(app.Group(""), container)
	handlers.SetupAccount(app.Group("/account"), container)

}
