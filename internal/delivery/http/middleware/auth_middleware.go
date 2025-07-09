package middleware

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/model"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JWTMiddleware(secret []byte) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: secret,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.ErrUnauthorized
		},
	})
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	auth, _ := ctx.Locals("auth").(*model.Auth)
	return auth
}

