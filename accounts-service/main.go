package main

import (
	"fmt"
	"log"

	"accounts-service/internals/datastruct"
	requests "accounts-service/internals/proto"
	"accounts-service/internals/repository"
	services "accounts-service/internals/services"
)

func main() {
	fmt.Println("Starting the accounts server")
	dao, err := repository.InitiatePostgresDatabase()
	if err != nil {
		log.Fatal("Encountered error connecting to accounts postgres database")
	}

	db := repository.GetDB()

	datastruct.MigrateWallets(db)

	walletService := services.NewWalletService(dao)
	requests.SetupService(&walletService)
	requests.SetupGRPCRequestsListener()
}
