package main

import (
	"accounts-service/internal/datastruct"
	"accounts-service/internal/repository"
	services "accounts-service/internal/services"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Starting the accounts server")
	dao, err := repository.InitiatePostgresDatabase()
	if err != nil {
		log.Fatal("Encountered error connecting to accounts postgres database")
	}

	db := repository.GetDB()

	datastruct.MigreateWallets(db)

	walletService := services.NewWalletService(dao)

	app := fiber.New()
	grpcServer := grpc.NewServer()
	if err := grpcServer.Serve(); err != nil {
		log.Fatalf("Failed to serve the accounts GRPC server over port: %v", err)
	}

	server_port := os.Getenv("PORT")
	app.Listen(fmt.Sprintf(":%s", server_port))
	log.Println("Listening on port =", server_port)
}
