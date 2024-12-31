package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/vanthang24803/go-api/internal/schema"
	"github.com/vanthang24803/go-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Define a middleware functions
func AuthMiddleware(db *mongo.Database) fiber.Handler {
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

		var user schema.User
		collection := db.Collection("users")

		err = collection.FindOne(context.Background(), bson.M{"_id": payload.Sub}).Decode(&user)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusUnauthorized).JSON(utils.NewAppError(401, "User not found"))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(500, "Internal server error"))
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func AuthorizeRoles(requiredRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		user, ok := c.Locals("user").(schema.User)

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.NewAppError(500, "Unauthorized"))
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
