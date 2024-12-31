package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vanthang24803/go-api/domain"
	"github.com/vanthang24803/go-api/functions"
	"github.com/vanthang24803/go-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	validate    *validator.Validate
	authService functions.AuthService
}

// NewAuthHandler initializes a new AuthHandler with dependencies.
func NewAuthHandler(db *mongo.Database) *AuthHandler {
	return &AuthHandler{
		validate:    validator.New(),
		authService: *functions.NewAuthService(db),
	}
}

// RegisterUser handles user registration requests.
func (h *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	// Create a new User struct
	user := new(domain.RegisterRequest)

	// Parse the request body into the user struct
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid input"))
	}

	// Validate the user struct
	if err := h.validate.Struct(user); err != nil {
		// If validation fails, return error with the message
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, err.Error()))
	}

	// Insert the user into the database
	result, err := h.authService.RegisterHandler(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(500, err.Error()))
	}

	// Return the created user object
	return c.JSON(utils.NewResponse(201, result))
}

func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {

	request := new(domain.LoginRequest)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid input"))
	}

	// Validate the user struct
	if err := h.validate.Struct(request); err != nil {
		// If validation fails, return error with the message
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, err.Error()))
	}

	account, err := h.authService.LoginHandler(request)

	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(utils.NewResponse(200, account))
}

// AuthRoutes defines the authentication-related routes.
func AuthRoutes(router fiber.Router, db *mongo.Database) {
	handler := NewAuthHandler(db)

	router.Post("/auth/register", handler.RegisterUser)

	router.Post("/auth/login", handler.LoginHandler)
}
