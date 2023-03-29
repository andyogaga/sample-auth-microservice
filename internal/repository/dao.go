package repository

import (
	"fmt"
)

type config struct {
	host     string
	database string
	port     string
	driver   string
	user     string
	password string
}

func InitiateDatabase() /**(*sql.DB, error) */ {
	fmt.Println("Setting up to connect to database")
}
