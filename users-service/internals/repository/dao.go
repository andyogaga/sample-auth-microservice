package repository

import (
	"fmt"
	"net/http"
	"users-service/internals/utils"

	"gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
)

type DAO interface {
	NewUserQuery() UsersQuery
	NewProfileQuery() IProfilesQuery
	BeginTransaction() *gorm.DB
}

type dao struct {
	PostgresDB *gorm.DB
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

type TGormOpen func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error)

func InitiatePostgresDatabase(config *DBConfig, gormOpen TGormOpen) (*dao, error) {
	fmt.Println("Setting up to connect to database")
	// Starting a database
	fmt.Println("Starting postgres database")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", config.Host, config.Port, config.User, config.DbName, config.Password)

	PostgresDB, err := gormOpen(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return nil, utils.NewError(http.StatusInternalServerError, "unexpected error", err)
	}
	return &dao{
		PostgresDB: PostgresDB,
	}, nil
}

func (d *dao) BeginTransaction() *gorm.DB {
	return d.PostgresDB.Begin()
}

func (d *dao) NewUserQuery() UsersQuery {
	return &usersQuery{
		PostgresDB: d.PostgresDB,
	}
}

func (d *dao) NewProfileQuery() IProfilesQuery {
	return &profilesQuery{
		PostgresDB: d.PostgresDB,
	}
}
