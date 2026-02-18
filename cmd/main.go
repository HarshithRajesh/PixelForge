package main

import (
	"log"

	"github.com/HarshithRajesh/PixelForge/internal/config"
	"github.com/HarshithRajesh/PixelForge/internal/handler"
	"github.com/HarshithRajesh/PixelForge/internal/processor"
	"github.com/HarshithRajesh/PixelForge/internal/repository"
	"github.com/HarshithRajesh/PixelForge/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	db, _ := config.ConnectDB()
	userRepo := repository.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()
	r.GET("/health", processor.Health)
	r.POST("/signup", userHandler.SignUp)
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
