package config

import (
	"log"
	"os"
	"strconv"
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

	// Redis configuration
	RedisAddr     string
	RedisUsername string
	RedisPassword string
	RedisDB       int
	RedisPort     string
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

		// Redis configuration
		redisAddr := os.Getenv("REDIS_ADDR")
		redisUsername := os.Getenv("REDIS_USERNAME")
		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisDB := os.Getenv("REDIS_DB")
		redisPort := os.Getenv("REDIS_PORT")

		if redisAddr == "" || redisDB == "" {
			log.Printf("Redis Details: REDIS_ADDR: %s, REDIS_USERNAME: %s, REDIS_PASSWORD: %s, REDIS_DB: %s", redisAddr, redisUsername, redisPassword, redisDB)
			log.Fatal("Redis configuration is not set in the environment variables")
		}

		dbInt, err := strconv.Atoi(redisDB)
		if err != nil {
			log.Fatalf("Invalid Redis DB number: %v", err)
		}

		cfg = &Config{
			JWTSecret: jwtSecret,

			// Database configuration
			DBUsername: username,
			DBPassword: password,
			DBHost:     host,
			DBPort:     port,
			DBName:     database,

			// Redis configuration
			RedisAddr:     redisAddr,
			RedisUsername: redisUsername,
			RedisPassword: redisPassword,
			RedisDB:       dbInt,
			RedisPort:     redisPort,
		}
	})

	return cfg
}
