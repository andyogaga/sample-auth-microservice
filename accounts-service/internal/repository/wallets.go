package repository

import (
	"accounts-service/internal/datastruct"
	"fmt"
)

type WalletsQuery interface {
	GetWalletById(walletId uint64) (*datastruct.Wallet, error)
}

type walletsQuery struct{}

func (w *walletsQuery) GetWalletById(walletID uint64) (*datastruct.Wallet, error) {
	walletModel := &datastruct.Wallet{ID: walletID}
	wallet := PostresDB.Model(&datastruct.Wallet{}).First(&walletModel)

	if wallet.Error != nil {
		return &datastruct.Wallet{}, fmt.Errorf("cannot get a transaction %v", wallet.Error)
	}

	return walletModel, nil
}
