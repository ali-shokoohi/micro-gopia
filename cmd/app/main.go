package main

import (
	"fmt"

	"github.com/ali-shokoohi/micro-gopia/config"
	"github.com/ali-shokoohi/micro-gopia/internal/api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Reading config file
	if err := config.Confs.Load(); err != nil {
		fmt.Printf("We have an error in loading config: %s", err.Error())
		return
	}

	// Import routes
	r := gin.Default()
	routes.Routes(r)

	// Start listening
	if err := r.Run(fmt.Sprintf("%s:%s", config.Confs.Service.HTTP.Host, config.Confs.Service.HTTP.Port)); err != nil {
		fmt.Printf("We have an error in start listening: %s", err.Error())
		return
	}
}
