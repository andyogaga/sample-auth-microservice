package services

import (
	"accounts-service/internals/datastruct"
	"accounts-service/internals/repository"
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
		return nil, err
	}
	return wallet, nil
}
