package database

import (
	"database/sql"
	"fmt"
	"gojek/library/config"

	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sql.DB, error) {
	cfg := config.NewConfig()
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	return sql.Open("postgres", connStr)
}
