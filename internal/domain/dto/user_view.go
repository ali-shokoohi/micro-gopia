package dto

import "gorm.io/gorm"

// UserViewDto - A dto of the User type for viewing
type UserViewDto struct {
	gorm.Model
	Name  string `json:"name"`
	Age   uint8  `json:"age"`
	Email string `json:"email"`
}
