package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ali-shokoohi/micro-gopia/internal/domain/dto"
	"github.com/ali-shokoohi/micro-gopia/internal/services"
	"github.com/ali-shokoohi/micro-gopia/scripts"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	UpdateUserByID(c *gin.Context)
	DeleteUserByID(c *gin.Context)
	Login(c *gin.Context)
}

// UserService represents the service that interacts with the user_repository.
type userHandler struct {
	userService services.UserService
}

// NewUserHandler returns a new instance of userHandler.
func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary      Create an user
// @Description  Create an user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        {object}   body      dto.UserCreateDto  true  "User creation body"
// @Success      201  {object}  dto.HttpUserSuccess
// @Failure      400  {object}  dto.HttpFailures
// @Failure      500  {object}  dto.HttpFailure
// @Router       /users [post]
func (uh *userHandler) CreateUser(c *gin.Context) {
	var userCreateDto *dto.UserCreateDto
	if err := c.BindJSON(&userCreateDto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailures{
			Status: "failure",
			Errors: []error{errors.New("Can't bind the request body. Check it!")},
		})
		return
	}
	userEntity, errs := uh.userService.CreateUser(c, userCreateDto)
	if errs != nil && len(errs) > 0 {
		switch errs[0].(type) {
		case dto.InternalServerError:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  errs[0],
			})
			return
		case dto.BadRequestError:
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailures{
				Status: "failure",
				Errors: errs,
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailures{
				Status: "failure",
				Errors: errs,
			})
			return
		}
	}
	userViewDto, err := scripts.UserEntityToUserViewDto(&userEntity)
	if err != nil {
		log.Printf("We could't transform user entity to user view, Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("We could't transform user entity to user view dto"),
		})
		return
	}
	c.JSON(http.StatusCreated, &dto.HttpUserSuccess{
		Status: "success",
		User:   userViewDto,
	})
}

// GetUsers godoc
// @Summary      List users
// @Description  Get a list of users with pagination
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page   path int false "Page number of users" default(0)
// @Param        limits path int false "limit of users in each page" default(10)
// @Success      200  {object}  dto.HttpUsersSuccess
// @Failure      400  {object}  dto.HttpFailure
// @Failure      500  {object}  dto.HttpFailure
// @Router       /users [get]
func (uh *userHandler) GetUsers(c *gin.Context) {
	var page, limits uint64 = 0, 10
	var err error
	if pageString, hasPage := c.GetQuery("page"); hasPage && pageString != "" {
		page, err = strconv.ParseUint(pageString, 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
				Status: "failure",
				Error:  errors.New("Please enter a valid page number"),
			})
			return
		}
	}
	if limitsString, hasLimits := c.GetQuery("limits"); hasLimits && limitsString != "" {
		limits, err = strconv.ParseUint(limitsString, 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
				Status: "failure",
				Error:  errors.New("Please enter a valid limits number"),
			})
			return
		}
	}
	userEntities, err := uh.userService.GetUsers(c, uint(page), uint(limits))
	if err != nil {
		switch err.(type) {
		case dto.BadRequestError:
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		case dto.InternalServerError:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		}

	}
	userViewDtos, err := scripts.UserEntitiesToUserViewDtos(userEntities)
	if err != nil {
		log.Printf("We could't transform user entities to user view dtos, Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("We could't transform user entities to user view dtos"),
		})
		return
	}
	c.JSON(http.StatusOK, &dto.HttpUsersSuccess{
		Status: "success",
		Users:  userViewDtos,
	})
}

// GetUserByID godoc
// @Summary      Show an user
// @Description  Get an user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path int true "User ID"
// @Success      200  {object}  dto.HttpUserSuccess
// @Failure      400  {object}  dto.HttpFailure
// @Failure      500  {object}  dto.HttpFailure
// @Router       /users/{id} [get]
func (uh *userHandler) GetUserByID(c *gin.Context) {
	idString, hasID := c.Params.Get("id")
	if !hasID || idString == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("Enter the ID correctly!"),
		})
		return
	}
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("Enter a valid ID number!"),
		})
		return
	}
	userEntity, err := uh.userService.GetUserByID(c, uint(id))
	if err != nil {
		switch err.(type) {
		case dto.BadRequestError:
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		case dto.InternalServerError:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		}
	}
	userViewDto, err := scripts.UserEntityToUserViewDto(&userEntity)
	if err != nil {
		log.Printf("We could't transform user entity to user view, Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("We could't transform user entity to user view dto"),
		})
		return
	}
	c.JSON(http.StatusOK, &dto.HttpUserSuccess{
		Status: "success",
		User:   userViewDto,
	})
}

