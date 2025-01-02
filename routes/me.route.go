package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vanthang24803/go-api/domain"
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

func (h *MeHandler) UpdateMeHandler(c *fiber.Ctx) error {
	payload, ok := c.Locals("user").(*utils.JwtPayload)

	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	req := new(domain.UpdateProfileRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid input"))
	}

	user, err := h.meService.UpdateProfileHandler(payload, req)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.JSON(utils.NewResponse(200, user))
}

func (h *MeHandler) UpdateAvatarHandler(c *fiber.Ctx) error {
	payload, ok := c.Locals("user").(*utils.JwtPayload)

	if !ok || payload == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewAppError(401, "Unauthorized"))
	}

	file, err := c.FormFile("avatar")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewAppError(400, "Invalid input"))
	}

	return c.JSON(utils.NewResponse(200, file))
}

func MeRoutes(router fiber.Router, db *mongo.Database) {
	handler := NewMeHandler(db)
	router.Get("/me", middlewares.AuthMiddleware(), handler.MeHandler)
	router.Put("/me", middlewares.AuthMiddleware(), handler.UpdateMeHandler)
	router.Put("/me/avatar", middlewares.AuthMiddleware(), handler.UpdateAvatarHandler)
}
