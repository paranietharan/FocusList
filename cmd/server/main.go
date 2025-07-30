package main

import (
	"FocusList/internal/config"
	"FocusList/internal/handler"
	"FocusList/internal/middleware"
	"FocusList/internal/repository"
	"FocusList/internal/service"
	"database/sql"
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

	connStr := "postgres://" + cfg.DBUsername + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authSvc := &service.AuthService{UserRepo: userRepo}
	authHandler := &handler.AuthHandler{AuthService: authSvc}

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	r.GET("/profile", middleware.AuthMiddleware("super_admin"), func(c *gin.Context) {
		email := c.GetString("userEmail")
		c.JSON(200, gin.H{
			"email": email,
			"role":  c.GetString("userRole"),
		})
	})

	r.Run(":8080")
}
