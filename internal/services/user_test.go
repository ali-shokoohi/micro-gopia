package services_test

import (
	"fmt"
	"testing"

	"github.com/ali-shokoohi/micro-gopia/config"
	"github.com/ali-shokoohi/micro-gopia/internal/datastore"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/dto"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/entities"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/repositories"
	"github.com/ali-shokoohi/micro-gopia/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
)

func setupUserService() (services.UserService, *gin.Context, func()) {
	// Setup the config
	if err := config.Confs.Load("../../config/config-debug.yaml"); err != nil {
		er := fmt.Sprintf("Can't load the config, Error: %s\n", err)
		panic(er)
	}
	// Setup the database
	datastore := datastore.NewDatabase(sqlite.Open(":memory:"))
	datastore.Gorm.AutoMigrate(&entities.UserEntity{})
	// Create a new instance of the userRepository
	repo := repositories.NewUserRepository(datastore)
	// Create a new instance of the userService
	service := services.NewUserService(repo)

	// Create a new gin context
	ctx := &gin.Context{}

	// Return the repository, context, and a cleanup function
	cleanup := func() {
		datastore.Gorm.Migrator().DropTable(&entities.UserEntity{})
	}

	return service, ctx, cleanup
}

func TestCreateUser(t *testing.T) {
	service, ctx, cleanup := setupUserService()
	defer cleanup()

	userCreate := &dto.UserCreateDto{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "*Password123#",
	}

	user, errs := service.CreateUser(ctx, userCreate)

	assert.Empty(t, errs, "There should be no error here on CreateUser")
	assert.NotZero(t, user.ID, "User's ID shouldn't be zero")
	assert.Equal(t, userCreate.Name, user.Name, "User's names should be equal")
	assert.Equal(t, userCreate.Age, user.Age, "User's ages should be equal")
	assert.Equal(t, userCreate.Email, user.Email, "User's emails should be equal")
	assert.Equal(t, userCreate.Password, user.Password, "User's hashed passwords should be equal")
}
func TestGetUsers(t *testing.T) {
	service, ctx, cleanup := setupUserService()
	defer cleanup()

	// Insert some users into the database for testing
	usersToInsert := []*entities.UserEntity{
		{
			Name:     "John Doe",
			Age:      25,
			Email:    "john@example.com",
			Password: "*Password123#",
		},
		{
			Name:     "Jane Smith",
			Age:      30,
			Email:    "jane@example.com",
			Password: "*Password456#",
		},
	}

	for _, user := range usersToInsert {
		_, errs := service.CreateUser(ctx, &dto.UserCreateDto{
			Name:     user.Name,
			Age:      user.Age,
			Email:    user.Email,
			Password: user.Password,
		})
		assert.Empty(t, errs, "There should be no error here on CreateUser")
	}

	// Get the users
	users, err := service.GetUsers(ctx, 0, 10)

	assert.NoError(t, err, "There should be no error here on GetUsers")
	assert.Len(t, users, len(usersToInsert))
	for index, user := range users {
		assert.NotZero(t, user.ID, "User's ID shouldn't be zero")
		assert.Equal(t, usersToInsert[index].Name, user.Name, "User's names should be equal")
		assert.Equal(t, usersToInsert[index].Age, user.Age, "User's ages should be equal")
		assert.Equal(t, usersToInsert[index].Email, user.Email, "User's emails should be equal")
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(usersToInsert[index].Password)), "User's hashed passwords should be equal")
	}
}

func TestGetUserByID(t *testing.T) {
	service, ctx, cleanup := setupUserService()
	defer cleanup()

	// Insert user into the database for testing
	userToInsert := entities.UserEntity{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "*Password123#",
	}

	userCreated, errs := service.CreateUser(ctx, &dto.UserCreateDto{
		Name:     userToInsert.Name,
		Age:      userToInsert.Age,
		Email:    userToInsert.Email,
		Password: userToInsert.Password,
	})

	assert.Empty(t, errs, "There should be no error here on CreateUser")

	// Get the user by ID
	user, err := service.GetUserByID(ctx, userCreated.ID)

	assert.NoError(t, err, "There should be no error here on GetUserByID")
	assert.NotZero(t, user.ID, "User's ID shouldn't be zero")
	assert.Equal(t, userToInsert.Name, user.Name, "User's names should be equal")
	assert.Equal(t, userToInsert.Age, user.Age, "User's ages should be equal")
	assert.Equal(t, userToInsert.Email, user.Email, "User's emails should be equal")
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userToInsert.Password)), "User's hashed passwords should be equal")
}

func TesUpdateUserByID(t *testing.T) {
	service, ctx, cleanup := setupUserService()
	defer cleanup()

	// Insert a user into the database for testing
	userToInsert := entities.UserEntity{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "*Password123#",
	}

	userCreated, errs := service.CreateUser(ctx, &dto.UserCreateDto{
		Name:     userToInsert.Name,
		Age:      userToInsert.Age,
		Email:    userToInsert.Email,
		Password: userToInsert.Password,
	})

	assert.Empty(t, errs, "There should be no error here on CreateUser")

	// Update the user from the database for testing
	userToUpdate := entities.UserEntity{
		Name:     "Junior John Doe",
		Age:      20,
		Email:    "new.john@example.com",
		Password: "*NewPassword123#",
	}
	// Update the user by ID
	userUpdated, errs := service.UpdateUserByID(ctx, userCreated.ID, &dto.UserUpdateDto{
		Name:     userToUpdate.Name,
		Age:      userToUpdate.Age,
		Email:    userToUpdate.Email,
		Password: userToUpdate.Password,
	})

	assert.Empty(t, errs, "There should be no error here on UpdateUserByID")
	assert.NotZero(t, userToUpdate.ID, "User's ID shouldn't be zero")
	assert.Equal(t, userToUpdate.Name, userUpdated.Name, "User's names should be equal")
	assert.Equal(t, userToUpdate.Age, userUpdated.Age, "User's ages should be equal")
	assert.Equal(t, userToUpdate.Email, userUpdated.Email, "User's emails should be equal")
	assert.Equal(t, userToUpdate.Password, userUpdated.Password, "User's hashed passwords should be equal")
}

func TestDeleteUserByID(t *testing.T) {
	service, ctx, cleanup := setupUserService()
	defer cleanup()

	// Insert user into the database for testing
	userToInsert := entities.UserEntity{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "*Password123#",
	}

	userToInsert, errs := service.CreateUser(ctx, &dto.UserCreateDto{
		Name:     userToInsert.Name,
		Age:      userToInsert.Age,
		Email:    userToInsert.Email,
		Password: userToInsert.Password,
	})

	assert.Empty(t, errs, "There should be no error here on CreateUser")

	// Delete the user by ID
	err := service.DeleteUserByID(ctx, userToInsert.ID)

	assert.NoError(t, err, "There should be no error here on DeleteUserByID")
}
