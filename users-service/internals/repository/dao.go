package repository

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DAO interface {
	NewUserQuery() UsersQuery
	NewProfileQuery() ProfilesQuery
}

var PostresDB *gorm.DB

type dao struct{}

func GetDB() *gorm.DB {
	return PostresDB
}

func InitiatePostgresDatabase() (*dao, error) {
	fmt.Println("Setting up to connect to database")

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDatabaseName := os.Getenv("POSTGRES_DATABASE_NAME")

	// Starting a database
	fmt.Println("Starting postgres database")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", postgresHost, postgresPort, postgresUser, postgresDatabaseName, postgresPassword)

	var err error
	PostresDB, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &dao{}, nil
}

func (d *dao) NewUserQuery() UsersQuery {
	return &usersQuery{}
}

func (d *dao) NewProfileQuery() ProfilesQuery {
	return &profilesQuery{}
}
