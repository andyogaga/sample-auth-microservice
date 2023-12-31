package datastruct

import (
	"time"

	"gorm.io/gorm"
)

const ProfileTableName = "profiles"

type Profile struct {
	gorm.Model
	ProfileId   string    `gorm:"primaryKey,unique"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Country     Countries `json:"country"`
}

type Countries string

const (
	NGN Countries = "NGN"
)

func MigrateProfiles(db *gorm.DB) error {
	err := db.AutoMigrate(&Profile{})
	return err
}
