package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func MigrateFromFile(db *sql.DB, migrationPath string) error {
	sqlBytes, err := os.ReadFile(filepath.Clean(migrationPath))
	if err != nil {
		return err
	}
	_, err = db.Exec(string(sqlBytes))
	return err
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	if username == "" || password == "" || host == "" || port == "" || database == "" {
		log.Fatal("Database connection details are not set in the environment variables")
		log.Printf("DB_USERNAME: %s, DB_PASSWORD: %s, DB_HOST: %s, DB_PORT: %s, DB_NAME: %s", username, password, host, port, database)
		return
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		database,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open DB connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	err = MigrateFromFile(db, "migrations/001_create_table.sql")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database migration completed successfully")
}
