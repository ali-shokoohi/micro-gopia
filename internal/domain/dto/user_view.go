package dto

import "time"

// UserViewDto - A dto of the User type for viewing
type UserViewDto struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Age       uint8     `json:"age"`
	Email     string    `json:"email"`
}
