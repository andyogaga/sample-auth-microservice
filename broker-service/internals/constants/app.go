package constants

import "os"

type Services string

var (
	USERS_SERVICE    = os.Getenv("USERS_SERVICE")
	ACCOUNTS_SERVICE = os.Getenv("ACCOUNTS_SERVICE")
	LISTENER_SERVICE = os.Getenv("LISTENER_SERVICE")
	BROKER_SERVICE   = os.Getenv("BROKER_SERVICE")
)

const (
	GRPC_PORT = "50002"
)
