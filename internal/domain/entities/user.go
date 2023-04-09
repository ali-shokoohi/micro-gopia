package entities

import (
	"gorm.io/gorm"
)

// UserEntity - An database entity of the User type
type UserEntity struct {
	gorm.Model
	Name     string `gorm:""`
	Age      uint8  `gorm:""`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}
