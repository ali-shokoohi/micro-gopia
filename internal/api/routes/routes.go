package routes

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.Group("/")
			{
				HomeRoutes(v1)
			}
		}
	}
}
