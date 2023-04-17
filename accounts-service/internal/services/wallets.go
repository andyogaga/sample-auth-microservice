package services

import (
	"log"

	"accounts-service/internal/datastruct"
	"accounts-service/internal/repository"
)

type WalletService interface {
	// CreateWallet(userId int64, currency string) (string, error)
	GetWallet(walletID uint64) (*datastruct.Wallet, error)
}

type walletService struct {
	dao repository.DAO
}

func NewWalletService(dao repository.DAO) WalletService {
	return &walletService{dao: dao}
}

func (w *walletService) GetWallet(walletID uint64) (*datastruct.Wallet, error) {
	var wallet *datastruct.Wallet
	var err error

	wallet, err = w.dao.NewWalletQuery().GetWalletById(walletID)
	if err != nil {
		log.Printf("user isn't authorized %v", err)
		return nil, err
	}
	return wallet, nil
}
