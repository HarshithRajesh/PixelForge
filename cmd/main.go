package main

import (
	"log"

	"github.com/HarshithRajesh/PixelForge/internal/config"
	"github.com/HarshithRajesh/PixelForge/internal/handler"
	"github.com/HarshithRajesh/PixelForge/internal/middleware"
	"github.com/HarshithRajesh/PixelForge/internal/processor"
	"github.com/HarshithRajesh/PixelForge/internal/repository"
	"github.com/HarshithRajesh/PixelForge/internal/user"
	"github.com/HarshithRajesh/PixelForge/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	rds := config.NewRedis()

	db, _ := config.ConnectDB()
	userRepo := repository.NewUserRepository(db)
	userService := user.NewUserService(userRepo, rds)
	userHandler := handler.NewUserHandler(userService, rds)

	store := storage.NewStorageRepository("storage")

	imageService := processor.NewImageManagement(userRepo, store)
	imageHandler := handler.NewImageManagementHandler(imageService)
	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	// AllowOrigins:     []string{os.Getenv("FRONTEND_ORIGIN")},
	// 	AllowMethods:     []string{"GET", "POST"},
	// 	AllowHeaders:     []string{"Content-Type", "Authorization"},
	// 	AllowCredentials: true,
	// }))

	r.GET("/health", processor.Health)
	r.POST("/signup", userHandler.SignUp)
	r.POST("/login", userHandler.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(rds))
	{
		protected.GET("/profile", processor.Profile)
		protected.POST("/image", imageHandler.ImageUpload)
	}

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
