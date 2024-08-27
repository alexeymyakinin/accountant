package middlewares

import (
	"accountant/internal/config"
	"accountant/internal/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"slices"
	"strings"
)

type JWTMiddlewareConfig struct {
	Skip []string
}

func NewJWTMiddleware(config JWTMiddlewareConfig) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		if slices.Contains(config.Skip, string(ctx.Request().URI().Path())) {
			return ctx.Next()
		}

		headers := ctx.GetReqHeaders()
		val, ok := headers["Authorization"]
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing Authorization header")
		}

		if len(val) != 1 {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header")
		}

		valSplit := strings.Split(val[0], " ")
		if len(valSplit) != 2 {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header")
		}

		if valSplit[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid auth schema %s", valSplit[0])
		}

		token, err := getToken(valSplit[1])
		if err != nil {
			slog.Error("asd", slog.String("err", err.Error()))
			return err
		}

		ctx.Locals("user", token.Claims.(*domain.JWTClaims))

		return ctx.Next()
	}
}

func GetClaimsFromContext(ctx fiber.Ctx) *domain.JWTClaims {
	return ctx.Locals("user").(*domain.JWTClaims)
}

func getToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &domain.JWTClaims{}, parseToken)
}

func parseToken(_ *jwt.Token) (interface{}, error) {
	return []byte(config.Get().JwtSigningKey), nil
}
