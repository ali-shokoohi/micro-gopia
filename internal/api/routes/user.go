package routes

import (
	"github.com/ali-shokoohi/micro-gopia/internal/api/handlers"
	"github.com/ali-shokoohi/micro-gopia/internal/auth"
	"github.com/ali-shokoohi/micro-gopia/internal/datastore"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/repositories"
	"github.com/ali-shokoohi/micro-gopia/internal/services"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	auth := auth.NewAuth()
	userRepository := repositories.NewUserRepository(datastore.DB)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	r.POST("/", userHandler.CreateUser)
	r.GET("/", userHandler.GetUsers)
	r.GET("/:id", userHandler.GetUserByID)
	r.POST("/login", userHandler.Login)
	r.Use(auth.JwtAuthMiddleware())
	r.PUT("/:id", userHandler.UpdateUserByID)
	r.DELETE("/:id", userHandler.DeleteUserByID)
	return r
}
