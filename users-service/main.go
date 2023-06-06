package main

import (
	"fmt"
	"log"

	constants "users-service/internals/constants"
	controller "users-service/internals/controllers"
	"users-service/internals/datastruct"
	events "users-service/internals/event"
	"users-service/internals/repository"
	services "users-service/internals/services"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting the users server")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	rabbitConn, err := events.ConnectToRabbitMQ()
	if err != nil {
		log.Panic(err)
	}
	defer rabbitConn.Close()
	config := events.NewRabbitMQConfig(rabbitConn, constants.USERS_SERVICE)

	dao, err := repository.InitiatePostgresDatabase()
	if err != nil {
		log.Fatal("Encountered error connecting to users postgres database")
	}

	db := repository.GetDB()
	datastruct.MigrateProfiles(db)
	datastruct.MigrateUsers(db)

	usersService := services.NewUserService(dao)
	controller.SetupService(&usersService)
	controller.SetupGRPCRequestsListener(&config)
}
