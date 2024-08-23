package handlers

import (
	"accountant/internal/dependencies"
	"accountant/internal/services"
	"github.com/gofiber/fiber/v3"
)

type account struct {
	accountService *services.AccountService
}

func (h *account) list(c fiber.Ctx) error {
	userId := c.Locals("user").(int)

	accounts, err := h.accountService.List(c.Context(), userId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(accounts)
}

func SetupAccount(router fiber.Router, container *dependencies.Container) {
	h := account{container.AccountService()}

	router.Get("/", h.list)
}
