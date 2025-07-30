package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string

	// Database configuration
	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println(".env file not found, continuing with system env variables")
		}

		// Get the Database Credentials from environment variables
		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		database := os.Getenv("DB_NAME")

		if username == "" || password == "" || host == "" || port == "" || database == "" {
			log.Printf("Database Details: DB_USERNAME: %s, DB_PASSWORD: %s, DB_HOST: %s, DB_PORT: %s, DB_NAME: %s", username, password, host, port, database)
			log.Fatal("Database connection details are not set in the environment variables")
		}

		jwtSecret := os.Getenv("JWT_SECRET_KEY")
		if jwtSecret == "" {
			log.Fatal("Missing JWT_SECRET_KEY in environment")
		}

		cfg = &Config{
			JWTSecret:  jwtSecret,
			DBUsername: username,
			DBPassword: password,
			DBHost:     host,
			DBPort:     port,
			DBName:     database,
		}
	})

	return cfg
}
