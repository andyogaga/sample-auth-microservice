package datastruct

import (
	"gorm.io/gorm"
)

const UserTableName = "users"

type User struct {
	gorm.Model
	UserId    string  `json:"userId" gorm:"unique,primaryKey"`
	Phone     string  `json:"phone" gorm:"unique"`
	Email     string  `json:"email" gorm:"unique,default:null"`
	Verified  bool    `json:"verified"`
	Role      Role    `json:"role" gorm:"default:lead"`
	ProfileId string  `json:"profileId" gorm:"default:null"`
	Profile   Profile `gorm:"foreignKey:ProfileId"`
	Password  string
	PIN       string
}

type Role string

const (
	ADMIN Role = "admin"
	USER  Role = "user"
	LEAD  Role = "lead"
)

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
