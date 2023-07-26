package main

import (
	"fmt"
	"log"
	"os"

	constants "users-service/internals/constants"
	controller "users-service/internals/controllers"
	"users-service/internals/datastruct"
	events "users-service/internals/event"
	"users-service/internals/repository"
	services "users-service/internals/services"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting the users server")
	if err := godotenv.Load(os.Getenv("SERVICE")); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	rabbitConn, err := events.ConnectToRabbitMQ()
	if err != nil {
		log.Panic(err)
	}
	defer rabbitConn.Close()
	config := events.NewRabbitMQConfig(rabbitConn, constants.USERS_SERVICE)

	dbConfig := repository.DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DATABASE_NAME"),
	}

	dao, err := repository.InitiatePostgresDatabase(&dbConfig, gorm.Open)
	if err != nil {
		log.Fatal("Encountered error connecting to users postgres database")
	}

	datastruct.MigrateProfiles(dao.PostgresDB)
	datastruct.MigrateUsers(dao.PostgresDB)

	profilesService := services.NewProfileService(dao)
	usersService := services.NewUserService(dao, profilesService)
	server := controller.SetupService(usersService)
	go controller.SetupGRPCRequestsListener(&config, server)
	fmt.Println("users service started successfully")
}
