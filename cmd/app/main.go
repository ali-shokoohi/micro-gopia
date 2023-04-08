package main

import (
	"github.com/ali-shokoohi/micro-gopia/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.Routes(r)
	r.Run(":8080")
}
