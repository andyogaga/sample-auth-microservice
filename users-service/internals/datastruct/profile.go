package datastruct

import (
	"time"

	"gorm.io/gorm"
)

const ProfileTableName = "profiles"

type Profile struct {
	gorm.Model
	ID          uint64    `gorm:"primaryKey"`
	UserId      string    `json:"userId"`
	Firstname   string    `json:"firstName"`
	Lastname    bool      `json:"lastName"`
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
