package main

import (
	"authentication-service/internal/repository"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Starting the authentication server")
	_, err := repository.InitiatePostgresDatabase()
	if err != nil {
		log.Fatal("Encountered error connecting to authentication postgres database")
	}
	app := fiber.New()
	// grpcServer := grpc.NewServer()
	// if err := grpcServer.Serve(); err != nil {
	// 	log.Fatalf("Failed to serve the Authentication GRPC server over port: %v", err)
	// }
	app.Get("/auth", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Kendi authentication server")
	})
	server_port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf(":%s", server_port))
	log.Println("Listening on port =", server_port)
}
