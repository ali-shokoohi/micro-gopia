package services

import (
	"log"
	"regexp"
	"time"

	"github.com/ali-shokoohi/micro-gopia/config"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/dto"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/entities"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/repositories"
	"github.com/ali-shokoohi/micro-gopia/scripts"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx *gin.Context, userCreate *dto.UserCreateDto) (*entities.UserEntity, []error)
	GetUsers(ctx *gin.Context, page, limits uint) ([]*entities.UserEntity, error)
	GetUserByID(ctx *gin.Context, userID uint) (*entities.UserEntity, error)
	UpdateUserByID(ctx *gin.Context, userID uint, userUpdate *dto.UserUpdateDto) (*entities.UserEntity, []error)
	DeleteUserByID(ctx *gin.Context, userID uint) error
	Login(ctx *gin.Context, userLogin *dto.UserLoginDto) (string, error)
}

// UserService represents the service that interacts with the user_repository.
type userService struct {
	userRepository repositories.UserRepository
}

// NewUserService returns a new instance of UserService.
func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

// CreateUser creates a new user if the provided user is valid, and returns an error otherwise.
func (userService *userService) CreateUser(ctx *gin.Context, userCreate *dto.UserCreateDto) (*entities.UserEntity, []error) {
	var errs []error
	if userCreate == nil {
		return nil, append(errs, dto.BadRequestError{Message: "User cannot be null"})
	}
	if userCreate.Name == "" {
		errs = append(errs, dto.BadRequestError{Message: "User must have a name"})
	}
	if userCreate.Email == "" {
		errs = append(errs, dto.BadRequestError{Message: "User must have an email"})
	}
	if userCreate.Age < 0 {
		errs = append(errs, dto.BadRequestError{Message: "User age cannot be negative"})
	}
	// Define a regular expression pattern for matching email addresses
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// Compile the pattern into a regular expression object
	regex, err := regexp.Compile(emailPattern)
	if err != nil {
		log.Printf("An error while compiling email regex pattern: %s", err.Error())
		errs = append(errs, dto.InternalServerError{Message: "Can't compiling email regex pattern"})
		return nil, errs
	}
	if !regex.MatchString(userCreate.Email) {
		errs = append(errs, dto.BadRequestError{Message: "invalid email address!"})
	}

	// check the password
	if !scripts.CheckPassword(userCreate.Password) {
		errs = append(errs, dto.BadRequestError{Message: "invalid password!"})
	}
	if len(errs) > 0 {
		return nil, errs
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("We can't create a hash password: %s", err.Error())
		errs = append(errs, dto.InternalServerError{Message: "Internal server error in creating hash password"})
		return nil, errs
	}
	userCreate.Password = string(hashedPassword)
	user, err := userService.userRepository.CreateUser(ctx, userCreate)
	if err != nil {
		return nil, append(errs, err)
	}
	return user, nil
}

// GetUsers returns a list of users from the database based on the specified page and limits.
func (userService *userService) GetUsers(ctx *gin.Context, page, limits uint) ([]*entities.UserEntity, error) {
	if page < 0 {
		return nil, dto.BadRequestError{Message: "Page shouldn't be negative!"}
	}
	if limits <= 0 {
		return nil, dto.BadRequestError{Message: "Limits should be bigger than zero!"}
	}
	return userService.userRepository.GetUsers(ctx, page, limits)
}

// GetUserById returns a user given an id, and an error if the id is not valid.
func (userService *userService) GetUserByID(ctx *gin.Context, userID uint) (*entities.UserEntity, error) {
	if userID <= 0 {
		return nil, dto.BadRequestError{Message: "Invalid user id"}
	}
	return userService.userRepository.GetUserByID(ctx, userID)
}

// UpdateUser updates an existing user if the provided user is valid, and returns an error otherwise.
func (userService *userService) UpdateUserByID(ctx *gin.Context, userID uint, userUpdate *dto.UserUpdateDto) (*entities.UserEntity, []error) {
	var errs []error
	user, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, append(errs, err)
	}

	if userUpdate.Email != "" { // Define a regular expression pattern for matching email addresses
		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		// Compile the pattern into a regular expression object
		regex, err := regexp.Compile(emailPattern)
		if err != nil {
			log.Printf("An error while compiling email regex pattern: %s", err.Error())
			errs = append(errs, dto.InternalServerError{Message: "Can't compiling email regex pattern"})
			return nil, errs
		}
		if !regex.MatchString(userUpdate.Email) {
			errs = append(errs, dto.BadRequestError{Message: "invalid email address!"})
		}
	}

	if userUpdate.Age < 0 {
		errs = append(errs, dto.BadRequestError{Message: "User age cannot be negative"})
	}

	if userUpdate.Password != "" {
		// check the password
		if !scripts.CheckPassword(userUpdate.Password) {
			errs = append(errs, dto.BadRequestError{Message: "invalid password!"})
		}
		if len(errs) > 0 {
			return nil, errs
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("We can't create a hash password: %s", err.Error())
			errs = append(errs, dto.InternalServerError{Message: "Internal server error in creating hash password"})
			return nil, errs
		}
		userUpdate.Password = string(hashedPassword)
	}

	if len(errs) > 0 {
		return nil, errs
	}
	user, err = userService.userRepository.UpdateUserByID(ctx, userID, userUpdate, user)
	if err != nil {
		return nil, append(errs, err)
	}
	return user, nil
}

// DeleteUser deletes an existing user given an id, and returns an error if the id is not valid.
func (userService *userService) DeleteUserByID(ctx *gin.Context, userID uint) error {
	if userID <= 0 {
		return dto.BadRequestError{Message: "Invalid user id"}
	}
	return userService.userRepository.DeleteUserByID(ctx, userID)
}

// Login verifies logins and return a token string
func (userService *userService) Login(ctx *gin.Context, userLogin *dto.UserLoginDto) (string, error) {
	user, err := userService.userRepository.GetUserByEmail(ctx, userLogin.Email)
	if err != nil {
		log.Printf("An Error while getting a user with email: %s", err.Error())
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		return "", dto.BadRequestError{Message: "Invalid password!"}
	}
	// Generate the token
	tk := &dto.Claims{UserID: user.ID}
	tk.ExpiresAt = time.Now().Unix() + int64(config.Confs.Service.Token.Expiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, err := token.SignedString([]byte(config.Confs.Service.Token.Password))
	if err != nil {
		log.Printf("An Error while Generating a user's token string: %s", err.Error())
		return "", dto.InternalServerError{Message: "An Error while Generating a user's token string"}
	}
	return tokenString, nil
}
