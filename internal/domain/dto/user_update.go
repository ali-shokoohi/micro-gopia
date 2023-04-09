package dto

// UserViewDto - A dto of the User type for updating
type UserUpdateDto struct {
	Name     string `json:"name"`
	Age      uint8  `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
