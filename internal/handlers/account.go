package handlers

import (
	"accountant/internal/dependencies"
	"accountant/internal/middlewares"
	"accountant/internal/services"
	"github.com/gofiber/fiber/v3"
)

type account struct {
	accountService *services.AccountService
}

func (h *account) list(ctx fiber.Ctx) error {
	user := middlewares.GetClaimsFromContext(ctx)
	userId := user.UserId

	accounts, err := h.accountService.List(ctx.Context(), userId)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(accounts)
}

func SetupAccount(router fiber.Router, container *dependencies.Container) {
	h := account{container.AccountService()}

	router.Get("/", h.list)
}
