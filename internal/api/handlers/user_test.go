package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ali-shokoohi/micro-gopia/config"
	"github.com/ali-shokoohi/micro-gopia/internal/api/handlers"
	"github.com/ali-shokoohi/micro-gopia/internal/datastore"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/dto"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/entities"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/repositories"
	"github.com/ali-shokoohi/micro-gopia/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
)

func setupUserHandler() (handlers.UserHandler, func()) {
	// Setup the config
	if err := config.Confs.Load("../../../config/config-debug.yaml"); err != nil {
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
	// Create a new instance of the userHandler
	handler := handlers.NewUserHandler(service)

	// Return the repository, context, and a cleanup function
	cleanup := func() {
		datastore.Gorm.Migrator().DropTable(&entities.UserEntity{})
	}

	return handler, cleanup
}

func TestCreateUser(t *testing.T) {
	userHandler, cleanup := setupUserHandler()
	defer cleanup()

	// Create a test HTTP request with a JSON payload
	userCreate := &dto.UserCreateDto{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Age:      24,
		Password: "*JohnPassword24#",
	}
	payload, err := json.Marshal(userCreate)
	assert.NoError(t, err, "There should be no error here on json.Marshal(userCreate)")
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(payload))
	assert.NoError(t, err, "There should be no error here on making creation request")

	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the CreateUser handler method
	userHandler.CreateUser(c)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, w.Code, "Response's status code should be 201 (Created)")

	// Validation the response body
	var responseBody dto.HttpUserSuccess
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")
	assert.NotZero(t, responseBody.User.ID, "Created user's ID shouldn't be zero")
	assert.Equal(t, userCreate.Name, responseBody.User.Name, "User's names should be equal")
	assert.Equal(t, userCreate.Age, responseBody.User.Age, "User's ages should be equal")
	assert.Equal(t, userCreate.Email, responseBody.User.Email, "User's emails should be equal")
}

func TestGetUsers(t *testing.T) {
	// Create a new instance of the user handler
	userHandler, cleanup := setupUserHandler()
	defer cleanup()

	// Create a test HTTP request with a JSON payload
	usersToInsert := []*dto.UserCreateDto{
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
		payload, err := json.Marshal(user)
		assert.NoError(t, err, "There should be no error here on json.Marshal(userCreate)")
		req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(payload))
		assert.NoError(t, err, "There should be no error here on making creation request")
		// Create a test HTTP response recorder
		w := httptest.NewRecorder()

		// Create a Gin context from the request and response recorder
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the CreateUser handler method
		userHandler.CreateUser(c)

		// Check the response status code
		assert.Equal(t, http.StatusCreated, w.Code, "Response's status code should be 201 (Created)")
	}

	// Create a test HTTP request with query parameters
	req, err := http.NewRequest(http.MethodGet, "/api/v1/users?page=0&limits=10", nil)
	assert.NoError(t, err, "There should be no error here on making getting users request")

	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the GetUsers handler method
	userHandler.GetUsers(c)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code, "Response's status code should be 200 (ok)")

	// Validation the response body
	var responseBody dto.HttpUsersSuccess
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")
	assert.NotEmpty(t, responseBody.Users, "Getting users shouldn't be empty")
	for index, user := range responseBody.Users {
		assert.NotZero(t, user.ID, "User's ID shouldn't be zero")
		assert.Equal(t, usersToInsert[index].Name, user.Name, "User's names should be equal")
		assert.Equal(t, usersToInsert[index].Age, user.Age, "User's ages should be equal")
		assert.Equal(t, usersToInsert[index].Email, user.Email, "User's emails should be equal")
	}
}

