package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ali-shokoohi/micro-gopia/config"
	"github.com/ali-shokoohi/micro-gopia/docs"
	"github.com/ali-shokoohi/micro-gopia/internal/api/routes"
	"github.com/ali-shokoohi/micro-gopia/pkg/migrations"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @title           Gopia
// @version         1.0
// @description     A RestAPI with Go
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Reading config file
	if err := config.Confs.Load(""); err != nil {
		fmt.Printf("We have an error in loading config: %s", err.Error())
		return
	}
	// Watch config file if it change
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		if err := config.Confs.Load(""); err != nil {
			fmt.Printf("We have an error in loading config: %s", err.Error())
			return
		}
	})

	log.Printf("Confs is: %v", config.Confs)

	// Check system arguments. like: main.go migrate
	args := os.Args
	if len(args) > 1 && args[1] == "migrate" {
		err := migrations.AutoMigrateDB()
		if err == nil {
			log.Println("Migration ended successfully!")
		} else {
			log.Printf("An error in migration, Error: %s", err.Error())
		}
		return
	}

	// programmatically set swagger info
	host := config.Confs.Service.HTTP.Host
	port := config.Confs.Service.HTTP.Port
	baseHost := fmt.Sprintf("%s:%s", host, port)
	docs.SwaggerInfo.Host = baseHost
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Import routes
	r := gin.Default()
	routes.Routes(r)

	// Start listening
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("We have an error in start listening: %s", err.Error())
		return
	}
}
