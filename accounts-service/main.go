package main

import (
	"fmt"
	"log"

	"accounts-service/internal/datastruct"
	requests "accounts-service/internal/proto"
	"accounts-service/internal/repository"
	services "accounts-service/internal/services"
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
