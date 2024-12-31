package interceptor

import (
	"time"

	"github.com/gofiber/fiber/v2"

	logger "github.com/vanthang24803/go-api/infra"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		err := c.Next()

		duration := time.Since(startTime)

		statusCode := c.Response().StatusCode()

		logger.Msg.Infof("%s %s - %d - %v", c.Method(), c.Path(), statusCode, duration)

		if err != nil {
			logger.Msg.With("error", err.Error()).Error("Request failed!")
		}

		return err
	}
}
