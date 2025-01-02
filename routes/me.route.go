package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vanthang24803/go-api/functions"
	"github.com/vanthang24803/go-api/middlewares"
	"github.com/vanthang24803/go-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeHandler struct {
	validate  *validator.Validate
	meService functions.MeService
}

func NewMeHandler(db *mongo.Database) *MeHandler {
	return &MeHandler{
		validate:  validator.New(),
		meService: *functions.NewMeService(db),
	}
}

func (h *MeHandler) MeHandler(c *fiber.Ctx) error {
	payload, ok := c.Locals("user").(*utils.JwtPayload)
	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	user, err := h.meService.GetProfileHandler(payload)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(utils.NewResponse(200, user))
}

func MeRoutes(router fiber.Router, db *mongo.Database) {
	handler := NewMeHandler(db)
	router.Get("/me", middlewares.AuthMiddleware(), handler.MeHandler)
}
