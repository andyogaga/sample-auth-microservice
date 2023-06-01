package main

import (
	"fmt"
	"log"

	"users-service/internals/datastruct"
	requests "users-service/internals/proto"
	"users-service/internals/repository"
	services "users-service/internals/services"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting the users server")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	dao, err := repository.InitiatePostgresDatabase()
	if err != nil {
		log.Fatal("Encountered error connecting to users postgres database")
	}

	db := repository.GetDB()

	datastruct.MigrateUsers(db)

	usersService := services.NewUserService(dao)
	requests.SetupService(&usersService)
	requests.SetupGRPCRequestsListener()
}
