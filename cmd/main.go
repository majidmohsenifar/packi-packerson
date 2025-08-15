package main

import (
	"fmt"
	"job-test/api"
	"job-test/config"

	_ "job-test/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Pack-Service
// @version 1.0
// @description This is a sample server for a Gin application.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @schemes http
// @BasePath /

func main() {
	config := config.LoadConfig()
	server := gin.Default()
	api.InitApi(server)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	if err := server.Run(fmt.Sprintf(":%d", config.App.Port)); err != nil {
		panic("Error run project on port!")
	}
}
