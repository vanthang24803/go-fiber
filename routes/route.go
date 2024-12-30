package routes

import "github.com/gofiber/fiber/v2"

func ApplyRoutes(router fiber.Router) {
	AuthRoutes(router)
}
