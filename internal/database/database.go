package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

func GetDB(cfg Config) (*sql.DB, error) {
	var err error
	once.Do(func() {
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		dbInstance, err = sql.Open("postgres", connStr)
		if err != nil {
			return
		}
		err = dbInstance.Ping()
	})
	return dbInstance, err
}

type Config struct {
	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}
