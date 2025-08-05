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

	// TODO: Check endpoint Need to remove this in production
	r.GET("/profile", middleware.AuthMiddleware("super_admin"), func(c *gin.Context) {
		email := c.GetString("userEmail")
		c.JSON(200, gin.H{
			"email": email,
			"role":  c.GetString("userRole"),
		})
	})
	// TODO: Check endpoint Need to remove this in production
	r.GET("/hi", middleware.AuthMiddleware("user"), func(c *gin.Context) {
		email := c.GetString("userEmail")
		c.JSON(200, gin.H{
			"email":   email,
			"message": "Hello, user!",
			"role":    c.GetString("userRole"),
		})
	})

	// TodoList Bucket Management
	r.POST("/buckets", middleware.AuthMiddleware("user"), bucketHandler.CreateBucket)
	r.GET("/buckets", middleware.AuthMiddleware("user"), bucketHandler.GetBuckets)
	r.GET("/buckets/:bucketID", middleware.AuthMiddleware("user"), bucketHandler.GetBucketByID)
	r.PUT("/buckets/:bucketID", middleware.AuthMiddleware("user"), bucketHandler.UpdateBucketName)
	r.PUT("/buckets/:bucketID/users", middleware.AuthMiddleware("user"), bucketHandler.AddUserToBucket)
	r.DELETE("/buckets/:bucketID", middleware.AuthMiddleware("user"), bucketHandler.DeleteBucket)
	//Update a user's permission in a bucket
	//Remove a user from a bucket
	//List all users and their permissions in the bucket
	/*


		-------------------------------------------------------------------


		-------------------------------------------------------------------
		3️⃣ TodoList Item Management (/buckets/:bucketID/items)
		-------------------------------------------------------------------
		- [POST]    /buckets/:bucketID/items
		    → Create a new todo item in the bucket

		- [GET]     /buckets/:bucketID/items
		    → List all items in the bucket

		- [GET]     /buckets/:bucketID/items/:itemID
		    → Get a specific todo item

		- [PUT]     /buckets/:bucketID/items/:itemID
		    → Update a todo item (description, status)

		- [DELETE]  /buckets/:bucketID/items/:itemID
		    → Delete a todo item

		-------------------------------------------------------------------
		4️⃣ Advanced Features (Optional Enhancements)
		-------------------------------------------------------------------
		- [PATCH]   /buckets/:bucketID/items/:itemID/complete
		    → Mark item as complete

		- [GET]     /buckets/:bucketID/items?isCompleted=true&q=urgent
		    → Filter/search items by status or keywords

		- Audit Logs
		    → Track who updated what and when (optional future enhancement)

		- Pagination & Sorting
		    → Add support: ?page=1&limit=10&sort=createdAt

		-------------------------------------------------------------------
		🔐 Authorization Middleware Suggestions:
		- Use "read", "write", "execute" permissions from `TodoListBucketUser`
		- Wrap route groups using `middleware.AuthMiddleware(role...)`
		- Determine permission on a per-bucket basis using middleware or service layer

		-------------------------------------------------------------------

		🛠️ Suggestion:
		Implement these endpoints incrementally starting with Bucket Management, then Collaboration, followed by Item Management. This keeps the system testable and avoids scope creep.

	*/

	r.Run(":8080")
}
