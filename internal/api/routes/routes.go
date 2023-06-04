package routes

import (
	"github.com/ali-shokoohi/micro-gopia/internal/api/middlewares"
	"github.com/gin-gonic/gin"                 // swagger embed files
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func Routes(r *gin.Engine) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.Use(middlewares.CORSMiddleware())
			v1.Group("/")
			{
				HomeRoutes(v1)
			}
			users := v1.Group("/users")
			{
				UserRoutes(users)
			}
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
