package datastruct

import (
	"time"

	"gorm.io/gorm"
)

const UserTableName = "users"

type User struct {
	ID        uint64    `gorm:"primaryKey"`
	UserId    string    `json:"userId"`
	Country   string    `json:"country"`
	Valid     bool      `json:"valid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
