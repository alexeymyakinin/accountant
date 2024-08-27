package api

import (
	"errors"
	"github.com/gofiber/fiber/v3"
)

type DetailedError struct {
	Detail string `json:"detail"`
}

func ErrorHandler(ctx fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	body := DetailedError{Detail: "Internal server error"}

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		body.Detail = e.Message
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return ctx.Status(code).JSON(body)
}
