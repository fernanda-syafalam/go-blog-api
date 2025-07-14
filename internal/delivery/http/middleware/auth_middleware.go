package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/common/response"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/model"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf"
)

func JWTMiddleware(k *koanf.Koanf) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendErrorResponse(c, response.TokenNotFound)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.SendErrorResponse(c, response.InvalidToken)
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString, k.String("jwt.secret"))
		if err != nil {
			fmt.Println(err)
			return utils.SendErrorResponse(c, response.Unauthorized)
		}

		userId, err := strconv.ParseUint(claims.UserID, 10, 64)
		if err != nil {
			return utils.SendErrorResponse(c, response.ServerError)
		}
		c.Locals("userID", uint(userId))
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	auth, _ := ctx.Locals("auth").(*model.Auth)
	return auth
}