func TestGetUserByID(t *testing.T) {
	// Create a new instance of the user handler
	userHandler, cleanup := setupUserHandler()
	defer cleanup()

	// Create a test HTTP request with a JSON payload
	userToInsert := dto.UserCreateDto{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "*Password123#",
	}

	payload, err := json.Marshal(userToInsert)
	assert.NoError(t, err, "There should be no error here on json.Marshal(userToInsert)")
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(payload))
	assert.NoError(t, err, "There should be no error here on making creation request")
	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the CreateUser handler method
	userHandler.CreateUser(c)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, w.Code, "Response's status code should be 201 (Created)")
	var createdUserResponse dto.HttpUserSuccess
	err = json.Unmarshal(w.Body.Bytes(), &createdUserResponse)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")

	// Create a test HTTP request with query parameters
	userIDString := strconv.Itoa(int(createdUserResponse.User.ID))
	req, err = http.NewRequest(http.MethodGet, "/api/v1/users/"+userIDString, nil)
	assert.NoError(t, err, "There should be no error here on making getting user by ID request")

	// Create a test HTTP response recorder
	w = httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: userIDString}}

	// Call the GetUserByID handler method
	userHandler.GetUserByID(c)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code, "Response's status code should be 200 (ok)")

	// Validation the response body
	var responseBody dto.HttpUserSuccess
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")
	assert.NotZero(t, responseBody.User.ID, "User's ID shouldn't be zero")
	assert.Equal(t, userToInsert.Name, responseBody.User.Name, "User's names should be equal")
	assert.Equal(t, userToInsert.Age, responseBody.User.Age, "User's ages should be equal")
	assert.Equal(t, userToInsert.Email, responseBody.User.Email, "User's emails should be equal")
}

func TestUpdateUserByID(t *testing.T) {
	// Create a new instance of the user handler
	userHandler, cleanup := setupUserHandler()
	defer cleanup()

	// Create a test HTTP request with a JSON payload
	userToInsert := dto.UserCreateDto{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "*Password123#",
	}

	payload, err := json.Marshal(userToInsert)
	assert.NoError(t, err, "There should be no error here on json.Marshal(userToInsert)")
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(payload))
	assert.NoError(t, err, "There should be no error here on making creation request")
	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the CreateUser handler method
	userHandler.CreateUser(c)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, w.Code, "Response's status code should be 201 (Created)")
	var createdUserResponse dto.HttpUserSuccess
	err = json.Unmarshal(w.Body.Bytes(), &createdUserResponse)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")

	// Create a test HTTP request for login
	userToLogin := dto.UserLoginDto{
		Email:    userToInsert.Email,
		Password: userToInsert.Password,
	}
	payload, err = json.Marshal(userToLogin)
	assert.NoError(t, err, "There should be no error here on json.Marshal(userToLogin)")
	req, err = http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(payload))
	assert.NoError(t, err, "There should be no error here on making login user request")
	// Create a test HTTP response recorder
	w = httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ = gin.CreateTestContext(w)
	c.Request = req

	// Call the UpdateUserByID handler method
	userHandler.Login(c)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code, "Response's status code should be 200 (ok)")

	// Validation the response body
	var loginResponseBody dto.HttpAccessTokenSuccess
	err = json.Unmarshal(w.Body.Bytes(), &loginResponseBody)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")
	assert.NotEmpty(t, loginResponseBody.Token, "Access token shouldn't be empty")
	// Create a test HTTP request with query parameter and payload
	userToUpdate := dto.UserCreateDto{
		Name:     "Junior John Doe",
		Age:      22,
		Email:    "new.john@example.com",
		Password: "*NewPassword123#",
	}
	payload, err = json.Marshal(userToUpdate)
	assert.NoError(t, err, "There should be no error here on json.Marshal(userToUpdate)")
	userIDString := strconv.Itoa(int(createdUserResponse.User.ID))
	req, err = http.NewRequest(http.MethodPut, "/api/v1/users/"+userIDString, bytes.NewBuffer(payload))
	assert.NoError(t, err, "There should be no error here on making Updating user by ID request")

	// Create a test HTTP response recorder
	w = httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: userIDString}}
	c.Request.Header.Set("Authorization", "Bearer "+loginResponseBody.Token)

	// Call the UpdateUserByID handler method
	userHandler.UpdateUserByID(c)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code, "Response's status code should be 200 (ok)")

	// Validation the response body
	var responseBody dto.HttpUserSuccess
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")
	assert.NotZero(t, responseBody.User.ID, "User's ID shouldn't be zero")
	assert.Equal(t, userToUpdate.Name, responseBody.User.Name, "User's names should be equal")
	assert.Equal(t, userToUpdate.Age, responseBody.User.Age, "User's ages should be equal")
	assert.Equal(t, userToUpdate.Email, responseBody.User.Email, "User's emails should be equal")
}

