package routes

import (
	"github.com/ali-shokoohi/micro-gopia/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func HomeRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	homeHandlers := &handlers.HomeHandler{}
	r.GET("/", homeHandlers.Root)
	return r
}
