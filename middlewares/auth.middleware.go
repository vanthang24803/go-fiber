package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/vanthang24803/go-api/utils"
)

// Define a middleware functions
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.NewAppError(401, "Missing authorization header"))
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.NewAppError(401, "Wrong authorization header format"))
		}

		token := bearerToken[1]

		payload, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.NewAppError(401, "Invalid token"))
		}

		// Save payload to context
		c.Locals("user", payload)

		return c.Next()
	}
}

func AuthorizeRoles(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		user, ok := c.Locals("user").(*utils.JwtPayload)

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.NewAppError(401, "Unauthorized"))
		}

		hasRequiredRole := lo.SomeBy(user.Roles, func(role string) bool {
			return lo.Contains(requiredRoles, role)
		})

		if !hasRequiredRole {
			return c.Status(fiber.StatusForbidden).JSON(utils.NewAppError(403, "Access denied"))
		}

		return c.Next()
	}
}
