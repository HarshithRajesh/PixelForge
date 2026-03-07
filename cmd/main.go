package main

import (
	"log"

	"github.com/HarshithRajesh/PixelForge/internal/config"
	"github.com/HarshithRajesh/PixelForge/internal/handler"
	"github.com/HarshithRajesh/PixelForge/internal/middleware"
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
	imageService := processor.NewImageManagement()
	imageHandler := handler.NewImageManagementHandler(imageService)
	r := gin.Default()
	r.GET("/health", processor.Health)
	r.POST("/signup", userHandler.SignUp)
	r.POST("/login", userHandler.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", processor.Profile)
		protected.POST("/image", imageHandler.ImageUpload)
	}

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
