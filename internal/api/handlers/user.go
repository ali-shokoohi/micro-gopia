package handlers

import (
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
}

// UserService represents the service that interacts with the user_repository.
type userHandler struct {
	userService services.UserService
}

// NewUserService returns a new instance of UserService.
func NewUserService(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (uh *userHandler) CreateUser(c *gin.Context) {
	var userCreateDto *dto.UserCreateDto
	if err := c.BindJSON(&userCreateDto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  "Can't bind the request body. Check it!",
		})
		return
	}
	userEntity, errs := uh.userService.CreateUser(c, userCreateDto)
	if errs != nil && len(errs) > 0 {
		strs := scripts.ErrorsSliceToStringSlice(errs)
		if strs != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": "bad",
				"errors": strs,
			})
			return
		}
	}
	userViewDto, err := scripts.UserEntityToUserViewDto(userEntity)
	if err != nil {
		log.Printf("We could't transform user entity to user view, Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "bad",
			"error":  "We could't transform user entity to user view dto",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"user":   userViewDto,
	})
}

func (uh *userHandler) GetUsers(c *gin.Context) {
	var page, limits uint64 = 0, 10
	var err error
	if pageString, hasPage := c.GetQuery("page"); hasPage && pageString != "" {
		page, err = strconv.ParseUint(pageString, 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": "bad",
				"error":  "Please enter a valid page number",
			})
			return
		}
	}
	if limitsString, hasLimits := c.GetQuery("limits"); hasLimits && limitsString != "" {
		limits, err = strconv.ParseUint(limitsString, 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": "bad",
				"error":  "Please enter a valid limits number",
			})
			return
		}
	}
	userEntities, err := uh.userService.GetUsers(c, uint(page), uint(limits))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  err.Error(),
		})
		return
	}
	userViewDtos, err := scripts.UserEntitiesToUserViewDtos(userEntities)
	if err != nil {
		log.Printf("We could't transform user entities to user view dtos, Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "bad",
			"error":  "We could't transform user entities to user view dtos",
		})
		return
	}
	if userViewDtos == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"users":  []dto.UserViewDto{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"users":  userViewDtos,
	})

}

func (uh *userHandler) GetUserByID(c *gin.Context) {
	idString, hasID := c.Params.Get("id")
	if !hasID || idString == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  "Enter the ID correctly!",
		})
		return
	}
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  "Enter a valid ID number!",
		})
		return
	}
	userEntity, err := uh.userService.GetUserByID(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  err.Error(),
		})
		return
	}
	userViewDto, err := scripts.UserEntityToUserViewDto(userEntity)
	if err != nil {
		log.Printf("We could't transform user entity to user view, Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "bad",
			"error":  "We could't transform user entity to user view dto",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"user":   userViewDto,
	})
}

func (uh *userHandler) UpdateUserByID(c *gin.Context) {
	idString, hasID := c.Params.Get("id")
	if !hasID || idString == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  "Enter the ID correctly!",
		})
		return
	}
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  "Enter a valid ID number!",
		})
		return
	}
	var userUpdateDto *dto.UserUpdateDto
	if err := c.BindJSON(&userUpdateDto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  "Can't bind the request body. Check it!",
		})
		return
	}
	userEntity, errs := uh.userService.UpdateUserByID(c, uint(id), userUpdateDto)
	if errs != nil && len(errs) > 0 {
		strs := scripts.ErrorsSliceToStringSlice(errs)
		if strs != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": "bad",
				"errors": strs,
			})
			return
		}
	}
	userViewDto, err := scripts.UserEntityToUserViewDto(userEntity)
	if err != nil {
		log.Printf("We could't transform user entity to user view, Error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": "bad",
			"error":  "We could't transform user entity to user view dto",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"user":   userViewDto,
	})
}

func (uh *userHandler) DeleteUserByID(c *gin.Context) {
	idString, hasID := c.Params.Get("id")
	if !hasID || idString == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  "Enter the ID correctly!",
		})
		return
	}
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  "Enter a valid ID number!",
		})
		return
	}
	err = uh.userService.DeleteUserByID(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status": "bad",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "User deleted successfully",
	})
}