func TestDeleteUserByID(t *testing.T) {
	// Create a new instance of the user handler
	userHandler, cleanup := setupUserHandler()
	defer cleanup()

	// Create a test HTTP request with a JSON payload
	userToInsert := dto.UserCreateDto{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "*Password123#",
	}

	payload, err := json.Marshal(userToInsert)
	assert.NoError(t, err, "There should be no error here on json.Marshal(userToInsert)")
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(payload))
	assert.NoError(t, err, "There should be no error here on making creation request")
	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the CreateUser handler method
	userHandler.CreateUser(c)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, w.Code, "Response's status code should be 201 (Created)")
	var createdUserResponse dto.HttpUserSuccess
	err = json.Unmarshal(w.Body.Bytes(), &createdUserResponse)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")

	// Create a test HTTP request for login
	userToLogin := dto.UserLoginDto{
		Email:    userToInsert.Email,
		Password: userToInsert.Password,
	}
	payload, err = json.Marshal(userToLogin)
	assert.NoError(t, err, "There should be no error here on json.Marshal(userToLogin)")
	req, err = http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(payload))
	assert.NoError(t, err, "There should be no error here on making login user request")
	// Create a test HTTP response recorder
	w = httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ = gin.CreateTestContext(w)
	c.Request = req

	// Call the UpdateUserByID handler method
	userHandler.Login(c)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code, "Response's status code should be 200 (ok)")

	// Validation the response body
	var loginResponseBody dto.HttpAccessTokenSuccess
	err = json.Unmarshal(w.Body.Bytes(), &loginResponseBody)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")
	assert.NotEmpty(t, loginResponseBody.Token, "Access token shouldn't be empty")
	// Create a test HTTP request with query parameter and payload
	userIDString := strconv.Itoa(int(createdUserResponse.User.ID))
	req, err = http.NewRequest(http.MethodDelete, "/api/v1/users/"+userIDString, nil)
	assert.NoError(t, err, "There should be no error here on making Deleting user by ID request")

	// Create a test HTTP response recorder
	w = httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: userIDString}}
	c.Request.Header.Set("Authorization", "Bearer "+loginResponseBody.Token)

	// Call the DeleteUserByID handler method
	userHandler.DeleteUserByID(c)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code, "Response's status code should be 200 (ok)")
}

func TestLogin(t *testing.T) {
	// Create a new instance of the user handler
	userHandler, cleanup := setupUserHandler()
	defer cleanup()

	// Create a test HTTP request with a JSON payload
	userToInsert := dto.UserCreateDto{
		Name:     "John Doe",
		Age:      25,
		Email:    "john@example.com",
		Password: "*Password123#",
	}

	payload, err := json.Marshal(userToInsert)
	assert.NoError(t, err, "There should be no error here on json.Marshal(userToInsert)")
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(payload))
	assert.NoError(t, err, "There should be no error here on making creation request")
	// Create a test HTTP response recorder
	w := httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the CreateUser handler method
	userHandler.CreateUser(c)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, w.Code, "Response's status code should be 201 (Created)")
	var createdUserResponse dto.HttpUserSuccess
	err = json.Unmarshal(w.Body.Bytes(), &createdUserResponse)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")

	// Create a test HTTP request for login
	userToLogin := dto.UserLoginDto{
		Email:    userToInsert.Email,
		Password: userToInsert.Password,
	}
	payload, err = json.Marshal(userToLogin)
	assert.NoError(t, err, "There should be no error here on json.Marshal(userToLogin)")
	req, err = http.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewBuffer(payload))
	assert.NoError(t, err, "There should be no error here on making login user request")
	// Create a test HTTP response recorder
	w = httptest.NewRecorder()

	// Create a Gin context from the request and response recorder
	c, _ = gin.CreateTestContext(w)
	c.Request = req

	// Call the UpdateUserByID handler method
	userHandler.Login(c)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code, "Response's status code should be 200 (ok)")

	// Validation the response body
	var loginResponseBody dto.HttpAccessTokenSuccess
	err = json.Unmarshal(w.Body.Bytes(), &loginResponseBody)
	assert.NoError(t, err, "There should be no error here on json.Unmarshal the response")
	assert.NotEmpty(t, loginResponseBody.Token, "Access token shouldn't be empty")
}
