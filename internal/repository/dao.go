package repository

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type config struct {
	host     string
	database string
	port     string
	driver   string
	user     string
	password string
}

func InitiatePostgresDatabase() (*gorm.DB, error) {
	fmt.Println("Setting up to connect to database")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDatabaseName := os.Getenv("POSTGRES_DATABASE_NAME")

	// Starting a database
	fmt.Println("Starting postgres database")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", postgresHost, postgresPort, postgresUser, postgresDatabaseName, postgresPassword)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
