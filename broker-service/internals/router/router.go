package router

import (
	"github.com/gofiber/fiber/v2"

	events "broker-service/internals/event"
	utils "broker-service/internals/utils"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func CreateRouter(config events.Config) *fiber.App {
	app := fiber.New()

	// middlewares
	app.Use(idempotency.New())
	app.Use(cors.New(cors.Config{
		// AllowOrigins: "https://gofiber.io, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(recover.New())
	app.Use(func(c *fiber.Ctx) error {
		utils.LogRequest(c, config)
		utils.SetupHttpServerInstance(c)
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the kendi server")
	})

	return app
}
