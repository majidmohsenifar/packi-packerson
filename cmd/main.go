package main

import (
	"embed"
	"job-test/api"
	"os"

	_ "job-test/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed static/*
var StaticFiles embed.FS

// @title Pack-Service
// @version 1.0
// @description This is a sample server for a Gin application.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @schemes http
// @BasePath /

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default for local dev
	}
	// config := config.LoadConfig()

	server := gin.Default()
	api.InitApi(server, StaticFiles)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	if err := server.Run(":" + port); err != nil {
		panic("Error run project on port!")
	}
}
