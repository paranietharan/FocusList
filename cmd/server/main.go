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

	// test
	testHandler := &handler.TestHandler{}

	// user
	userRepo := repository.NewUserRepository(db)
	authSvc := &service.AuthService{
		UserRepo:  userRepo,
		CacheRepo: c,
	}
	authHandler := &handler.AuthHandler{AuthService: authSvc}

	// bucket
	bucketRepo := repository.NewTodoListBucketRepository(db)
	bucketSvc := &service.TodoListBucketService{
		TodoListBucketRepo: bucketRepo,
	}
	bucketHandler := &handler.BucketHandler{
		TodoListBucketService: bucketSvc,
	}

	// Initialize Gin router
	r := gin.Default()

	// Login and Singup routes
	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)
	r.POST("/verify-email", authHandler.VerifyEmail)
	r.POST("/reset-password", authHandler.ForgotPassword)
	r.POST("/confirm-reset-password", authHandler.ConfirmResetPassword)

	// Test Endpoint
	r.GET("/test", middleware.AuthMiddleware("user"), testHandler.UserDetailsCheck)

	// Bucket
	r.POST("/buckets", middleware.AuthMiddleware("user"), bucketHandler.CreateBucket)
	r.GET("/buckets", middleware.AuthMiddleware("user"), bucketHandler.GetBuckets)
	r.GET("/buckets/:bucketID", middleware.AuthMiddleware("user"), bucketHandler.GetBucketByID)
	r.PUT("/buckets/:bucketID", middleware.AuthMiddleware("user"), bucketHandler.UpdateBucketName)
	r.DELETE("/buckets/:bucketID", middleware.AuthMiddleware("user"), bucketHandler.DeleteBucket)

	// Bucket Management
	r.PUT("/buckets/:bucketID/users", middleware.AuthMiddleware("user"), bucketHandler.AddUserToBucket)
	r.DELETE("/buckets/:bucketID/users/:userEmail", middleware.AuthMiddleware("user"), bucketHandler.RemoveUserFromBucket)
	r.GET("/buckets/:bucketID/users/", middleware.AuthMiddleware("user"), bucketHandler.GetBucketUsers)

	//r.POST("/buckets/:bucketID/items", middleware.AuthMiddleware("user"), itemHandler.CreateItem)
	//r.GET("/buckets/:bucketID/items", middleware.AuthMiddleware("user"), itemHandler.GetItems)
	//r.GET("/buckets/:bucketID/items/:itemID", middleware.AuthMiddleware("user"), itemHandler.GetItemByID)
	//r.PUT("/buckets/:bucketID/items/:itemID", middleware.AuthMiddleware("user"), itemHandler.UpdateItem)
	//r.DELETE("/buckets/:bucketID/items/:itemID", middleware.AuthMiddleware("user"), itemHandler.DeleteItem)
	//r.PATCH("/buckets/:bucketID/items/:itemID/complete", middleware.AuthMiddleware("user"), itemHandler.MarkItemComplete)

	r.Run(":8080")
}
