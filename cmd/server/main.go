package main

import (
	"FocusList/internal/cache"
	"FocusList/internal/config"
	"FocusList/internal/database"
	"FocusList/internal/handler"
	"FocusList/internal/middleware"
	"FocusList/internal/repository"
	"FocusList/internal/service"
	"context"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {
	// Load the configuration
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("Failed to load configuration")
	}

	dbconfig := database.Config{
		DBUsername: cfg.DBUsername,
		DBPassword: cfg.DBPassword,
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBName:     cfg.DBName,
	}
	db, err := database.GetDB(dbconfig)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c := cache.NewCacheService(cache.RedisConfig{
		Addr:     cfg.RedisAddr,
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		Port:     cfg.RedisPort,
	}, context.Background())

	err = c.Ping()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	userRepo := repository.NewUserRepository(db)
	authSvc := &service.AuthService{
		UserRepo:  userRepo,
		CacheRepo: c,
	}
	authHandler := &handler.AuthHandler{AuthService: authSvc}

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/verify-email", authHandler.VerifyEmail)
	r.POST("/login", authHandler.Login)

	r.GET("/profile", middleware.AuthMiddleware("super_admin"), func(c *gin.Context) {
		email := c.GetString("userEmail")
		c.JSON(200, gin.H{
			"email": email,
			"role":  c.GetString("userRole"),
		})
	})

	r.GET("/hi", middleware.AuthMiddleware("user"), func(c *gin.Context) {
		email := c.GetString("userEmail")
		c.JSON(200, gin.H{
			"email":   email,
			"message": "Hello, user!",
			"role":    c.GetString("userRole"),
		})
	})

	r.Run(":8080")
}
