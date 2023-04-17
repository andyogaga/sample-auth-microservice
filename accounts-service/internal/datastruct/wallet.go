package datastruct

import (
	"time"

	"gorm.io/gorm"
)

const WalletTableName = "wallets"

type Wallet struct {
	ID        uint64    `gorm:"primaryKey"`
	UserId    string    `json:"userId"`
	Currency  string    `json:"currency"`
	Balance   string    `json:"balance"`
	Valid     bool      `json:"valid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MigreateWallets(db *gorm.DB) error {
	err := db.AutoMigrate(&Wallet{})
	return err
}
