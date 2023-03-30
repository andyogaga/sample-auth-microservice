package main

import (
	"fmt"
	"kendi-api/internal/repository"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Starting the server")
	_, err := repository.InitiatePostgresDatabase()
	if err != nil {
		log.Fatal("Encountered error connecting to postgres database")
	}
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Kendi server")
	})
	server_port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf(":%s", server_port))
	log.Println("Listening on port =", server_port)
}
