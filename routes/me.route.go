package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vanthang24803/go-api/functions"
	"github.com/vanthang24803/go-api/internal/schema"
	"github.com/vanthang24803/go-api/middlewares"
	"github.com/vanthang24803/go-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeHandler struct {
	validate *validator.Validate
	meSevice functions.MeService
}

func NewMeHandler(db *mongo.Database) *MeHandler {
	return &MeHandler{
		validate: validator.New(),
		meSevice: *functions.NewMeService(db),
	}
}

func (h *MeHandler) MeHandler(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(schema.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(491, "Unauthorized"))
	}

	return c.JSON(utils.NewResponse(200, user))
}

func MeRoutes(router fiber.Router, db *mongo.Database) {
	handler := NewMeHandler(db)
	router.Get("/me", middlewares.AuthMiddleware(db), handler.MeHandler)
}
