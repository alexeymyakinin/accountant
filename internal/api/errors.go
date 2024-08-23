package api

import (
	"errors"
	"github.com/gofiber/fiber/v3"
)

type DetailedError struct {
	Detail string `json:"detail"`
}

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return c.Status(code).JSON(DetailedError{Detail: err.Error()})
}
