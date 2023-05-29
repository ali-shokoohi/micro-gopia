package services

import (
	"errors"
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
		return nil, append(errs, errors.New("User cannot be nil"))
	}
	if userCreate.Name == "" {
		return nil, append(errs, errors.New("User must have a name"))
	}
	if userCreate.Email == "" {
		return nil, append(errs, errors.New("User must have an email"))
	}
	if userCreate.Age < 0 {
		return nil, append(errs, errors.New("User age cannot be negative"))
	}
	// Define a regular expression pattern for matching email addresses
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// Compile the pattern into a regular expression object
	regex, err := regexp.Compile(emailPattern)
	if err != nil {
		log.Printf("An error while compiling email regex pattern: %s", err.Error())
		errs = append(errs, err)
		return nil, errs
	}
	if !regex.MatchString(userCreate.Email) {
		errs = append(errs, errors.New("invalid email address!"))
	}

	// check the password
	if !scripts.CheckPassword(userCreate.Password) {
		errs = append(errs, errors.New("invalid password!"))
	}
	if len(errs) > 0 {
		return nil, errs
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("We can't create a hash password: %s", err.Error())
		errs = append(errs, errors.New("Internal server error in creating hash password"))
		return nil, errs
	}
	userCreate.Password = string(hashedPassword)
	user, err := userService.userRepository.CreateUser(ctx, userCreate)
	return user, append(errs, err)
}

// GetUsers returns a list of users from the database based on the specified page and limits.
func (userService *userService) GetUsers(ctx *gin.Context, page, limits uint) ([]*entities.UserEntity, error) {
	if page < 0 {
		return nil, errors.New("Page shouldn't be negative!")
	}
	if limits <= 0 {
		return nil, errors.New("Limits should be bigger than zero!")
	}
	return userService.userRepository.GetUsers(ctx, page, limits)
}

// GetUserById returns a user given an id, and an error if the id is not valid.
func (userService *userService) GetUserByID(ctx *gin.Context, userID uint) (*entities.UserEntity, error) {
	if userID <= 0 {
		return nil, errors.New("Invalid user id")
	}
	return userService.userRepository.GetUserByID(ctx, userID)
}

// UpdateUser updates an existing user if the provided user is valid, and returns an error otherwise.
func (userService *userService) UpdateUserByID(ctx *gin.Context, userID uint, userUpdate *dto.UserUpdateDto) (*entities.UserEntity, []error) {
	claim, err := scripts.CurrentTokenClaim(ctx)
	if err != nil {
		log.Printf("Can't get sender ID from the request's token: %s", err.Error())
		return nil, []error{errors.New("Can't get sender ID from the request's token")}
	}
	if claim.UserID != userID {
		return nil, []error{errors.New("Permission denied")}
	}
	if userUpdate.Age < 0 {
		return nil, []error{errors.New("User age cannot be negative")}
	}
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
			errs = append(errs, err)
			return nil, errs
		}
		if !regex.MatchString(userUpdate.Email) {
			errs = append(errs, errors.New("invalid email address!"))
		}
	}

	if userUpdate.Password != "" {
		// check the password
		if !scripts.CheckPassword(userUpdate.Password) {
			errs = append(errs, errors.New("invalid password!"))
		}
		if len(errs) > 0 {
			return nil, errs
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("We can't create a hash password: %s", err.Error())
			errs = append(errs, errors.New("Internal server error in creating hash password"))
			return nil, errs
		}
		userUpdate.Password = string(hashedPassword)
	}

	if len(errs) > 0 {
		return nil, errs
	}
	user, err = userService.userRepository.UpdateUserByID(ctx, userID, userUpdate, user)
	return user, append(errs, err)
}

// DeleteUser deletes an existing user given an id, and returns an error if the id is not valid.
func (userService *userService) DeleteUserByID(ctx *gin.Context, userID uint) error {
	if userID <= 0 {
		return errors.New("Invalid user id")
	}
	claim, err := scripts.CurrentTokenClaim(ctx)
	if err != nil {
		log.Printf("Can't get sender ID from the request's token: %s", err.Error())
		return errors.New("Can't get sender ID from the request's token")
	}
	if claim.UserID != userID {
		return errors.New("Permission denied")
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
		return "", errors.New("Invalid password!")
	}
	// Generate the token
	tk := &dto.Claims{UserID: user.ID}
	tk.ExpiresAt = time.Now().Unix() + int64(config.Confs.Service.Token.Expiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, err := token.SignedString([]byte(config.Confs.Service.Token.Password))
	if err != nil {
		log.Printf("An Error while Generating a user's token string: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
