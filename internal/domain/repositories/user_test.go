package repositories_test

import (
	"testing"

	"github.com/ali-shokoohi/micro-gopia/internal/datastore"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/dto"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/entities"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/repositories"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
)

func setupUserRepository() (repositories.UserRepository, *gin.Context, func()) {
	// Setup the database
	datastore := datastore.NewDatabase(sqlite.Open(":memory:"))
	datastore.Gorm.AutoMigrate(&entities.UserEntity{})
	// Create a new instance of the userRepository
	repo := repositories.NewUserRepository(datastore)

	// Create a new gin context
	ctx := &gin.Context{}

	// Return the repository, context, and a cleanup function
	cleanup := func() {
		datastore.Gorm.Migrator().DropTable(&entities.UserEntity{})
	}

	return repo, ctx, cleanup
}

func TestCreateUser(t *testing.T) {
	repo, ctx, cleanup := setupUserRepository()
	defer cleanup()

	userCreate := &dto.UserCreateDto{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "*Password123#",
	}

	user := &entities.UserEntity{}
	err := repo.CreateUser(ctx, userCreate, user)

	assert.NoError(t, err, "There should be no error here on CreateUser")
	assert.NotZero(t, user.ID, "User's ID shouldn't be zero")
	assert.Equal(t, userCreate.Name, user.Name, "User's names should be equal")
	assert.Equal(t, userCreate.Age, user.Age, "User's ages should be equal")
	assert.Equal(t, userCreate.Email, user.Email, "User's emails should be equal")
	assert.Equal(t, userCreate.Password, user.Password, "User's passwords should be equal")
}

func TestGetUsers(t *testing.T) {
	repo, ctx, cleanup := setupUserRepository()
	defer cleanup()

	// Insert some users into the database for testing
	usersToInsert := []*entities.UserEntity{
		{
			Name:     "John Doe",
			Age:      25,
			Email:    "john@example.com",
			Password: "password123",
		},
		{
			Name:     "Jane Smith",
			Age:      30,
			Email:    "jane@example.com",
			Password: "password456",
		},
	}

	for _, user := range usersToInsert {
		err := repo.CreateUser(ctx, &dto.UserCreateDto{
			Name:     user.Name,
			Age:      user.Age,
			Email:    user.Email,
			Password: user.Password,
		}, user)
		assert.NoError(t, err, "There should be no error here on CreateUser")
	}

	// Get the users
	var users []*entities.UserEntity
	err := repo.GetUsers(ctx, 0, 10, &users)

	assert.NoError(t, err, "There should be no error here on GetUsers")
	assert.Len(t, users, len(usersToInsert))
	for index, user := range users {
		assert.NotZero(t, user.ID, "User's ID shouldn't be zero")
		assert.Equal(t, usersToInsert[index].Name, user.Name, "User's names should be equal")
		assert.Equal(t, usersToInsert[index].Age, user.Age, "User's ages should be equal")
		assert.Equal(t, usersToInsert[index].Email, user.Email, "User's emails should be equal")
		assert.Equal(t, usersToInsert[index].Password, user.Password, "User's passwords should be equal")
	}
}

func TestGetUserByID(t *testing.T) {
	repo, ctx, cleanup := setupUserRepository()
	defer cleanup()

	// Insert user into the database for testing
	userToInsert := &entities.UserEntity{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "password123",
	}

	err := repo.CreateUser(ctx, &dto.UserCreateDto{
		Name:     userToInsert.Name,
		Age:      userToInsert.Age,
		Email:    userToInsert.Email,
		Password: userToInsert.Password,
	}, userToInsert)

	assert.NoError(t, err, "There should be no error here on CreateUser")

	// Get the user by ID
	var user entities.UserEntity
	err = repo.GetUserByID(ctx, userToInsert.ID, &user)

	assert.NoError(t, err, "There should be no error here on GetUserByID")
	assert.NotZero(t, user.ID, "User's ID shouldn't be zero")
	assert.Equal(t, userToInsert.Name, user.Name, "User's names should be equal")
	assert.Equal(t, userToInsert.Age, user.Age, "User's ages should be equal")
	assert.Equal(t, userToInsert.Email, user.Email, "User's emails should be equal")
	assert.Equal(t, userToInsert.Password, user.Password, "User's passwords should be equal")
}

func TesUpdateUserByID(t *testing.T) {
	repo, ctx, cleanup := setupUserRepository()
	defer cleanup()

	// Insert a user into the database for testing
	userToInsert := &entities.UserEntity{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "password123",
	}

	err := repo.CreateUser(ctx, &dto.UserCreateDto{
		Name:     userToInsert.Name,
		Age:      userToInsert.Age,
		Email:    userToInsert.Email,
		Password: userToInsert.Password,
	}, userToInsert)

	assert.NoError(t, err, "There should be no error here on CreateUser")

	// Update the user from the database for testing
	userToUpdate := &entities.UserEntity{
		Name:     "Junior John Doe",
		Age:      20,
		Email:    "new.john@example.com",
		Password: "newPassword123",
	}
	// Update the user by ID
	err = repo.UpdateUserByID(ctx, userToInsert.ID, &dto.UserUpdateDto{
		Name:     userToUpdate.Name,
		Age:      userToUpdate.Age,
		Email:    userToUpdate.Email,
		Password: userToUpdate.Password,
	}, userToUpdate)

	assert.NoError(t, err, "There should be no error here on UpdateUserByID")
	assert.NotZero(t, userToUpdate.ID, "User's ID shouldn't be zero")
	assert.Equal(t, userToInsert.Name, userToUpdate.Name, "User's names should be equal")
	assert.Equal(t, userToInsert.Age, userToUpdate.Age, "User's ages should be equal")
	assert.Equal(t, userToInsert.Email, userToUpdate.Email, "User's emails should be equal")
	assert.Equal(t, userToInsert.Password, userToUpdate.Password, "User's passwords should be equal")
}

func TestDeleteUserByID(t *testing.T) {
	repo, ctx, cleanup := setupUserRepository()
	defer cleanup()

	// Insert user into the database for testing
	userToInsert := &entities.UserEntity{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "password123",
	}

	err := repo.CreateUser(ctx, &dto.UserCreateDto{
		Name:     userToInsert.Name,
		Age:      userToInsert.Age,
		Email:    userToInsert.Email,
		Password: userToInsert.Password,
	}, userToInsert)

	assert.NoError(t, err, "There should be no error here on CreateUser")

	// Delete the user by ID
	err = repo.DeleteUserByID(ctx, userToInsert.ID)

	assert.NoError(t, err, "There should be no error here on DeleteUserByID")
}