// UpdateUserByID godoc
// @Summary      Update an user
// @Description  Update an user by ID and new data
// @Tags         users
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Param        {object}   body      dto.UserUpdateDto  true  "User update body"
// @Param        id   path int true "User ID"
// @Success      200  {object}  dto.HttpUserSuccess
// @Failure      400  {object}  dto.HttpFailures
// @Failure      403  {object}  dto.HttpFailure
// @Failure      500  {object}  dto.HttpFailure
// @Router       /users/{id} [put]
func (uh *userHandler) UpdateUserByID(c *gin.Context) {
	idString, hasID := c.Params.Get("id")
	if !hasID || idString == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailures{
			Status: "failure",
			Errors: []error{errors.New("Enter the ID correctly!")},
		})
		return
	}
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailures{
			Status: "failure",
			Errors: []error{errors.New("Enter a valid ID number!")},
		})
		return
	}
	claim, err := scripts.CurrentTokenClaim(c)
	if err != nil {
		log.Printf("Can't get sender ID from the request's token: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("Can't get sender ID from the request's token"),
		})
		return
	}
	if claim.UserID != uint(id) {
		c.AbortWithStatusJSON(http.StatusForbidden, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("Permission denied"),
		})
		return
	}
	var userUpdateDto *dto.UserUpdateDto
	if err := c.BindJSON(&userUpdateDto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailures{
			Status: "failure",
			Errors: []error{errors.New("Can't bind the request body. Check it!")},
		})
		return
	}
	userEntity, errs := uh.userService.UpdateUserByID(c, uint(id), userUpdateDto)
	if errs != nil && len(errs) > 0 {
		switch errs[0].(type) {
		case dto.InternalServerError:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  errs[0],
			})
			return
		case dto.BadRequestError:
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailures{
				Status: "failure",
				Errors: errs,
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailures{
				Status: "failure",
				Errors: errs,
			})
			return
		}
	}
	userViewDto, err := scripts.UserEntityToUserViewDto(&userEntity)
	if err != nil {
		log.Printf("We could't transform user entity to user view, Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("We could't transform user entity to user view dto"),
		})
		return
	}
	c.JSON(http.StatusOK, &dto.HttpUserSuccess{
		Status: "success",
		User:   userViewDto,
	})
}

// DeleteUserByID godoc
// @Summary      Delete an user
// @Description  Delete an user By ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @Param        id   path int true "User ID"
// @Success      200  {object}  dto.HttpSuccess
// @Failure      400  {object}  dto.HttpFailure
// @Failure      403  {object}  dto.HttpFailure
// @Failure      500  {object}  dto.HttpFailure
// @Router       /users/{id} [delete]
func (uh *userHandler) DeleteUserByID(c *gin.Context) {
	idString, hasID := c.Params.Get("id")
	if !hasID || idString == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("Enter the ID correctly!"),
		})
		return
	}
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("Enter a valid ID number!"),
		})
		return
	}
	claim, err := scripts.CurrentTokenClaim(c)
	if err != nil {
		log.Printf("Can't get sender ID from the request's token: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("Can't get sender ID from the request's token"),
		})
		return
	}
	if claim.UserID != uint(id) {
		c.AbortWithStatusJSON(http.StatusForbidden, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("Permission denied"),
		})
		return
	}
	err = uh.userService.DeleteUserByID(c, uint(id))
	if err != nil {
		switch err.(type) {
		case dto.BadRequestError:
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		case dto.InternalServerError:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		}
	}
	c.JSON(http.StatusOK, &dto.HttpSuccess{
		Status:  "success",
		Message: "User deleted successfully",
	})
}

// Login godoc
// @Summary      User login
// @Description  Login with username and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        {object}   body      dto.UserLoginDto  true  "User login body"
// @Success      200  {object}  dto.HttpAccessTokenSuccess
// @Failure      400  {object}  dto.HttpFailures
// @Failure      500  {object}  dto.HttpFailure
// @Router       /users/login [post]
func (uh *userHandler) Login(c *gin.Context) {
	var userLoginDto dto.UserLoginDto
	if err := c.BindJSON(&userLoginDto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
			Status: "failure",
			Error:  errors.New("Can't bind the request body. Check it!"),
		})
		return
	}
	token, err := uh.userService.Login(c, &userLoginDto)
	if err != nil {
		switch err.(type) {
		case dto.BadRequestError:
			c.AbortWithStatusJSON(http.StatusBadRequest, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		case dto.InternalServerError:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, &dto.HttpFailure{
				Status: "failure",
				Error:  err,
			})
			return
		}
	}
	c.JSON(http.StatusOK, &dto.HttpAccessTokenSuccess{
		Status: "success",
		Token:  token,
	})
}
