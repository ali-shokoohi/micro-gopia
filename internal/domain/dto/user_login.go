package dto

// UserViewDto - A dto of the User type for creation
type UserLoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
