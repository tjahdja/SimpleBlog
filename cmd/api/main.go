package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"SimpleBlog/internal/entity"
	"SimpleBlog/internal/handler"
	"SimpleBlog/internal/middleware"
	"SimpleBlog/internal/repository"
	"SimpleBlog/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=root dbname=simple_blog port=5433 sslmode=disable"
	}

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
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Post{}, &entity.Comment{})
	if err != nil {
		log.Fatalf("Failed to run schema migrations: %v", err)
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting application server", "port", ":8080")

	// Start the server
	if err := r.Run(":" + port); err != nil {
		slog.Error("failed to start server", "err", err)
	}
}
