package interceptor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vanthang24803/go-api/utils"
)

func NotFoundRoute() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.SendStatus(404)

		return c.JSON(utils.NewAppError(404, "Not found route!"))
	}
}
