package datastruct

import (
	"gorm.io/gorm"
)

const UserTableName = "users"

type User struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey"`
	Phone     string `json:"phone" gorm:"unique"`
	Email     string `json:"email" gorm:"unique"`
	UserId    string `json:"userId" gorm:"unique"`
	Verified  bool   `json:"verified"`
	Role      Role   `json:"role"`
	ProfileId string
	Profile   Profile `gorm:"foreignKey:ProfileId"`
}

type Role string

const (
	ADMIN Role = "admin"
	USER  Role = "user"
)

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
