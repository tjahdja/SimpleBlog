package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/tjahdja/SimpleBlog/internal/handler"
	"github.com/tjahdja/SimpleBlog/internal/middleware"
	"github.com/tjahdja/SimpleBlog/internal/repository"
	"github.com/tjahdja/SimpleBlog/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/tjahdja/SimpleBlog/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           SimpleBlog REST API
// @version         1.0
// @description     A production-ready Clean Architecture blogging backend built with Go and Gin.
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apiKey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type 'Bearer ' followed by your JWT token payload.

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dsn := os.Getenv("DATABASE_URL")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super_secret_blog_key_change_me_in_production"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewGORMUserRepository(db)
	postRepo := repository.NewGORMPostRepository(db)
	commentRepo := repository.NewGORMCommentRepository(db)

	userService := service.NewGORMUserService(userRepo, jwtSecret)
	postService := service.NewGORMPostService(postRepo)
	commentService := service.NewGORMCommentService(commentRepo)

	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	r.GET("/posts", postHandler.ListPosts)
	r.GET("/posts/:id", postHandler.GetPostByID)
	r.GET("/posts/:id/comments", commentHandler.GetCommentsByPostID)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.POST("/posts", postHandler.CreatePost)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)

		protected.POST("/posts/:id/comments", commentHandler.CreateComment)
		protected.DELETE("/comments/:id", commentHandler.DeleteComment)
	}

	slog.Info("Starting application server", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
