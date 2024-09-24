package middleware

import (
	"IFEST/helpers"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return helpers.HttpUnauthorized(c, "unauthorized, invalid token")
		}

		token := auth[7:]
		claims := helpers.UserClaims{}

		if _, err := helpers.DecodeJWT(token, &claims); err != nil {
			return helpers.HttpUnauthorized(c, "unauthorized, invalid token")
		}

		c.Locals("userID", claims.ID)
		return c.Next()
	}
}
