package interceptor

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vanthang24803/go-api/utils"
)

func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			if appError, ok := err.(*utils.AppError); ok {
				return c.Status(appError.Code).JSON(appError)
			}

			return c.Status(500).JSON(fiber.Map{
				"message":   "Internal server error",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		}

		return nil
	}
}
