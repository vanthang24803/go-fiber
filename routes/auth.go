package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vanthang24803/go-api/internal/schema"
	"github.com/vanthang24803/go-api/queries"
	"github.com/vanthang24803/go-api/utils"
)

// Declare validate and db globally
var validate *validator.Validate

func init() {
	// Initialize validator
	validate = validator.New()
}

func AuthRoutes(router fiber.Router) {
	router.Post("/register", func(c *fiber.Ctx) error {
		// Create a new User struct
		user := new(schema.User)

		// Parse the request body into the user struct
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid input"))
		}

		// Validate the user struct
		if err := validate.Struct(user); err != nil {
			// If validation fails, return error with the message
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, err.Error()))
		}

		// Insert the user into the database
		user, err := queries.Register(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(500, err.Error()))
		}

		// Return the created user object
		return c.JSON(utils.NewResponse(201, user))
	})
}
