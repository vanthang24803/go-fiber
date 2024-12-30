package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/vanthang24803/go-api/infra"
	"github.com/vanthang24803/go-api/interceptor"
	"github.com/vanthang24803/go-api/internal/database"
	"github.com/vanthang24803/go-api/routes"
	"github.com/vanthang24803/go-api/utils"
)

func init() {
	app := fiber.New()

	config := infra.GetConfig()

	err := database.Connect(config.DatabaseConnection)

	if err != nil {
		infra.Msg.Errorf("Error connecting to database: %s", err.Error())
	}

	app.Use(interceptor.Logger())
	app.Use(interceptor.ErrorHandler())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(utils.NewResponse(200, "Hello world!"))
	})

	v1 := app.Group("/api/v1")

	routes.ApplyRoutes(v1)

	app.Use(interceptor.NotFoundRoute())

	infra.Msg.Infof("Server is running on port %s 🚀", config.Port)

	app.Listen(fmt.Sprintf(":%s", config.Port))
}
