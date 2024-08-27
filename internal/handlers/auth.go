package handlers

import (
	"accountant/internal/dependencies"
	"accountant/internal/services"
	"errors"
	"github.com/gofiber/fiber/v3"
)

type authHandler struct {
	authService *services.AuthService
}

func (h *authHandler) login(c fiber.Ctx) error {
	req := loginRequest{}
	err := c.Bind().Body(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	t, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.UserNotFoundError):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return err
		}
	}

	return c.JSON(fiber.Map{"token": t})
}

func (h *authHandler) register(c fiber.Ctx) error {
	req := registerRequest{}
	err := c.Bind().Body(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	t, err := h.authService.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.EmailAlreadyExistsError):
			return fiber.NewError(fiber.StatusBadRequest, services.EmailAlreadyExistsError.Error())
		default:
			return err
		}
	}

	return c.JSON(fiber.Map{"token": t})
}

func SetupAuth(router fiber.Router, container *dependencies.Container) {
	h := authHandler{
		authService: container.AuthService(),
	}

	router.Post("/login", h.login)
	router.Post("/register", h.register)
}
