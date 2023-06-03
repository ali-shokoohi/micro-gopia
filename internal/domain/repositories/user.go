package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/ali-shokoohi/micro-gopia/internal/datastore"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/dto"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserRepository defines five methods for creating, reading, updating, and deleting user data from the database
type UserRepository interface {
	CreateUser(ctx *gin.Context, userCreate *dto.UserCreateDto, user *entities.UserEntity) error
	GetUsers(ctx *gin.Context, page, limits uint, users *[]*entities.UserEntity) error
	GetUserByID(ctx *gin.Context, userID uint, user *entities.UserEntity) error
	GetUserByEmail(ctx *gin.Context, email string, user *entities.UserEntity) error
	UpdateUserByID(ctx *gin.Context, userID uint, userUpdate *dto.UserUpdateDto, userEntity *entities.UserEntity) error
	DeleteUserByID(ctx *gin.Context, userID uint) error
}

// userRepository implements the interface by defining methods that interact with the database using gorm.
type userRepository struct {
	db     *datastore.Database // Assuming you have a struct named Datastore that implements the necessary DB methods
	gormDB *gorm.DB
}

// NewUserRepository returns a new instance of userRepository
func NewUserRepository(db *datastore.Database) UserRepository {
	return &userRepository{db: db, gormDB: db.GetDatabase()}
}

// CreateUser validates user data such as email and password using regular expressions, generates a hashed password,
// And then creates and inserts a new user into the database using gorm.
func (r *userRepository) CreateUser(ctx *gin.Context, userCreate *dto.UserCreateDto, user *entities.UserEntity) error {
	user.Name = userCreate.Name
	user.Age = userCreate.Age
	user.Email = userCreate.Email
	user.Password = userCreate.Password

	// Insert the user into the database
	if err := r.gormDB.Create(user).Error; err != nil {
		log.Printf("We can't insert new user into database: %s", err.Error())

		return dto.InternalServerError{Message: "Internal error at inserting new user into database"}
	}

	return nil
}

// GetUsers returns a list of users from the database based on the specified page and limits.
func (r *userRepository) GetUsers(ctx *gin.Context, page, limits uint, users *[]*entities.UserEntity) error {
	// Query the users from the database
	offset := page * limits
	if err := r.gormDB.Limit(int(limits)).Offset(int(offset)).Find(&users).Error; err != nil {

		log.Printf("We can't get Users Error: %s", err.Error())
		return dto.InternalServerError{Message: "Internal server error at getting users"}
	}

	return nil
}

// GetUserByID returns a user from the database based on the specified userID.
func (r *userRepository) GetUserByID(ctx *gin.Context, userID uint, user *entities.UserEntity) error {

	// Query the user by ID from the database
	if err := r.gormDB.First(user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BadRequestError{Message: fmt.Sprintf("User not found with ID: %d", userID)}
		}
		log.Printf("We can't get User with ID: %d, Error: %s", userID, err.Error())
		return dto.InternalServerError{Message: "Internal server error at getting user with ID"}
	}

	return nil
}

// GetUserByEmail returns a user from the database based on the specified email.
func (r *userRepository) GetUserByEmail(ctx *gin.Context, email string, user *entities.UserEntity) error {
	user.Email = email

	// Query the user by Email from the database
	if err := r.gormDB.Where(user).Find(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BadRequestError{Message: fmt.Sprintf("User not found with Email: %s", email)}
		}
		log.Printf("We can't get User with Email: %s, Error: %s", email, err.Error())
		return dto.InternalServerError{Message: "Internal server error at getting user with Email"}
	}

	return nil
}

// UpdateUserByID updates a user in the database based on the specified userID and userUpdate data,
// Which may include a new email, password, age or name.
func (r *userRepository) UpdateUserByID(ctx *gin.Context, userID uint, userUpdate *dto.UserUpdateDto, userEntity *entities.UserEntity) error {

	// Update the user's properties
	if userUpdate.Name != "" {
		userEntity.Name = userUpdate.Name
	}
	if userUpdate.Age != 0 {
		userEntity.Age = userUpdate.Age
	}
	if userUpdate.Email != "" {
		userEntity.Email = userUpdate.Email
	}
	if userUpdate.Password != "" {
		userEntity.Password = userUpdate.Password
	}

	// Save the updated user to the database
	if err := r.gormDB.Save(userEntity).Error; err != nil {
		log.Printf("Internal error at saving updated user into database. UserID: %d, Error: %s", userID, err.Error())
		return dto.InternalServerError{Message: "Internal error at saving updated user into database"}
	}

	return nil
}

// DeleteUserByID deletes a user from the database based on the specified userID.
func (r *userRepository) DeleteUserByID(ctx *gin.Context, userID uint) error {
	user := &entities.UserEntity{}

	// Query the user by ID from the database
	if err := r.GetUserByID(ctx, userID, user); err != nil {
		return err
	}

	// Delete the user from the database
	if err := r.gormDB.Delete(user).Error; err != nil {
		log.Printf("Internal error at deleting a user from database. UserID: %d, Error: %s", userID, err.Error())
		return dto.InternalServerError{Message: "Internal error at deleting a user from database"}
	}

	return nil
}
